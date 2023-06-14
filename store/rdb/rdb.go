/*
 * @Author: reel
 * @Date: 2023-05-16 22:16:53
 * @LastEditors: reel
 * @LastEditTime: 2023-06-14 20:30:45
 * @Description: 关系数据库配置
 */
package rdb

import (
    "os"
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
    TABLE_INIT_ALL   = "all"
    TABLE_SYS_APIDOC = "sys_apidoc"
    TABLE_SYS_MEANS  = "sys_means"
)

type rdbStore struct {
    db          *gorm.DB
    dial        gorm.Dialector
    migrateList []func() error
}

type Store interface {
    rdbP()
    Name() string
    Start() error
    Stop() error
    Status() int8
    Register(t Tabler, fs ...func())
    SetConfig(fs ...dsn.DsnFunc) error
    AutoMigrate() (err error)
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
    migrateList: make([]func() error, 0, 100),
}

func New() (s Store) {

    return rdb
}

type ping struct {
    Content string
}

func (f *ping) TableName() string {
    return "sys_ping"
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

    // 用于测试数据库是否正常连接
    // TODO:更换为系统常用表
    var p = &ping{
        Content: "用于测试数据库链接情况, 请勿删除",
    }
    if !rdb.db.Migrator().HasTable(p) {
        err := rdb.db.Migrator().CreateTable(p)
        if err != nil {
            return err
        }
    }

    return nil
}
func (rdb *rdbStore) Status() int8 {
    err := rdb.db.FirstOrCreate(&ping{"用于测试数据库链接情况, 请勿删除"}).Error
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
func (r *rdbStore) Register(t Tabler, fs ...func()) {
    r.migrateList = append(r.migrateList, func() (err error) {
        if strings.Contains(env.Active().DBInit(), t.TableName()) ||
            env.Active().DBInit() == TABLE_INIT_ALL ||
            // api文档表和资源表在开发模式下, 每次均重置
            (env.Active().Value() == env.ENV_MODE_DEV && t.TableName() == TABLE_SYS_APIDOC) ||
            (env.Active().Value() == env.ENV_MODE_DEV && t.TableName() == TABLE_SYS_MEANS) {
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
                f()
            }
        }
        return
    })
}

func (r *rdbStore) AutoMigrate() (err error) {
    for _, fs := range r.migrateList {
        e := fs()
        if e != nil {
            err = errorx.Wrap(e, "迁移失败")
        }
    }
    return
}

type Tabler interface {
    TableName() string
    // GetID() uint
}
