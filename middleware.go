/*
 * @Author: reel
 * @Date: 2023-07-19 00:08:08
 * @LastEditors: reel
 * @LastEditTime: 2024-10-05 00:04:08
 * @Description: 常用的中间件
 */
package core

import (
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/fbs-io/core/logx"
	"github.com/fbs-io/core/pkg/consts"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/pkg/trace"
	"github.com/fbs-io/core/session"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	COOKIE_SID               = "ASID"
	SINGULAR_TYPE_COOKIE     = "cookie"
	SINGULAR_TYPE_SID        = "sid"
	SINGULAR_TYPE_TOKEN      = "token"
	SINGULAR_TYPE_CSRF_TOKEN = "CSRF-TOKEN"
)

var (
	STATIC_PATH_PREFIX = "/static/"
)

func SetStaticPathPrefix(prefix string) {
	STATIC_PATH_PREFIX = prefix
}

// 跨域处理中间件
func CorsMiddleware(c Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin") //请求头部
		if origin != "" && !strings.Contains(ctx.FullPath(), STATIC_PATH_PREFIX) {
			//接收客户端发送的origin （重要！）
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			//服务器支持的所有跨域请求的方法
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Content-Length, X-CSRF-TOKEN, Token, session, TRACE-ID")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, X-CSRF-TOKEN, SID, TRACE-ID")
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
		if strings.Contains(ctx.Request.RequestURI, STATIC_PATH_PREFIX) {
			ctx.Next()
			return
		}
		startTime := time.Now()
		logx.APP.Info("http请求开始", logx.F("status", ctx.Writer.Status()),
			logx.Context(ctx),
			logx.F("client_ip", ctx.ClientIP()),
			logx.F("req_method", ctx.Request.Method),
			logx.F("req_url", ctx.Request.RequestURI),
		)
		defer func() {
			if err := recover(); err != nil {
				logx.APP.Warn("http请求发生错误", logx.F("status", ctx.Writer.Status()),
					logx.Context(ctx),
					logx.F("client_ip", ctx.ClientIP()),
					logx.F("req_method", ctx.Request.Method),
					logx.F("req_url", ctx.Request.RequestURI),
					logx.F("error", fmt.Sprintf("msg: %v, stack: %s", err, string(debug.Stack()))),
				)
				ctx.JSON(200, errno.ERRNO_SYSTEM.ToMapWithError(errorx.Errorf("%v", err)))
				ctx.Abort()
				return
			}
			logx.APP.Info("http请求结束", logx.F("status", ctx.Writer.Status()),
				logx.Context(ctx),
				logx.DiffTime(startTime),
				logx.F("client_ip", ctx.ClientIP()),
				logx.F("req_method", ctx.Request.Method),
				logx.F("req_url", ctx.Request.RequestURI),
			)
		}()

		ctx.Next()

	}
}

var (
	allowSource = make(map[string]bool, 1000) // 资源访问例外
)

// 请使用 method:path 的方式定义资源
// 比如 POST:/ajax/login
func AddAllowSource(resoures ...string) {
	for _, resoure := range resoures {
		allowSource[resoure] = true
	}
}

func GetAllowSource(ctx *gin.Context) bool {
	return allowSource[requestKey(ctx)]
}

func requestKey(ctx *gin.Context) string {
	return fmt.Sprintf("%s:%s", ctx.Request.Method, ctx.FullPath())
}

type signatureOption struct {
	sessionType string // session类型, 以那种方式存储session
}

type SignatureOptFunc func(*signatureOption)

func SetSeesionType(sessionType string) SignatureOptFunc {
	return func(o *signatureOption) {
		o.sessionType = sessionType
	}
}

// 校验签名中间件
// 如果没有登陆, 则会给一个默认的签名
// Singular: 默认 token 模式， 同时可以选择cookie，sid, csrftoken方式
func SignatureMiddleware(c Core, sop ...SignatureOptFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.Request.RequestURI, STATIC_PATH_PREFIX) {
			ctx.Next()
			return
		}

		startTime := time.Now()
		logx.APP.Debug("SignatureMiddleware 签名校验开始",
			logx.Context(ctx),
		)
		defer logx.APP.Debug("SignatureMiddleware 签名校验结束",
			logx.Context(ctx),
			logx.DiffTime(startTime),
		)
		if allowSource[requestKey(ctx)] {
			logx.APP.Debug("SignatureMiddleware 例外的请求接口无需校验",
				logx.Context(ctx),
			)
			return
		}

		var (
			sessionKey    string
			sessiionValue string
			err           error
		)
		var opt = &signatureOption{
			sessionType: SINGULAR_TYPE_CSRF_TOKEN,
		}
		for _, fc := range sop {
			fc(opt)
		}
		switch opt.sessionType {
		case SINGULAR_TYPE_COOKIE:
			sessionKey, sessiionValue, err = c.Session().GetWithCookie(ctx.Request)
			c.Session().SetWithCookie(ctx.Writer, sessionKey, sessiionValue)
		case SINGULAR_TYPE_SID:
			sessionKey, sessiionValue, err = c.Session().GetWithSid(ctx.Request)
			c.Session().SetWithSid(ctx.Writer, sessionKey, sessiionValue)

		case SINGULAR_TYPE_CSRF_TOKEN:
			sessionKey, sessiionValue, err = c.Session().GetWithCsrfToken(ctx.Request)
			c.Session().SetWithCsrfToken(ctx.Writer, sessionKey, sessiionValue)

		default:
			sessionKey, sessiionValue, err = c.Session().GetWithToken(ctx.Request)
			// 更新session过期时间
			c.Session().SetWithToken(sessionKey, sessiionValue)
		}
		if err != nil {
			if err == session.ERROR_SESSION_ELSE_LOGIN {
				ctx.JSON(200, errno.ERRNO_AUTH_ELSE_LOGIN.ToMapWithError(err))
			} else {
				ctx.JSON(200, errno.ERRNO_AUTH_NOT_LOGIN.ToMapWithError(err))
			}
			logx.APP.Debug("SignatureMiddleware 签名校验不通过",
				logx.E(err),
				logx.Context(ctx),
			)
			ctx.Abort()
			return
		}
		// 用户鉴权成功后, 把用户信息写入上下文用于数据的查询,记录等
		ctx.Set(CTX_AUTH, sessiionValue)
		logx.APP.Debug("SignatureMiddleware 签名校验通过, 用户信息写入上下文传递",
			logx.Context(ctx),
		)

	}
}

// 参数自动生成中间件
//
// 会生成参数结构体以及gorm.DB
//
// 同时根据约束, 自动完成参数校验
func ParamsMiddleware(c Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		logx.APP.Debug("ParamsMiddleware 参数生成开始",
			logx.Context(ctx),
		)
		var params interface{}

		rt := requestParams[requestKey(ctx)]
		if rt == nil {
			logx.APP.Debug("ParamsMiddleware 无法获取参数配置, 无法生成参数",
				logx.Context(ctx),
			)
			return
		}
		rv := reflect.New(rt)
		params = rv.Interface()
		var err error
		switch ctx.ContentType() {
		case formContent:
			err = ctx.ShouldBindQuery(params)
		case jsonContent:
			err = ctx.ShouldBindJSON(&params)
		default:
			err = ctx.ShouldBindQuery(params)
		}
		if err != nil {
			logx.APP.Debug("ParamsMiddleware 参数生成发生错误",
				logx.Context(ctx),
				logx.EV(err),
			)
			ctx.JSON(200, errno.ERRNO_PARAMS_BIND.ToMapWithError(errorx.Wrap(err, "校验参数发生错误")))
			ctx.Abort()
			return
		}
		ctx.Set(CTX_PARAMS, params)
		ctx.Set(CTX_REFLECT_VALUE, rv)
		logx.APP.Debug("ParamsMiddleware 参数生成结束, 参数写入上下文",
			logx.Context(ctx),
			logx.DiffTime(startTime),
			logx.F("params", params),
		)
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

// 链路ID校验
//
// 如果没有trace_id, 不允许访问
func TraceMiddleware(c Core, hasUiTrace bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.Request.Header.Get(consts.REQUEST_HEADER_TRACE_ID)
		if traceID == "" && !allowSource[requestKey(ctx)] && !strings.Contains(ctx.Request.RequestURI, STATIC_PATH_PREFIX) {
			if hasUiTrace {
				ctx.JSON(200, errno.ERRNO_NO_TRACE_REQUESTS.ToMap())
				ctx.Abort()
				return
			}
			traceID = uuid.New().String()
		}
		trace := &trace.Trace{TraceID: traceID}
		ctx.Set(consts.CTX_TRACE_ID, trace)
		ctx.Next()
	}
}
