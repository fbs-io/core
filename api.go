/*
 * @Author: reel
 * @Date: 2023-05-21 20:49:21
 * @LastEditors: reel
 * @LastEditTime: 2024-03-27 04:47:04
 * @Description: 用于管理在core上的各个模块, 如cache, db, handle等
 */
package core

import (
	"github.com/fbs-io/core/internal/config"
	"github.com/fbs-io/core/pkg/mux"
	"github.com/fbs-io/core/session"
	"github.com/fbs-io/core/store/cache"
	"github.com/fbs-io/core/store/rdb"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

func (c *core) Engine() *gin.Engine {
	return c.ams.Engine()
}

func (c *core) Cache() cache.Store {
	return c.cache
}

func (c *core) RDB() rdb.Store {
	return c.rdb
}

func (c *core) Session() session.Session {
	return c.session
}

// 限流器
func (c *core) Limiter() *rate.Limiter {
	return c.limiter
}

// 配置
func (c *core) Config() *config.Config {
	return c.config
}

// 增加执行的程序
func (c *core) AddStartJobList(startjobs ...mux.StartJobFunc) {
	c.ams.AddStartJobList(startjobs...)
}
