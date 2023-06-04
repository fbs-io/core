/*
 * @Author: reel
 * @Date: 2023-05-16 07:50:35
 * @LastEditors: reel
 * @LastEditTime: 2023-06-04 11:27:08
 * @Description: 请填写简介
 */
package cache

import (
    "core/pkg/errorx"
    "core/store/dsn"
    "os"
    "time"

    "github.com/tidwall/buntdb"
)

// TODO: 增加日志
type cacheStore struct {
    db  *buntdb.DB
    dsn *dsn.Dsn
    // log logx.Logger
}
type Store interface {
    cacheP()
    Name() string
    Status() int8
    Start() error
    Stop() error
    Shrink() error
    Get(string) string
    Del(key string) error
    Set(key, value string, funcs ...OptFunc) (err error)
    SetConfig(funcs ...dsn.DsnFunc) (err error)
}

var _ Store = (*cacheStore)(nil)

var cache = &cacheStore{}

func New() Store {
    return cache
}

func (c *cacheStore) cacheP()      {}
func (c *cacheStore) Name() string { return "Cache" }

func (c *cacheStore) Get(key string) (value string) {
    var err error
    err = c.db.View(func(tx *buntdb.Tx) error {
        value, err = tx.Get(key)
        if err != nil {
            return err
        }
        return err
    })

    return
}

func (c *cacheStore) Set(key, value string, funcs ...OptFunc) (err error) {
    var opt = &buntdb.SetOptions{}
    for _, f := range funcs {
        f(opt)
    }
    err = c.db.Update(func(tx *buntdb.Tx) error {
        _, _, err = tx.Set(key, value, opt)
        return err
    })
    return
}

func (c *cacheStore) Del(key string) (err error) {
    err = c.db.Update(func(tx *buntdb.Tx) error {
        _, err = tx.Delete(key)
        return err
    })
    return
}

func (c *cacheStore) Shrink() error {
    return c.db.Shrink()
}

type OptFunc func(*buntdb.SetOptions)

// 设置过期时间
func SetTTL(ttl time.Duration) OptFunc {
    return func(so *buntdb.SetOptions) {
        so.Expires = true
        so.TTL = ttl * time.Second
    }
}

func (c *cacheStore) Start() error {
    if c.dsn == nil {
        return errorx.New("没有可用的DSN")
    }
    db, err := buntdb.Open(c.dsn.Link())
    if err != nil {
        return errorx.Wrap(err, "cache 启动失败")
    }
    db.Shrink()
    c.db = db
    c.Set("cache::status", "Y")
    return nil
}
func (c *cacheStore) Status() int8 {
    if c.db == nil {
        return -1
    }
    s := c.Get("cache::status")
    if s == "Y" {
        return 1
    }
    return 0
}

func (c *cacheStore) Stop() error {
    return c.db.Close()
}

func (c *cacheStore) SetConfig(dsnFuns ...dsn.DsnFunc) (err error) {
    cacheDsn := dsn.NewCacheDsn()

    for _, optf := range dsnFuns {
        optf(cacheDsn)
    }

    if cacheDsn.Link() == "" {
        return errorx.New("dsn 为空, 请检查配置及 DB 类型是否正确")
    }
    if cacheDsn.Type == dsn.DSN_TYPE_LOCAL {
        _, err = os.Stat(cacheDsn.Path)
        if err != nil {
            err = os.MkdirAll(cacheDsn.Path, 0766)
            if err != nil {
                return
            }
        }
    } else {
        return errorx.New("缓存链接类型不正确")
    }
    c.dsn = cacheDsn
    return nil
}
