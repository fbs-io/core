/*
 * @Author: reel
 * @Date: 2023-05-16 22:16:53
 * @LastEditors: reel
 * @LastEditTime: 2023-06-20 06:48:01
 * @Description: 关系数据库配置
 */
package rdb

import (
    "fmt"
    "os"
    "reflect"
    "strings"
    "time"

    "github.com/fbs-io/core/pkg/env"
    "github.com/fbs-io/core/pkg/errorx"
    "github.com/fbs-io/core/store/dsn"

    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
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
    Name() string
    Start() error
    Stop() error
    Status() int8
    Register(t Tabler, fs ...RegisterFunc)
    SetConfig(fs ...dsn.DsnFunc) error
    DB() *gorm.DB
    // Create(t Tabler) (err error)
    // Delete(t Tabler, bfs ...BuildFunc) (err error)
    // CreateInBatches(ts interface{}) (err error)
    // Query(t Tabler, bfs ...BuildFunc) (err error)
    // Updates(t Tabler, bfs ...BuildFunc) (err error)
    // Queries(ts interface{}, bfs ...BuildFunc) (err error)
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

func (rdb *rdbStore) rdbP()        {}
func (rdb *rdbStore) Name() string { return "RDB" }
func (rdb *rdbStore) DB() *gorm.DB { return rdb.db }

func (rdb *rdbStore) Start() (err error) {
    if rdb.dial == nil {
        return errorx.New("没有可用的DSN")
    }
    rdb.db, err = gorm.Open(rdb.dial, &gorm.Config{
        Logger: logger.Default.LogMode(env.Active().GormLogLevel()),
    })
    if err != nil {
        return errorx.Wrap(err, "数据库链接失败")
    }
    // 配置链接时间等
    sqlDB, err := rdb.db.DB()
    if err != nil {
        return err
    }

    // TODO: 自定义配置连接池大小
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return rdb.autoMigrate()
}

// 通过查询获取当前db状态
// 如果没有表, 默认当前不可用
func (rdb *rdbStore) Status() int8 {
    if rdb.statTab == nil {
        return 0
    }
    err := rdb.db.FirstOrCreate(&(rdb.statTab)).Error
    if err != nil {
        return 0
    }
    return 1
}

func (r *rdbStore) Stop() error {
    sqlDB, err := r.db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

func (r *rdbStore) SetConfig(optfs ...dsn.DsnFunc) (err error) {
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
        rdb.dial = sqlite.Open(link)

    case dsn.DSN_TYPE_PGSQL:
        rdb.dial = postgres.Open(link)
    case dsn.DSN_TYPE_MYSQL:
        rdb.dial = mysql.Open(link)
    }
    return nil
}

// 迁移表
// 如果表结构发生变化, 需要启动时设置 dbinit=tableaName 参数, 删除旧表并创建新表
// 在表创建后, 可以执行一些自定义的方法, 主要用于初始化数据写入
func (r *rdbStore) Register(t Tabler, fs ...RegisterFunc) {
    r.tablers[t.TableName()] = t
    if r.statTab == nil {
        r.statTab = t
    }
    r.migrateList = append(r.migrateList, func() (err error) {
        // 从命令行重置所有表或部分表
        if strings.Contains(env.Active().DBInit(), t.TableName()) ||
            env.Active().DBInit() == TABLE_INIT_ALL {
            err = r.db.Migrator().DropTable(t)
            if err != nil {
                return
            }
        }
        if !r.db.Migrator().HasTable(t) {
            err = r.db.AutoMigrate(t)
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
