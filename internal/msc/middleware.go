/*
 * @Author: reel
 * @Date: 2023-05-28 18:06:12
 * @LastEditors: reel
 * @LastEditTime: 2023-06-04 19:26:31
 * @Description: 中间件
 */
package msc

import (
    "core/logx"
    "core/pkg/errno"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"

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
            // c.Header("Access-Control-Allow-Credentials", "true")
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

// cookie 中间件
func (m *handler) validCookie() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if allowPath[ctx.FullPath()] {
            ctx.Next()
            return
        }
        if !m.config.IsLoad {
            ctx.Next()
            return
        }

        sid, err := ctx.Cookie(COOKIE_SID)
        ck := m.cache.Get(sid)
        if err == nil && ck != "" {
            m.cache.Set(COOKIE_SID, sid)
            ctx.SetCookie(COOKIE_SID, sid, 0, "/", "", false, true)
            ctx.Next()
            return
        }
        // ctx.HTML(200, "index.html", "login")
        // ctx.Request.URL.Path = "/"
        ctx.JSON(200, errno.ERRNO_AUTH_NOT_LOGIN.ToMap())
        ctx.Abort()
    }
}

// func (m *handler) validParams() gin.HandlerFunc {
//     return func(ctx *gin.Context) {
//         // ctx.ShouldBindJSON()
//         // key := ctx.Request.Method + ":" + ctx.FullPath()
//         // rt := c.requestParams[key]
//         // if rt != nil {
//         //     params := reflect.New(rt).Interface()
//         //     var err error
//         //     switch ctx.Request.Method {
//         //     case "GET", "DELETE":
//         //         err = ctx.ShouldBindWith(&params, binding.Query)
//         //     case "POST", "PUT":
//         //         err = ctx.ShouldBindJSON(&params)
//         //     }
//         //     if err != nil {
//         //         ctx.JSON(200, errno.ERRNO_PARAMS_BIND.ToMapWithError(errorx.Wrap(err, "校验参数发生错误")))
//         //         ctx.Abort()
//         //         return
//         //     }
//         //     ctx.Set(consts.CTX_PARAMS, params)
//         // }
//         ctx.Next()
//     }
// }
