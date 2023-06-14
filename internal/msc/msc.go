/*
 * @Author: reel
 * @Date: 2023-05-17 22:49:53
 * @LastEditors: reel
 * @LastEditTime: 2023-06-14 20:26:15
 * @Description: 后台管理中心
 */
package msc

import (
    "fmt"
    "net/http"

    "github.com/fbs-io/core/cron"
    "github.com/fbs-io/core/internal/config"
    "github.com/fbs-io/core/internal/msc/views"
    "github.com/fbs-io/core/logx"
    "github.com/fbs-io/core/service"
    "github.com/fbs-io/core/session"
    "github.com/fbs-io/core/store/cache"

    "github.com/gin-gonic/gin"
)

type handler struct {
    config    *config.Config
    cache     cache.Store
    session   session.Session
    cron      cron.Cron
    procinfos []map[string]interface{}
    sysinfo   map[string]interface{}
    srvStatus []map[string]interface{}
}

func Init(engine *gin.Engine, conf *config.Config, cache cache.Store, cron cron.Cron) {
    m := &handler{
        config:    conf,
        cache:     cache,
        session:   session.New(session.Store(cache)),
        cron:      cron,
        procinfos: make([]map[string]interface{}, 0, 20),
        sysinfo:   make(map[string]interface{}, 0),
        srvStatus: make([]map[string]interface{}, 0),
    }

    m.cron.AddJob(
        func() {
            info := m.getProcessInfo()
            data := []map[string]interface{}{}
            for _, i := range info {
                data = append(data, map[string]interface{}{
                    "pid":        i.PID,
                    "pname":      i.PName,
                    "cpupercent": fmt.Sprintf("%.2f", i.CpuPercent) + "%",
                    "meminfo":    fmt.Sprintf("%d MB", i.MemInfo),
                    "io":         fmt.Sprintf("%d ", i.IO),
                })
                m.procinfos = data
            }
            logx.Sys.Info("进程信息", logx.F("details", data))
        },
        "系统进程Top20",
        20,
    )
    m.cron.AddJob(
        func() {
            logx.Sys.Info("系统信息", logx.Details(m.getSysInfo()))
            m.srvStatus = service.Status()
            logx.Sys.Info("服务状态", logx.F("details", m.srvStatus))
        },
        "当前系统资源使用率及服务状态查询",
        2,
    )

    // 加载中间件
    engine.Use(m.log())
    engine.Use(m.cors())
    engine.Use(m.signature())

    // 加载静态资源
    engine.GET("/static/*filepath", func(ctx *gin.Context) {
        staticSrv := http.FileServer(http.FS(views.Static))
        staticSrv.ServeHTTP(ctx.Writer, ctx.Request)
    })

    engine.GET("/", m.index())
    ajax := engine.Group("ajax")

    // 登陆相关
    {
        ajax.POST("login", m.login())

    }

    {
        // 服务器信息
        ajax.GET("/sysinfo", m.sysInfo())
        ajax.GET("/hostinfo", m.hostInfo())
        ajax.GET("/srvstatus", m.getSrvStatus())
        ajax.POST("/srvrestart", m.setSrvRestart())
    }

}
