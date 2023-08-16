/*
 * @Author: reel
 * @Date: 2023-05-28 18:06:12
 * @LastEditors: reel
 * @LastEditTime: 2023-08-01 23:41:27
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
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("Content-type", "application/json")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, gin.H{"code": 0})
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

// var notWriteLogPath = map[string]string{
//     "/static/": "/static/",
// }

// 日志中间件
func (m *handler) log() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		startTime := time.Now()
		ctx.Next()
		if strings.Contains(ctx.Request.RequestURI, "/static/") {
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
	"/ajax/login": true,
	// "/":             true,
	"/ajax/install": true,
}

// 校验签名
// 如果没有登陆, 则会给一个默认的签名
func (m *handler) signature() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if allowPath[ctx.FullPath()] {
			ctx.Next()
			return
		}
		if !m.config.IsLoad {
			ctx.Next()
			return
		}

		sessionKey, sessiionValue, err := m.session.GetWithToken(ctx.Request)
		// 更新session过期时间
		m.session.SetWithToken(sessionKey, sessiionValue)
		if err != nil && sessiionValue == m.session.CookieName() {
			ctx.JSON(200, errno.ERRNO_AUTH_NOT_LOGIN.ToMapWithError(err))
			ctx.Abort()
			return
		}
		// 无论是否有获取到cookie, 均需要重新设置cookie
		ctx.Next()

	}
}
