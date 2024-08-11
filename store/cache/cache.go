/*
 * @Author: reel
 * @Date: 2023-05-16 07:50:35
 * @LastEditors: reel
 * @LastEditTime: 2024-05-11 12:21:19
 * @Description: 请填写简介
 */
package cache

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/dsn"

	"github.com/tidwall/buntdb"
)

const (
	cacheStatusKey   = "cache::status"
	cacheStatusValue = "Y"
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
	SetIndex(indexName, index string, fs ...func(a string, b string) bool) error
	GetByIndex(indexName string, iterator func(key string, value string) bool) error
	GetWithObj(key string, obj interface{}) error
	SetWithObj(key string, obj interface{}, funcs ...OptFunc) error
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
		value = strings.TrimSpace(value)
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
	return nil
}

// 通过插入和写入操作, 测试缓存功能是否正常
func (c *cacheStore) Status() int8 {
	err := c.Set(cacheStatusKey, cacheStatusValue)
	if err != nil {
		return 0
	}

	err = c.db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(cacheStatusKey)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return 0
	}
	return 1
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

// 设置索引
func (c *cacheStore) SetIndex(indexName, index string, fs ...func(a string, b string) bool) error {
	return c.db.CreateIndex(indexName, index, fs...)
}

// 通过索引查询
//
// 结果集放在iterator函数内实现, 如: func(k,v)bool
func (c *cacheStore) GetByIndex(indexName string, iterator func(key string, value string) bool) error {
	return c.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(indexName, iterator)
	})

}

// 设置对象缓存, 基于set方法增加json序列化方法
func (c *cacheStore) SetWithObj(key string, obj interface{}, funcs ...OptFunc) error {
	vb, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return c.Set(key, string(vb), funcs...)
}

// 获取缓存对象, 需要传入对象指针
func (c *cacheStore) GetWithObj(key string, obj interface{}) error {
	vbs := c.Get(key)
	if vbs == "" {
		return errors.New("没有缓存数据")
	}
	return json.Unmarshal([]byte(vbs), obj)
}
