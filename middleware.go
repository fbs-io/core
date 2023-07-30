/*
 * @Author: reel
 * @Date: 2023-07-19 00:08:08
 * @LastEditors: reel
 * @LastEditTime: 2023-07-23 22:44:10
 * @Description: 常用的中间件
 */
package core

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/fbs-io/core/logx"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/gin-gonic/gin"
)

const (
	COOKIE_SID = "ASID"
)

// 跨域处理中间件
func CorsMiddleware(c Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			ctx.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Content-type", "application/json")
		}

		//允许类型校验
		if method == "OPTIONS" {
			ctx.JSON(200, errno.ERRNO_OK.ToMap())
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		ctx.Next()
	}
}

// 日志中间件
func LogMiddleware(c Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		startTime := time.Now()
		ctx.Next()
		if strings.Contains(ctx.Request.RequestURI, "/static/") {
			return
		}
		endTime := time.Now()
		logx.Sys.Debug("http请求", logx.F("status", ctx.Writer.Status()),
			logx.F("diff_time", fmt.Sprintf("%d ns", endTime.Sub(startTime))),
			logx.F("client_ip", ctx.ClientIP()),
			logx.F("req_method", ctx.Request.Method),
			logx.F("req_url", ctx.Request.RequestURI),
		)
	}
}

var allowResource = make(map[string]bool, 100)

// 请使用 method:path 的方式定义资源
// 比如 POST:/ajax/login
func AddAllowResource(resoures ...string) {

	for _, resoure := range resoures {
		allowResource[resoure] = true
	}
}

func resource(ctx *gin.Context) string {
	return fmt.Sprintf("%s:%s", ctx.Request.Method, ctx.FullPath())
}

// 校验签名中间件
// 如果没有登陆, 则会给一个默认的签名
func SignatureMiddleware(c Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if allowResource[resource(ctx)] {
			ctx.Next()
			return
		}
		sessionKey, sessiionValue, err := c.Session().GetWithToken(ctx.Request)
		fmt.Println(sessionKey, sessiionValue, err)
		if err != nil {
			ctx.JSON(200, errno.ERRNO_AUTH_NOT_LOGIN.ToMapWithError(err))
			ctx.Abort()
			return
		}
		// 更新session过期时间
		// c.Session().SetWithToken(sessionKey, sessiionValue)
		c.Session().SetWithCookie(ctx.Writer, sessionKey, sessiionValue)
		// 无论是否有获取到cookie, 均需要重新设置cookie
		ctx.Next()

	}
}

// 参数自动生成中间件
//
// 会生成参数结构体以及gorm.DB
//
// 同时根据约束, 自动完成参数校验
func ParamsMiddleware(c Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res := resource(ctx)
		rt := requestParams[res]
		if rt == nil {
			ctx.Next()
			return
		}
		rv := reflect.New(rt)
		params := rv.Interface()
		var err error
		switch ctx.ContentType() {
		case formContent:
			err = ctx.ShouldBindQuery(&params)
		case jsonContent:
			err = ctx.ShouldBindJSON(&params)
		}
		if err != nil {
			ctx.JSON(200, errno.ERRNO_PARAMS_BIND.ToMapWithError(errorx.Wrap(err, "校验参数发生错误")))
			// ctx.JSON(errno.ERRNO_PARAMS_BIND.ToMapWithError(errorx.Wrap(err, "校验参数发生错误")))
			ctx.Abort()
			return
		}
		// ctx.CtxSet()
		ctx.Set(CTX_PARAMS, params)
		tx := c.RDB().BuildQueryWithParams(rv)
		ctx.Set(CTX_TX, tx)
	}
}

// 限流器
func LimiterMiddleware(c Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !c.Limiter().Allow() {
			ctx.JSON(200, errno.ERRNO_TOO_MANY_REQUESTS.ToMap())
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
