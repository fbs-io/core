/*
 * @Author: reel
 * @Date: 2023-05-11 23:25:29
 * @LastEditors: reel
 * @LastEditTime: 2024-03-27 04:41:41
 * @Description: 管理核心组件的启动和运行
 */
package core

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fbs-io/core/cron"
	"github.com/fbs-io/core/internal/config"
	"github.com/fbs-io/core/internal/msc"
	"github.com/fbs-io/core/internal/pem"
	"github.com/fbs-io/core/pkg/env"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/pkg/mux"
	"github.com/fbs-io/core/session"
	"github.com/fbs-io/core/store/cache"
	"github.com/fbs-io/core/store/rdb"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

type core struct {
	config  *config.Config
	msc     mux.Mux // 用于管理整个服务
	ams     mux.Mux // 用于应用管理
	cron    cron.Cron
	cache   cache.Store
	rdb     rdb.Store
	session session.Session
	limiter *rate.Limiter
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

	// 基于gin.Engine封装, 用于快速开发

	Group(elativePath string, handlers ...HandlerFunc) RouterGroup

	// 缓存
	Cache() cache.Store

	// 关系数据库
	RDB() rdb.Store

	// session
	Session() session.Session

	// 限流器
	Limiter() *rate.Limiter

	// 配置
	Config() *config.Config

	// 添加启动时的一些作业
	AddStartJobList(...mux.StartJobFunc)
}

func (c *core) coreP() {}

func New(funcs ...FuncCores) (Core, error) {
	var opt = &options{
		limitSize:   200 * runtime.NumCPU(),
		limitNumber: 200 * runtime.NumCPU(),
	}
	for _, fs := range funcs {
		fs(opt)
	}

	env.Init()
	gin.SetMode(env.Active().Mode())
	dms, err := mux.New(
		mux.SetHost(env.Active().MscAddr()),
		mux.SetName("MSC"),
		mux.SetTimeout(180),
		mux.SetMaxReadTime(180*time.Second),
		mux.SetMaxWriteTIme(180*time.Second),
	)
	if err != nil {
		return nil, errorx.Wrap(err, "初始化后台管理服务发生错误")
	}
	ams, err := mux.New(
		mux.SetHost(":80"),
		mux.SetName("AMS"),
		mux.SetTimeout(180),
		mux.SetMaxReadTime(180*time.Second),
		mux.SetMaxWriteTIme(180*time.Second),
	)
	if err != nil {
		return nil, errorx.Wrap(err, "初始化应用管理服务发生错误")
	}
	c := &core{
		msc:     dms,
		ams:     ams,
		rdb:     rdb.New(),
		cron:    cron.New(),
		cache:   cache.New(),
		config:  &config.Config{},
		limiter: rate.NewLimiter(rate.Limit(opt.limitNumber), opt.limitSize),
	}
	c.session = session.New(session.Store(c.cache), session.Prefix("app::session"))
	// 配置中心和其他服务分开启动和关闭

	msc.Init(c.msc.Engine(), c.config, c.cache, c.cron)

	if _, err := pem.GetPems(); err != nil {
		c.msc.Engine().POST("/msc/ajax/install", c.installHandler())
	}
	err = c.msc.Start()
	if err != nil {
		return c, errorx.Wrap(err, fmt.Sprintf("启动%s失败", c.msc.Name()))
	}

	return c, nil
}
