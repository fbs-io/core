/*
 * @Author: reel
 * @Date: 2023-05-16 22:16:53
 * @LastEditors: reel
 * @LastEditTime: 2023-11-10 06:55:08
 * @Description: 关系数据库配置
 */
package rdb

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/fbs-io/core/logx"
	"github.com/fbs-io/core/pkg/env"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/dsn"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	// 数据分区模式
	// 默认值, 不分区
	SHADING_MODEL_NOT int8 = iota

	// 数据都在一张表中
	SHADING_MODEL_ONE

	// 按表分区, 根据分区字段值, 设置表后缀, 如果表设置了分区
	SHADING_MODEL_TABLE

	// 按库(schema)分区, 根据分区字段值, 设置不同的库名(schema)后缀
	SHADING_MODEL_DB

	TABLE_INIT_ALL = "all"

	// 分区列名
	SHADING_KEY             = "sk"
	TX_SHADING_MODEL_KEY    = "tx_sharding_model"
	TX_SHADING_TABLE_KEY    = "tx_sharding_table"
	TX_CONDITION_BUILD_KEY  = "tx_condition_build"
	TX_SUB_QUERY_COLUMN_KEY = "tx_sub_query_column"
	TX_DATA_PERMISSION_KEY  = "tx_data_permission"
)

type rdbStore struct {
	db                  *gorm.DB
	dial                gorm.Dialector
	statTab             Tabler
	tablers             map[string]Tabler
	isRunning           bool
	migrateList         []func() error
	shardingTable       map[string][]string // 仅仅写入注册了分区表的表
	shardingModel       int8
	shardingSuffixs     []interface{}   //分区后缀
	shardingAllTable    map[string]bool // 模型注册时, 只要包含了分区字段的表, 都会写入到该map中, 用于回调函数判断是否增加分区字段
	dataPermissionTable map[string]bool // 模型注册时, 只要包含了权限字段的表, 都会写入到该map中, 用于回调函数判断是否增加分区字段
}

type Store interface {
	rdbP()
	DB() *gorm.DB

	// service.Service接口
	Name() string
	Stop() error
	Start() error
	Status() int8

	// DB 相关
	SetConfig(fs ...dsn.DsnFunc) error

	// 注册表结构, 同时允许注册时写入函数, 如初始化创建部分数据等
	Register(t Tabler, fs ...RegisterFunc) Store

	// 根据条件结构体生成查询tx
	BuildQuery(cb *Condition) (tx *gorm.DB)

	// 通过参数反射值构建普通查询tx用于应用端使用
	// 大多数查询通过使用该方法即可
	BuildQueryWithParams(params reflect.Value) *gorm.DB

	// 通过参数反射值构建普通查询tx用于应用端使用
	//
	// 此方法适用于表中有ID(主键)的字段, 优化了翻页查询性能
	BuildQueryWihtSubQryID(cb *Condition) (tx *gorm.DB)

	// 设置分区模式
	//
	// SHADING_MODEL_NOT : 不分区, 默认值, 数据都在一张表中,
	//
	// SHADING_MODEL_TABLE: 按表分区, 根据分区字段值, 设置表后缀
	//
	// TODO:SHADING_MODEL_DB: 按库(schema)分区, 根据分区字段值, 设置不同的库名(schema)后缀
	//
	// 该模式适用于使用cores 上下文ctx.TX()方式生成的 gorm.DB, 且在上下文中传入了分区字段, 会自动构建查询条件, 配合 ShardingModel使用,可以自动写入分区字段
	//
	// 如果直接使用gorm.DB, 该设置并不会生效
	SetShardingModel(model int8, suffix []interface{})

	// 获取分区模式
	ShardingModel() (model int8)

	// 获取可分区表
	ShardingTable() map[string][]string

	// SHADING_MODEL_TABLE(按表分区) 模式下, 自定义添加可用于分区的表
	AddShardingTable(table string)

	// 添加启动前的前置执行程序
	AddMigrateList(fs ...func() error)

	//添加分区后缀, 同时会重置表分区,自动迁移新增表分区的结构
	AddShardingSuffixs(suffixs string) (err error)
}

var _ Store = (*rdbStore)(nil)

var rdb = &rdbStore{
	db:                  &gorm.DB{},
	tablers:             make(map[string]Tabler, 1000),
	migrateList:         make([]func() error, 0, 100),
	shardingTable:       make(map[string][]string, 100),
	shardingSuffixs:     make([]interface{}, 0, 100),
	shardingAllTable:    make(map[string]bool, 100),
	dataPermissionTable: make(map[string]bool, 100),
}

func New() (s Store) {
	return rdb
}

func (store *rdbStore) rdbP()        {}
func (store *rdbStore) Name() string { return "RDB" }
func (store *rdbStore) DB() *gorm.DB { return store.db }

func (store *rdbStore) Start() (err error) {
	if store.dial == nil {
		return errorx.New("没有可用的DSN")
	}

	store.db, err = gorm.Open(store.dial, &gorm.Config{
		Logger: logx.DB.LogMode(env.Active().GormLogLevel()),
	})
	if err != nil {
		return errorx.Wrap(err, "数据库链接失败")
	}
	// 配置链接时间等
	sqlDB, err := store.db.DB()
	if err != nil {
		return err
	}

	// TODO: 自定义配置连接池大小
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 注册回调函数
	store.registerCallbacks()
	err = store.autoMigrate()
	if err != nil {
		return errorx.Wrap(err, "表迁移失败")
	}
	store.isRunning = true
	return
}

// 通过查询获取当前db状态
// 如果没有表, 默认当前不可用
func (store *rdbStore) Status() int8 {
	if !store.isRunning {
		return -1
	}
	if store.statTab == nil {
		return -1
	}
	err := store.db.FirstOrCreate(&(store.statTab)).Error
	if err != nil {
		return -1
	}
	return 1
}

func (store *rdbStore) Stop() error {
	sqlDB, err := store.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (store *rdbStore) SetConfig(optfs ...dsn.DsnFunc) (err error) {
	rdbDsn := dsn.NewDBDsn()
	for _, optf := range optfs {
		optf(rdbDsn)
	}
	link := rdbDsn.Link()
	if link == "" {
		return errorx.New("dsn 为空")
	}

	// 本地 db 进行的判断, db 为 sqlite
	switch rdbDsn.Type {
	case dsn.DSN_TYPE_LOCAL, dsn.DSN_TYPE_SQLITE: // 本地db使用 sqlite
		_, err = os.Stat(rdbDsn.Path)
		if err != nil {
			err = os.MkdirAll(rdbDsn.Path, 0766)
			if err != nil {
				return err
			}
		}
		store.dial = sqlite.Open(link)

	case dsn.DSN_TYPE_PGSQL:
		store.dial = postgres.Open(link)
	case dsn.DSN_TYPE_MYSQL:
		store.dial = mysql.Open(link)
	}
	return nil
}

// 迁移表
// 如果表结构发生变化, 需要启动时设置 dbinit=tableaName 参数, 删除旧表并创建新表
// 在表创建后, 可以执行一些自定义的方法, 主要用于初始化数据写入
func (store *rdbStore) Register(t Tabler, fs ...RegisterFunc) Store {
	store.tablers[t.TableName()] = t

	store.migrateList = append(store.migrateList, func() (err error) {
		// 根据实际使用, 系统资源表将在最后被加载
		store.statTab = t
		err = store.AutoShardingTable(t.TableName())
		if err != nil {
			return
		}
		// 从命令行重置所有表或部分表
		if strings.Contains(env.Active().DBInit(), t.TableName()) ||
			env.Active().DBInit() == TABLE_INIT_ALL {
			if store.db.Migrator().HasTable(t) {
				err = store.db.Migrator().DropTable(t)
				if err != nil {
					return
				}
				err = store.db.AutoMigrate(t)
				if err != nil {
					return
				}
				for _, f := range fs {
					err = f()
					if err != nil {
						return
					}
				}
			}
		} else {
			err = store.db.AutoMigrate(t)
			if err != nil {
				return
			}
		}

		return

	})
	return store
}

func (store *rdbStore) autoMigrate() (err error) {
	for _, fs := range store.migrateList {
		e := fs()
		if e != nil {
			err = errorx.Wrap(e, "表迁移操作失败")
		}
	}
	return
}

type Tabler interface {
	TableName() string
}

// 请传入一个有 TableName 方法的 结构体切片/数组
//
// 例子: []*A{&A{C1:"abc"}, &A{C1:"b"}}
func (store *rdbStore) CreateInBatches(ts interface{}) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errorx.New(fmt.Sprintf("TableName method not found at index 1, error: %v", e))
			return
		}
	}()
	reflectValue := reflect.Indirect(reflect.ValueOf(ts))
	var tableName string
	switch reflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		if reflectValue.Len() > 0 {
			v := reflectValue.Index(0)
			rm := v.MethodByName("TableName")
			v2 := rm.Call(make([]reflect.Value, 0))[0]
			tableName = fmt.Sprintf("%v", v2)
		}
		err = store.db.Table(tableName).CreateInBatches(ts, reflectValue.Len()).Error
		if err != nil {
			return errorx.Wrap(err, "batches create error")
		}
	default:
		return errorx.New("please enter a slice or array containing the TableName method")
	}
	return
}

type RegisterFunc func() error

func IsUniqueError(err error) (ok bool) {
	return strings.Contains(err.Error(), "UNIQUE constraint failed:")
}
