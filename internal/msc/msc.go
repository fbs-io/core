/*
 * @Author: reel
 * @Date: 2023-05-17 22:49:53
 * @LastEditors: reel
 * @LastEditTime: 2023-06-05 07:57:49
 * @Description: 后台管理中心
 */
package msc

import (
    "net/http"

    "github.com/fbs-io/core/internal/config"
    "github.com/fbs-io/core/internal/msc/views"
    "github.com/fbs-io/core/store/cache"

    "github.com/gin-gonic/gin"
)

type handler struct {
    config *config.Config
    cache  cache.Store
}

func Init(engine *gin.Engine, conf *config.Config, cache cache.Store) {
    m := &handler{
        config: conf,
        cache:  cache,
    }

    // 加载中间件
    engine.Use(m.log())
    engine.Use(m.cors())
    engine.Use(m.validCookie())

    // 加载静态资源
    engine.GET("/static/*filepath", func(ctx *gin.Context) {
        staticSrv := http.FileServer(http.FS(views.Static))
        staticSrv.ServeHTTP(ctx.Writer, ctx.Request)
    })

    engine.GET("/", m.index())
    ajax := engine.Group("ajax")

    // 登陆相关
    {
        ajax.POST("/login", m.login())
        // ajax.POST("/resetpwd", m.resetpwd())

    }
    // ajax := engine.Group("ajax")
    // if conf.IsLoad {
    // 	{
    // 		// 初始化配置相关
    // 		// ajax.GET("/defaultconfig", m.defaultConfig())
    // 	}

    // }

}
