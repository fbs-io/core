/*
 * @Author: reel
 * @Date: 2023-06-15 07:35:00
 * @LastEditors: reel
 * @LastEditTime: 2024-10-08 00:43:14
 * @Description: 基于gin的上下文进行封装
 */
package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/fbs-io/core/logx"
	"github.com/fbs-io/core/pkg/consts"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/cache"
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
	CTX_AUTH_CACHE          = "ctx_auth_cache"               // 上下文的用户数据缓存
	CTX_LOG_CONTENT         = "ctx_operate_log_content"      // 上下文的操作内容
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
	JSON(data interface{}, funcs ...FuncOperateOpt)

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

	// URI 获取 unescape 后的 Request.URL.RequestURI()
	URI() string
	// Host 获取 Request.Host
	Host() string
	// Path 获取 请求的路径 Request.URL.Path (不附带 querystring)
	Path() string
	// Method 获取 Request.Method
	Method() string

	// resource 获取 请求方式和全路径拼接好的字符串
	//
	// 如GET:/api/v1/userlist
	Resource() string

	// 生成资源code
	//
	// 格式: get:api:user:list
	ResourceCode() string

	// 终止并返回信息
	AbortWithError(interface{})

	// gin的next方法
	Next()

	// gin的abort方法
	Abort()

	// 上下文相关
	Ctx() *gin.Context

	// 获取上下文的core
	Core() Core

	// 获取上下文中的用户
	Auth() string

	// 获取上下文中的分区
	ShardingKey() (sk string)

	// 获取接口的参数
	CtxGetParams() any

	// gin上下文设置
	// CtxGet 获取上下文自定义的一些参数
	CtxGet(key string) interface{}

	// 设置自定义参数在上下文中
	CtxSet(key string, v interface{})

	// 返回通过参数构建好查询参数参数的gorm.DB
	TX(optFunc ...TxOptsFunc) *gorm.DB

	// 生成新的db查询,不含参数及预置的sql条件
	NewTX(optFunc ...TxOptsFunc) *gorm.DB

	// 生成带有子查询的db对象
	SubQueryTX() *gorm.DB

	// 设置缓存
	CacheSet(key, value string, funcs ...cache.OptFunc) error

	// 获取缓存
	CacheGet(key string) string

	// 删除缓存
	CacheDelete(key string) error

	// 设置缓存对象
	CacheSetWithObj(key string, value interface{}, funcs ...cache.OptFunc) error

	//获取缓存对象
	CacheGetWithObj(key string, result interface{}) error

	// 获取分区DB对象, 用于事物处理
	ShardingTx() *gorm.DB

	// 打印日志

	// info日志
	//
	// 默认传递上下文至日志中
	LogInfo(msg string, infoF ...logx.EntityFunc)

	// debug日志
	//
	// 默认传递上下文至日志中
	LogDebug(msg string, infoF ...logx.EntityFunc)

	// warn日志
	//
	// 默认传递上下文至日志中
	LogWarn(msg string, infoF ...logx.EntityFunc)

	// error日志
	//
	// 默认传递上下文至日志中
	LogError(msg string, infoF ...logx.EntityFunc)

	// Fatal日志
	//
	// 默认传递上下文至日志中
	LogFatal(msg string, infoF ...logx.EntityFunc)
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
//
// 格式:GET:/api/xxxxx
func (c *context) Resource() string {
	return requestKey(c.ctx)
}

// 生成资源code
//
// 格式: get:api:user:list
func (c *context) ResourceCode() string {
	return fmt.Sprintf("%s%s", strings.ToLower(c.ctx.Request.Method), strings.ReplaceAll(c.ctx.FullPath(), "/", ":"))
}

// URI unescape后的uri
func (c *context) URI() string {
	uri, _ := url.QueryUnescape(c.ctx.Request.URL.RequestURI())
	return uri
}

func (c *context) HTML(name string, obj interface{}) {
	c.ctx.HTML(200, name+".html", obj)
}

func (c *context) JSON(data interface{}, funcs ...FuncOperateOpt) {

	en, ok := data.(errno.Errno)
	if ok {

		if en.Code() != 0 {
			resource := resourcesMap[c.ResourceCode()]
			if resource != nil {
				en = en.Api(resource.Desc)
			}
		}
		c.ctx.JSON(en.HTTPCode(), en.ToMap())

		// 操作日志
		setOperateLog(c, en, funcs...)
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

func (c *context) Abort() {
	c.ctx.Abort()
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
	tx = ctx.Core().RDB().DB().Where("1 = 1").WithContext(ctx.ctx.Copy())
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
		tx = tx.Set(rdb.TX_CONDITION_BUILD_KEY, cb)
		tx = tx.Set(rdb.TX_SUB_QUERY_COLUMN_KEY, "id")
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

// 只保留上下文参数的db对象
func (ctx *context) NewTX(optFunc ...TxOptsFunc) *gorm.DB {
	tx := ctx.Core().RDB().DB().Where("1=1").WithContext(ctx.Ctx())

	for k, v := range ctx.ctx.Copy().Keys {
		tx = tx.Set(k, v)
	}
	return tx
}

// 返回子查询的db对象
func (ctx *context) SubQueryTX() *gorm.DB {
	return ctx.TX(SetTxMode(TX_QRY_MODE_SUBID))
}

func (ctx *context) ShardingKey() (sk string) {
	skI, ok := ctx.ctx.Get(CTX_SHARDING_KEY)
	if !ok {
		return
	}
	return skI.(string)
}

// 设置缓存
func (ctx *context) CacheSet(key, value string, funcs ...cache.OptFunc) error {
	return ctx.core.Cache().Set(key, value, funcs...)
}

// 获取缓存
func (ctx *context) CacheGet(key string) string {
	return ctx.core.Cache().Get(key)
}

// 设置缓存对象
func (ctx *context) CacheSetWithObj(key string, value interface{}, funcs ...cache.OptFunc) error {
	vb, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ctx.core.Cache().Set(key, string(vb), funcs...)
}

// 根据key获取单个缓存对象
func (ctx *context) CacheGetWithObj(key string, res interface{}) error {
	vbs := ctx.core.Cache().Get(key)
	if vbs == "" {
		return errors.New("没有缓存数据")
	}
	return json.Unmarshal([]byte(vbs), res)
}

// 删除缓存
func (ctx *context) CacheDelete(key string) error {
	return ctx.core.Cache().Del(key)
}

// 获取分区gorm对象
//
// 主要用于事物处理, 正常增删改查可不用该方法
//
// 无上下文的自动处理的查询参数
func (ctx *context) ShardingTx() *gorm.DB {
	tx := ctx.core.RDB().GetShardingDB(ctx.ShardingKey())
	tx = tx.Where("1=1")
	for k, v := range ctx.ctx.Copy().Keys {
		tx = tx.Set(k, v)
	}
	tx = tx.Set(consts.CTX_SHARDING_DB, ctx.ShardingKey()).WithContext(ctx.Ctx())
	return tx
}

// info日志
//
// 默认传递上下文至日志中
func (ctx *context) LogInfo(msg string, infoF ...logx.EntityFunc) {
	logx.APP.Info(msg, append(infoF, logx.Context(ctx.ctx))...)
}

// debug日志
//
// 默认传递上下文至日志中
func (ctx *context) LogDebug(msg string, infoF ...logx.EntityFunc) {
	logx.APP.Debug(msg, append(infoF, logx.Context(ctx.ctx))...)
}

// warn日志
//
// 默认传递上下文至日志中
func (ctx *context) LogWarn(msg string, infoF ...logx.EntityFunc) {
	logx.APP.Warn(msg, append(infoF, logx.Context(ctx.ctx))...)
}

// error日志
//
// 默认传递上下文至日志中
func (ctx *context) LogError(msg string, infoF ...logx.EntityFunc) {
	logx.APP.Error(msg, append(infoF, logx.Context(ctx.ctx))...)
}

// Fatal日志
//
// 默认传递上下文至日志中
func (ctx *context) LogFatal(msg string, infoF ...logx.EntityFunc) {
	logx.APP.Fatal(msg, append(infoF, logx.Context(ctx.ctx))...)
}

func setOperateLog(ctx Context, en errno.Errno, funcs ...FuncOperateOpt) {
	defer func() {
		setFreeCtx(ctx)
		if err := recover(); err != nil {
			logx.Sys.Error("写入操作日志发生错误", logx.F("status", ctx.Ctx().Writer.Status()),
				logx.Context(ctx.Ctx()),
			)
			return
		}
	}()

	// 设置配置项
	opt := &operateOpt{}
	for _, fs := range funcs {
		fs(opt)
	}
	// 如果配置项都为空, 则不需要写入操作日志
	if !opt.isSet && opt.content == "" && opt.result == nil {
		return
	}
	resource := resourcesMap[fmt.Sprintf("%s%s", strings.ToLower(ctx.Ctx().Request.Method), strings.Replace(ctx.Ctx().FullPath(), "/", ":", -1))]

	auth := ""
	authI, ok := ctx.Ctx().Get(CTX_AUTH)
	if ok {
		auth = authI.(string)
	}
	res := "成功"

	if en.Code() != errno.ERRNO_OK.Code() {
		res = "失败"
	}

	content := fmt.Sprintf("%s%s%v%s", auth, resource.Desc, opt.result, res)
	if opt.result == nil {
		content = fmt.Sprintf("%s%s%s", auth, resource.Desc, res)
		if res == "失败" {
			if en.Details() == nil {
				content = fmt.Sprintf("%s%s%s, 错误:%s", auth, resource.Desc, res, en.Message())
			} else {
				content = fmt.Sprintf("%s%s%s, 错误:%s details:%v", auth, resource.Desc, res, en.Message(), en.Details())
			}
		}
	}

	if opt.content != "" {
		content = opt.content
	}
	operateLog := &OperateLog{
		IP:        ctx.Ctx().ClientIP(),
		User:      auth,
		Content:   content,
		Result:    res,
		Method:    ctx.Ctx().Request.Method,
		Api:       ctx.Ctx().Request.RequestURI,
		ApiName:   resource.Desc,
		TraceID:   ctx.Ctx().Request.Header.Get(consts.REQUEST_HEADER_TRACE_ID),
		OperateID: ctx.Ctx().Request.Header.Get(consts.REQUEST_HEADER_OPERATE_ID),
	}
	operateLog.ShadingKey = ctx.ShardingKey()
	err := ctx.Core().RDB().DB().Create(operateLog).Error
	if err != nil {
		logx.Sys.Error("写入操作日志失败", logx.EV(err))
	}

}
