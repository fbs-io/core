/*
 * @Author: reel
 * @Date: 2023-05-21 20:49:21
 * @LastEditors: reel
 * @LastEditTime: 2023-05-21 20:56:18
 * @Description: 用于管理在core上的各个模块, 如cache, db, handle等
 */
package core

import (
	"core/store/cache"
	"core/store/rdb"

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
