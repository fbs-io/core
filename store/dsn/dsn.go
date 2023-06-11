/*
 * @Author: reel
 * @Date: 2023-05-16 20:02:15
 * @LastEditors: reel
 * @LastEditTime: 2023-06-11 15:48:13
 * @Description: 配置数据库，缓存的链接, 支持本地缓存和本地数据库
 */
package dsn

import (
    "fmt"
    "path"

    "github.com/fbs-io/core/logx"
    "github.com/fbs-io/core/pkg/env"
)

const (
    DSN_TYPE_LOCAL   = "local"
    DSN_TYPE_SQLITE  = "sqlite"
    DSN_TYPE_REDIS   = "redis"
    DSN_TYPE_MYSQL   = "mysql"
    DSN_TYPE_PGSQL   = "postgres"
    MYSQL_DB_DSN_KYE = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
    PGSQL_DB_DSN_KYE = "user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai"
)

type Dsn struct {
    link      string
    Type      string
    Path      string
    Name      string
    Host      string
    Port      string
    User      string
    Pwd       string
    LocalName string
    Log       logx.Logger
}

// TODO: 增加 mysql postgre redis 等
func (d *Dsn) Link() string {
    switch d.Type {
    case DSN_TYPE_LOCAL, DSN_TYPE_SQLITE:
        d.link = path.Join(d.Path, d.Name)
    case "redis":
    case "postgres":
        d.link = fmt.Sprintf(
            PGSQL_DB_DSN_KYE,
            d.User,
            d.Pwd,
            d.Host,
            d.Port,
            d.Name,
        )
    case "mysql":
        d.link = fmt.Sprintf(
            MYSQL_DB_DSN_KYE,
            d.User,
            d.Pwd,
            d.Host,
            d.Port,
            d.Name,
        )
    }
    return d.link
}

type DsnFunc func(Dsn *Dsn)

// 直接设置数据库链接,而不需要再设置主机端口等
// TODO: 配置 redis, mysql 等数据库支持
func SetLink(link string) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.link = link
    }
}

// 设置本地 db 路径
func SetPath(path string) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.Path = path
    }
}

// 设置 db 名称
func SetName(name string) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.Name = name
    }
}

// 设置 db 名称
func SetLocalName(name string) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.LocalName = name
    }
}

func SetHost(host string) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.Host = host
    }
}

func SetPort(port string) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.Port = port
    }
}

func SetUser(user string) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.User = user
    }
}

func SetPwd(pwd string) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.Pwd = pwd
    }
}
func SetLog(log logx.Logger) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.Log = log
    }
}

func SetType(t string) DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.Type = t
    }
}

func NewCacheDsn() *Dsn {
    return &Dsn{
        Type: DSN_TYPE_LOCAL,
        Path: path.Join(env.Active().DataPath(), "db/cache/"),
        Name: env.Active().AppName() + ".cache",
        // Log: logx.New(logx.SetLogPath("cache")),
    }
}

func NewDBDsn() *Dsn {
    return &Dsn{
        Type: DSN_TYPE_LOCAL,
        Path: "db/rdb/",
        Name: env.Active().AppName() + ".db",
    }
}

// 配置默认端口为 5432
func PGDefaultDsn() DsnFunc {
    return func(Dsn *Dsn) {
        Dsn.Type = DSN_TYPE_PGSQL
        Dsn.Port = "5432"
        Dsn.Name = "postgres"
    }
}

func NewPGDsn() *Dsn {
    return &Dsn{
        Type: DSN_TYPE_PGSQL,
        Port: "5432",
        Name: "postgres",
    }
}
