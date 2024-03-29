/*
 * @Author: reel
 * @Date: 2023-06-15 07:35:00
 * @LastEditors: reel
 * @LastEditTime: 2024-03-27 04:47:12
 * @Description: 基于gin的上下文进行封装
 */
package core

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/fbs-io/core/pkg/consts"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type context struct {
	ctx  *gin.Context
	core Core
}

const (
	CTX_TX                  = "ctx_tx"                       // 上下文的数据库信息
	CTX_PARAMS              = "ctx_params"                   // 上下文的参数
	CTX_AUTH                = consts.CTX_AUTH                // 上下文的操作用户
	CTX_REFLECT_VALUE       = "reflect_value"                // 上下文中的反射值,用于自动校验并生成参数
	CTX_SHARDING_KEY        = consts.CTX_SHARDING_KEY        // 上下文的数据分区
	CTX_DATA_PERMISSION_KEY = consts.CTX_DATA_PERMISSION_KEY // 上下文的数据权限

	// 通过ctx生成查询tx的方式
	// 适用于表中有id的查询, 通过子查询优化分页性能
	TX_QRY_MODE_SUBID = "subid"
	TX_QRY_DELETE     = true
)

var (
	// 定义请求参数map用于便于中间件拦截
	requestParams = make(map[string]reflect.Type)
)

type Context interface {
	// ShouldBindQuery 反序列化 querystring
	// tag: `form:"xxx"` (注：不要写成 query)
	ShouldBindQuery(obj interface{}) error

	// ShouldBindPostForm 反序列化 postform (querystring会被忽略)
	// tag: `form:"xxx"`
	ShouldBindPostForm(obj interface{}) error

	// ShouldBindForm 同时反序列化 querystring 和 postform;
	// 当 querystring 和 postform 存在相同字段时，postform 优先使用。
	// tag: `form:"xxx"`
	ShouldBindForm(obj interface{}) error

	// ShouldBindJSON 反序列化 postjson
	// tag: `json:"xxx"`
	ShouldBindJSON(obj interface{}) error

	// ShouldBindURI 反序列化 path 参数(如路由路径为 /user/:name)
	// tag: `uri:"xxx"`
	ShouldBindURI(obj interface{}) error

	// Redirect 重定向
	Redirect(code int, location string)

	// HTML 返回界面
	HTML(name string, obj interface{})

	// 返回 Json
	JSON(data interface{})

	// Header 获取 Header 对象
	Header() http.Header

	// GetHeader 获取 Header
	GetHeader(key string) string

	// SetHeader 设置 Header
	SetHeader(key, value string)

	// SetCookie 设置cookie
	SetCookie(key, value string)

	// Cookie 根据cookie的key获取值
	Cookie(key string) (value string, err error)

	// RequestInputParams 获取所有参数
	RequestInputParams() url.Values

	// RequestPostFormParams  获取 PostForm 参数
	RequestPostFormParams() url.Values

	// Request 获取 Request 对象
	Request() *http.Request

	// Method 获取 Request.Method
	Method() string
	// Host 获取 Request.Host
	Host() string
	// Path 获取 请求的路径 Request.URL.Path (不附带 querystring)
	Path() string
	// URI 获取 unescape 后的 Request.URL.RequestURI()
	URI() string
	// resource 获取 请求方式和全路径拼接好的字符串
	// 如GET:/api/v1/userlist
	Resource() string

	// 终止并返回信息
	AbortWithError(interface{})

	// CtxGet 获取上下文自定义的一些参数
	CtxGet(key string) interface{}

	// 设置自定义参数在上下文中
	CtxSet(key string, v interface{})

	// 获取二次封装的参数
	CtxGetParams() any

	Next()

	Ctx() *gin.Context

	// 返回通过参数构建好查询参数参数的gorm.DB
	TX(optFunc ...TxOptsFunc) *gorm.DB

	Core() Core

	// 获取用户
	Auth() string

	// 生成新的db查询
	NewTX(optFunc ...TxOptsFunc) *gorm.DB
}

var _ Context = (*context)(nil)

// 定义上下文池, 减少内存频繁申请开销, 提高性能
var ctxPool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

// 新建一个上下文
func NewCtx(c Core, ctx *gin.Context) Context {
	ct := ctxPool.Get().(*context)
	ct.ctx = ctx
	ct.core = c
	return ct
}

// 回收上下文
func setFreeCtx(ctx Context) {
	ct := ctx.(*context)
	ct.ctx = nil
	ctxPool.Put(ct)
}

// 获取自定义参数
func (ctx *context) CtxGetParams() any {
	p, ok := ctx.ctx.Get(CTX_PARAMS)
	if ok {
		return p
	}
	return nil
}

// ShouldBindQuery 反序列化querystring
// tag: `form:"xxx"` (注：不要写成query)
func (c *context) ShouldBindQuery(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Query)
}

// ShouldBindPostForm 反序列化 postform (querystring 会被忽略)
// tag: `form:"xxx"`
func (c *context) ShouldBindPostForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.FormPost)
}

// ShouldBindForm 同时反序列化querystring和postform;
// 当querystring和postform存在相同字段时，postform优先使用。
// tag: `form:"xxx"`
func (c *context) ShouldBindForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Form)
}

// ShouldBindJSON 反序列化postjson
// tag: `json:"xxx"`
func (c *context) ShouldBindJSON(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.JSON)
}

// ShouldBindURI 反序列化path参数(如路由路径为 /user/:name)
// tag: `uri:"xxx"`
func (c *context) ShouldBindURI(obj interface{}) error {
	return c.ctx.ShouldBindUri(obj)
}

// Redirect 重定向
func (c *context) Redirect(code int, location string) {
	c.ctx.Redirect(code, location)
}

func (c *context) Header() http.Header {
	header := c.ctx.Request.Header

	clone := make(http.Header, len(header))
	for k, v := range header {
		value := make([]string, len(v))
		copy(value, v)

		clone[k] = value
	}
	return clone
}

func (c *context) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}

func (c *context) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

func (c *context) SetCookie(key, value string) {
	c.ctx.SetCookie(key, value, 3600, "/", "", false, true)
}

func (c *context) Cookie(key string) (value string, err error) {
	return c.ctx.Cookie(key)
}

// Method 请求的method
func (c *context) Method() string {
	return c.ctx.Request.Method
}

// Host 请求的host
func (c *context) Host() string {
	return c.ctx.Request.Host
}

// Path 请求的路径(不附带querystring)
func (c *context) Path() string {
	return c.ctx.Request.URL.Path
}

// Path 请求的路径(不附带querystring)
func (c *context) Resource() string {
	return fmt.Sprintf("%s:%s", c.ctx.Request.Method, c.ctx.FullPath())
}

// URI unescape后的uri
func (c *context) URI() string {
	uri, _ := url.QueryUnescape(c.ctx.Request.URL.RequestURI())
	return uri
}

func (c *context) HTML(name string, obj interface{}) {
	c.ctx.HTML(200, name+".html", obj)
}

func (c *context) JSON(data interface{}) {
	en, ok := data.(errno.Errno)
	if ok {
		c.ctx.JSON(en.HTTPCode(), en.ToMap())
		return
	}
	c.ctx.JSON(200, data)
}

// Request 获取 Request
func (c *context) Request() *http.Request {
	return c.ctx.Request
}

// RequestInputParams 获取所有参数
func (c *context) RequestInputParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.Form
}

// RequestPostFormParams 获取 PostForm 参数
func (c *context) RequestPostFormParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.PostForm
}

// 获取gin的上下文
func (c *context) Ctx() *gin.Context {
	return c.ctx
}

// 终止路由并返回错误信息
func (c *context) AbortWithError(err interface{}) {
	c.ctx.JSON(200, err)
	c.ctx.Abort()
}

func (c *context) CtxGet(key string) interface{} {
	v, ok := c.ctx.Get(key)
	if ok {
		return v
	}
	return nil
}

func (c *context) CtxSet(key string, v interface{}) {
	c.ctx.Set(key, v)
}

func (c *context) Next() {
	c.ctx.Next()
}

type txOpts struct {
	mode      string
	qryDelete bool
	tableName string
}

type TxOptsFunc func(*txOpts)

// 设置查询方式
//
// TX_QRY_MODE_SUBID 表示带id的子查询, 注意: 使用TX_QRY_MODE_SUBID, 必须配合使用 SetTxSubTable 设置表名
func SetTxMode(mode string) TxOptsFunc {
	return func(txo *txOpts) {
		txo.mode = mode
	}
}

// 设置子查询表明
func SetTxSubTable(table string) TxOptsFunc {
	return func(txo *txOpts) {
		txo.tableName = table
	}
}

func QryDelete() TxOptsFunc {
	return func(txo *txOpts) {
		txo.qryDelete = true
	}
}

// 通过传入参数设置查询方式
func (ctx *context) TX(optFunc ...TxOptsFunc) (tx *gorm.DB) {
	txopt := &txOpts{
		mode: "",
	}

	for _, optfunc := range optFunc {
		optfunc(txopt)
	}
	sk, _ := ctx.CtxGet(CTX_SHARDING_KEY).(string)
	tx = ctx.Core().RDB().DB().Where("1 = 1")

	// 如果没有参数, 直接返回
	rvi, ok := ctx.ctx.Get(CTX_REFLECT_VALUE)
	if !ok {
		return
	}

	cb := rdb.GenConditionWithParams(rvi.(reflect.Value))
	cb.QryDelete = txopt.qryDelete
	cb.ShardingKey = sk
	if txopt.tableName != "" {
		cb.TableName = txopt.tableName
	}
	switch txopt.mode {
	case TX_QRY_MODE_SUBID:
		tx.Set(rdb.TX_CONDITION_BUILD_KEY, cb)
		tx.Set(rdb.TX_SUB_QUERY_COLUMN_KEY, "id")
	default:
		tx = ctx.Core().RDB().BuildQuery(cb)
	}
	for k, v := range ctx.ctx.Copy().Keys {
		tx.Set(k, v)
	}
	return
}

func (c *context) Core() Core {
	return c.core
}

func (c *context) Auth() (auth string) {
	authI, ok := c.ctx.Get(CTX_AUTH)
	if !ok {
		return
	}
	return authI.(string)
}

func (ctx *context) NewTX(optFunc ...TxOptsFunc) *gorm.DB {
	tx := ctx.Core().RDB().DB().Where("1=1")

	for k, v := range ctx.ctx.Copy().Keys {
		tx.Set(k, v)
	}
	return tx
}
