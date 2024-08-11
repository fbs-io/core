/*
 * @Author: reel
 * @Date: 2023-05-16 20:17:56
 * @LastEditors: reel
 * @LastEditTime: 2024-07-07 20:48:18
 * @Description: 系统配置相关操作
 */
package core

import (
	"encoding/json"
	"fmt"

	"path"

	"github.com/fbs-io/core/internal/pem"
	"github.com/fbs-io/core/logx"
	"github.com/fbs-io/core/pkg/encrypt"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/service"
	"github.com/fbs-io/core/store/dsn"

	"github.com/gin-gonic/gin"
)

func (c *core) install() (err error) {
	defer func() {
		if err == nil {
			c.config.IsLoad = true
		}
	}()

	err = logx.Init(logx.SetDataPath(c.config.DataPath))
	if err != nil {
		return errorx.Wrap(err, "日志服务初始化失败")
	}
	logx.Sys.Info("完成日志配置")

	// 通过service统一管理服务的启动和停止
	service.Append(c.cron)
	logx.Sys.Info("完成定时服务配置")

	// 缓存
	err = c.cache.SetConfig(
		dsn.SetType(c.config.CacheType),
		dsn.SetPath(path.Join(c.config.DataPath, "cache")),
		dsn.SetHost(c.config.CacheHost),
		dsn.SetPort(c.config.CachePort),
		dsn.SetName(c.config.CacheName),
		dsn.SetUser(c.config.CacheUser),
		dsn.SetPwd(c.config.CachePwd),
	)
	if err != nil {
		return errorx.Wrap(err, "缓存服务初始化失败")
	}
	// 定时对缓存进行清理
	c.cron.AddJob(func() { c.cache.Shrink() }, "cache 空间释放", 3600)

	service.Append(c.cache)
	logx.Sys.Info("完成缓存配置")

	// db初始化
	err = c.rdb.SetConfig(
		dsn.SetType(c.config.DbType),
		dsn.SetPath(path.Join(c.config.DataPath, "db")),
		dsn.SetHost(c.config.DBHost),
		dsn.SetPort(c.config.DBPort),
		dsn.SetName(c.config.DBName),
		dsn.SetUser(c.config.DBUser),
		dsn.SetPwd(c.config.DBPwd),
	)
	if err != nil {
		return errorx.Wrap(err, "DB服务初始化失败")
	}

	// 把资源map转成list, 初始化写入资源列表
	s := &Resources{}
	c.rdb.Register(s, func() error { return c.rdb.DB().Table(s.TableName()).CreateInBatches(resources, len(resources)).Error })

	service.Append(c.rdb)
	logx.Sys.Info("完成数据库配置")

	// 重设web端口
	c.ams.SetAddr(fmt.Sprintf(":%s", c.config.Port))

	service.Append(c.ams)
	logx.Sys.Info("完成APP服务配置")

	return
}

// 安装
func (c *core) installHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			err  error
			data []byte
			str  string
		)
		defer func() {
			e := recover()
			if e != nil {
				logx.Sys.Error("系统初始化发生错误", logx.F("err", fmt.Sprintf("%v", e)))
				ctx.JSON(200, errno.ERRNO_SYSTEM.ToMapWithError(errorx.Errorf("%v", e)))
			}
		}()
		if c.config.IsLoad {
			ctx.JSON(200, errno.ERRNO_IS_INSTALL.ToMap())
			return
		}
		pem.RemovePem()
		err = ctx.ShouldBindJSON(&(c.config))
		if err != nil {
			ctx.JSON(200, errno.ERRNO_PARAMS_BIND.ToMapWithError(err))
			return
		}
		data, err = json.Marshal(c.config)
		if err != nil {
			ctx.JSON(200, errno.ERRNO_PARAMS_INVALID.ToMapWithError(err))
			return
		}
		str, err = encrypt.InternalEncode(data)
		if err != nil {
			ctx.JSON(200, errno.ERRNO_PARAMS_INVALID.ToMapWithError(err))
			return
		}
		err = c.install()
		if err != nil {
			ctx.JSON(200, errno.ERRNO_INIT.ToMapWithError(err))
			return
		}
		if err := service.Start(); err != nil {
			c.config.IsLoad = false
			service.Stop()
			ctx.JSON(200, errno.ERRNO_PARAMS_INVALID.ToMapWithError(err))
			return
		}
		err = pem.UpdatePems(str)
		if err != nil {
			ctx.JSON(200, errno.ERRNO_SYSTEM.ToMapWithError(err))
			return
		}
		ctx.JSON(200, errno.ERRNO_OK.ToMap())
	}
}
