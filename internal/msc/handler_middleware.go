/*
 * @Author: reel
 * @Date: 2023-05-28 18:06:12
 * @LastEditors: reel
 * @LastEditTime: 2024-03-27 04:46:34
 * @Description: 中间件
 */
package msc

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fbs-io/core/logx"
	"github.com/fbs-io/core/pkg/errno"

	"github.com/gin-gonic/gin"
)

const (
	COOKIE_SID = "SID"
)

// 跨域处理中间件
func (m *handler) cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin") //请求头部
		if origin != "" && !strings.Contains(ctx.FullPath(), "website") {
			//接收客户端发送的origin （重要！）
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			//服务器支持的所有跨域请求的方法
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Content-Length, X-CSRF-TOKEN, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, X-CSRF-TOKEN, SID")
			//设置缓存时间
			ctx.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			ctx.Header("Access-Control-Allow-Credentials", "true")
			// ctx.Header("Content-type", ctx.ContentType())
		}

		//允许类型校验
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, gin.H{"code": 0})
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
func (m *handler) log() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		startTime := time.Now()
		ctx.Next()
		if strings.Contains(ctx.Request.RequestURI, "/mscui/") {
			return
		}
		endTime := time.Now()
		if m.config.IsLoad {
			logx.Sys.Debug("http请求", logx.F("status", ctx.Writer.Status()),
				logx.F("diff_time", fmt.Sprintf("%d ns", endTime.Sub(startTime))),
				logx.F("client_ip", ctx.ClientIP()),
				logx.F("req_method", ctx.Request.Method),
				logx.F("req_url", ctx.Request.RequestURI),
			)
		}
	}
}

var allowPath = map[string]bool{
	"/msc/ajax/login":   true,
	"/msc/ajax/install": true,
	"/":                 true,
	"":                  true,
}

const (
	initConfig = `const APP_CONFIG = { APP_NAME: "FBS Manager System Center",API_URL:  "", APP_INIT: true}`
)

// 校验签名
// 如果没有登陆, 则会给一个默认的签名
func (m *handler) signature() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.RequestURI == `/website/config.js` && !m.config.IsLoad {
			ctx.Header("Content-Type", "application/javascript")
			ctx.String(200, initConfig)
			ctx.Abort()
			return
		}

		if allowPath[ctx.FullPath()] {
			ctx.Next()
			return
		}
		if !m.config.IsLoad {
			ctx.Next()
			return
		}
		if strings.Contains(ctx.Request.RequestURI, "/website/") {
			ctx.Next()
			return
		}
		sessionKey, sessiionValue, err := m.session.GetWithCsrfToken(ctx.Request)
		m.session.SetWithCsrfToken(ctx.Writer, sessionKey, sessiionValue)

		// 更新session过期时间
		if err != nil {
			ctx.JSON(200, errno.ERRNO_AUTH_NOT_LOGIN.ToMapWithError(err))
			ctx.Abort()
			return
		}
		// 无论是否有获取到cookie, 均需要重新设置cookie
		ctx.Next()

	}
}
