/*
 * @Author: reel
 * @Date: 2023-05-16 22:16:53
 * @LastEditors: reel
 * @LastEditTime: 2023-08-15 23:25:13
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
	TABLE_INIT_ALL = "all"

// TABLE_SYS_APIDOC = "sys_apidoc"
// TABLE_SYS_MEANS  = "sys_means"
)

type rdbStore struct {
	db          *gorm.DB
	dial        gorm.Dialector
	statTab     Tabler
	tablers     map[string]Tabler
	migrateList []func() error
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
	Register(t Tabler, fs ...RegisterFunc)
	BuildQueryWithParams(params reflect.Value) *gorm.DB
}

var _ Store = (*rdbStore)(nil)

var rdb = &rdbStore{
	db:          &gorm.DB{},
	tablers:     make(map[string]Tabler, 1000),
	migrateList: make([]func() error, 0, 100),
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
		// Logger: logger.Default.LogMode(env.Active().GormLogLevel()),
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

	return store.autoMigrate()
}

// 通过查询获取当前db状态
// 如果没有表, 默认当前不可用
func (store *rdbStore) Status() int8 {
	if store.statTab == nil {
		return 0
	}
	err := store.db.FirstOrCreate(&(store.statTab)).Error
	if err != nil {
		return 0
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
func (store *rdbStore) Register(t Tabler, fs ...RegisterFunc) {
	store.tablers[t.TableName()] = t

	store.migrateList = append(store.migrateList, func() (err error) {
		// 根据实际使用, 系统资源表将在最后被加载
		store.statTab = t

		// 从命令行重置所有表或部分表
		if strings.Contains(env.Active().DBInit(), t.TableName()) ||
			env.Active().DBInit() == TABLE_INIT_ALL {
			err = store.db.Migrator().DropTable(t)
			if err != nil {
				return
			}
		}
		if !store.db.Migrator().HasTable(t) {
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
		return
	})
}

func (r *rdbStore) autoMigrate() (err error) {
	for _, fs := range r.migrateList {
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
func (r *rdbStore) CreateInBatches(ts interface{}) (err error) {
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
		err = r.db.Table(tableName).CreateInBatches(ts, reflectValue.Len()).Error
		if err != nil {
			return errorx.Wrap(err, "batches create error")
		}
	default:
		return errorx.New("please enter a slice or array containing the TableName method")
	}
	return
}

type RegisterFunc func() error
