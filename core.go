/*
 * @Author: reel
 * @Date: 2023-05-11 23:25:29
 * @LastEditors: reel
 * @LastEditTime: 2023-06-06 06:03:24
 * @Description: 管理核心组件的启动和运行
 */
package core

import (
    "fmt"

    "github.com/fbs-io/core/cron"
    "github.com/fbs-io/core/internal/config"
    "github.com/fbs-io/core/internal/msc"
    "github.com/fbs-io/core/internal/pem"
    "github.com/fbs-io/core/pkg/env"
    "github.com/fbs-io/core/pkg/errorx"
    "github.com/fbs-io/core/pkg/mux"
    "github.com/fbs-io/core/service"
    "github.com/fbs-io/core/store/cache"
    "github.com/fbs-io/core/store/rdb"

    "github.com/gin-gonic/gin"
)

type core struct {
    config   *config.Config
    msc      mux.Mux // 用于管理整个服务
    ams      mux.Mux // 用于应用管理
    services []service.Service
    cron     cron.Cron
    cache    cache.Store
    rdb      rdb.Store
}

var _ Core = (*core)(nil)

type Core interface {
    // 私有方法, 防止和其他借口冲突
    coreP()

    // 封装了服务启动和关闭
    // 方便快速启动
    Run()

    // 关闭整个服务
    Shutdown()

    // gin的engine, 用于原生gin方法
    // 可以更灵活的实现开发
    Engine() *gin.Engine

    // 缓存
    Cache() cache.Store

    // 关系数据库
    RDB() rdb.Store
}

func (c *core) coreP() {}

func New() (Core, error) {
    env.Init()

    gin.SetMode(env.Active().Mode())
    dms, err := mux.New(
        mux.SetHost(env.Active().MscAddr()),
        mux.SetName("DMSC"),
    )
    if err != nil {
        return nil, errorx.Wrap(err, "初始化后台管理服务发生错误")
    }
    ams, err := mux.New(
        mux.SetHost(":80"),
        mux.SetName("AMS"),
    )
    if err != nil {
        return nil, errorx.Wrap(err, "初始化应用管理服务发生错误")
    }
    c := &core{
        msc:    dms,
        ams:    ams,
        rdb:    rdb.New(),
        cron:   cron.New(),
        cache:  cache.New(),
        config: &config.Config{},
    }

    // 配置中心和其他服务分开启动和关闭

    msc.Init(c.msc.Engine(), c.config, c.cache)

    if _, err := pem.GetPems(); err != nil {
        c.msc.Engine().POST("/ajax/install", c.installHandler())
    }
    err = c.msc.Start()
    if err != nil {
        return c, errorx.Wrap(err, fmt.Sprintf("启动%s失败", c.msc.Name()))
    }

    return c, nil
}
