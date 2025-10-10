package core

import (
	// "fbs/internal/core/apidoc"
	// "fbs/internal/core/means"
	// "fbs/pkg/convx"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(Context)

func wrapHandlers(c Core, handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		funcs[i] = func(ginCtx *gin.Context) {
			ctx := NewCtx(c, ginCtx)
			defer setFreeCtx(ctx)
			h(ctx)
		}
	}
	return funcs
}

// RouterGroup 包装gin的RouterGroup
//
// 传入相对路由路径和名称, 用于在菜单中进行注册
type RouterGroup interface {
	// 创建分组
	Group(api, apiName string, handlers ...HandlerFunc) RouterGroup

	// 设置中间件
	Use(middleware ...gin.HandlerFunc) RouterGroup

	// 获取core
	Core() Core

	IRoutes

	RouterResource
}

var _ IRoutes = (*router)(nil)

// IRoutes 包装gin的IRoutes
type IRoutes interface {
	// Any(string, ...HandlerFunc)
	//需要填写相对路由路径, 名称, 参数, 及中间件, 用于在 api 文档和菜单中注册
	//参数为如果为空, 该方法不会在 api 文档中进行注册
	GET(relativePath, pathName string, params any, handlers ...HandlerFunc) (resource *Resources)
	//需要填写相对路由路径, 名称, 参数, 及中间件, 用于在 api 文档和菜单中注册
	//参数为如果为空, 该方法不会在 api 文档中进行注册
	PUT(relativePath, pathName string, params any, handlers ...HandlerFunc) (resource *Resources)
	//需要填写相对路由路径, 名称, 参数, 及中间件, 用于在 api 文档和菜单中注册
	//参数为如果为空, 该方法不会在 api 文档中进行注册
	POST(relativePath, pathName string, params any, handlers ...HandlerFunc) (resource *Resources)
	//需要填写相对路由路径, 名称, 参数, 及中间件, 用于在 api 文档和菜单中注册
	//参数为如果为空, 该方法不会在 api 文档中进行注册
	DELETE(relativePath, pathName string, params any, handlers ...HandlerFunc) (resource *Resources)
	// TODO: 以后根据业务进行扩展
	// PATCH(string, ...HandlerFunc)
	// OPTIONS(string, ...HandlerFunc)
	// HEAD(string, ...HandlerFunc)
}

type router struct {
	group    *gin.RouterGroup
	resource *Resources
	core     *core
}

var (
	routers = make(map[string]*router, 100)
	lock    = &sync.Mutex{}
)

// 获取路由
// 如果路由已存在, 直接返回, 防止重复生成
func getRouter(relativePath string) (rout *router) {
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
	return routers[relativePath]
}

// 设置路由
// 同时对路由设置进行加锁
// 设置路由时, 同时完成对资源表的写入
func setRouter(relativePath, pathName string, r *router) {
	lock.Lock()
	defer lock.Unlock()

	if routers[relativePath] != nil {
		return
	}
	if pathName == "" {
		pathName = relativePath
	}
	resource := r.genResources(relativePath, pathName, "")

	// 默认路由组为菜单, 均需要授权才能访问
	resource.Type = SOURCE_TYPE_MENU
	resource.IsRouter = SOURCE_ROUTER_IS
	resourcesMap[resource.Code] = resource
	resources = append(resources, resource)
	r.resource = resource
	routers[relativePath] = r

}

// 顶层路由分组, 使用路由的入口
func (c *core) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	r := getRouter(relativePath)
	if r == nil {
		r = &router{
			group: c.Engine().Group(relativePath, wrapHandlers(c, handlers...)...),
			core:  c,
		}
	}
	setRouter(relativePath, relativePath, r)
	return r
}

// 对gin路由分组的封装
func (r *router) Group(relativePath, pathName string, handlers ...HandlerFunc) RouterGroup {
	group := r.group.Group(relativePath, wrapHandlers(r.core, handlers...)...)

	rout := getRouter(relativePath)
	if rout == nil {
		rout = &router{
			group: group,
			core:  r.core,
		}
	}
	setRouter(relativePath, pathName, rout)

	return rout
}

// Get请求方式封装
//
// 参数如果为空, 该方法不会被记录在资源表中
func (r *router) GET(relativePath, pathName string, params any, handlers ...HandlerFunc) (source *Resources) {
	// handlers = append([]HandlerFunc{r.validParams()}, handlers...)
	r.group.GET(relativePath, wrapHandlers(r.core, handlers...)...)
	return r.operation("GET", relativePath, pathName, params)
}

// Post请求方式封装
//
// 参数如果为空, 该方法不会被记录在资源表中
func (r *router) POST(relativePath, pathName string, params any, handlers ...HandlerFunc) (source *Resources) {
	r.group.POST(relativePath, wrapHandlers(r.core, handlers...)...)
	return r.operation("POST", relativePath, pathName, params)
}

// Delete请求方式封装
//
// 参数如果为空, 该方法不会被记录在资源表中
func (r *router) DELETE(relativePath, pathName string, params any, handlers ...HandlerFunc) (source *Resources) {
	r.group.DELETE(relativePath, wrapHandlers(r.core, handlers...)...)
	return r.operation("DELETE", relativePath, pathName, params)
}

// Put请求方式封装
//
// 参数如果为空, 该方法不会被记录在资源表中
func (r *router) PUT(relativePath, pathName string, params any, handlers ...HandlerFunc) (source *Resources) {
	r.group.PUT(relativePath, wrapHandlers(r.core, handlers...)...)
	return r.operation("PUT", relativePath, pathName, params)
}

func (r *router) PATCH(relativePath string, handlers ...HandlerFunc) {
	r.group.PATCH(relativePath, wrapHandlers(r.core, handlers...)...)
}

func (r *router) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	r.group.OPTIONS(relativePath, wrapHandlers(r.core, handlers...)...)
}

func (r *router) HEAD(relativePath string, handlers ...HandlerFunc) {
	r.group.HEAD(relativePath, wrapHandlers(r.core, handlers...)...)
}

// 处理参数生成逻辑
func (r *router) operation(method, relativePath, pathName string, params any) (source *Resources) {
	if relativePath == "" {
		relativePath = "/"
	}
	var (
		paramStr, acceptType string
	)
	source = r.genResources(relativePath, pathName, method)
	// 默认资源均需要授权才能访问
	source.Type = SOURCE_TYPE_PERMISSION
	resources = append(resources, source)
	resourcesMap[source.Code] = source

	if params != nil {
		rt := reflect.TypeOf(params)
		requestParams[fmt.Sprintf("%s:%s/%s", method, r.group.BasePath(), relativePath)] = rt
		paramStr, acceptType = r.genResourcesParams(method, relativePath, rt)

	}
	// 每个接口的参数存放在变量中便于后面查询使用
	source.Params, source.AcceptType = paramStr, acceptType

	return
}

// 用于生成系统资源结构体
//
// 用于API文档和前端菜单
//
// 同时可用于权限及自动生成gorm查询参数
func (r *router) genResources(relativePath, name, method string) *Resources {
	basePaths := strings.Split(r.group.BasePath(), "/")[1:]
	method = strings.ToLower(method)
	s := &Resources{}
	// s.Meta = make(map[string]any, 10)
	fullpath := basePaths
	var metaType = "menu"
	if method != "" {
		basePaths = append(basePaths, relativePath)
		fullpath = append([]string{method}, basePaths...)
		s.Method = method
		metaType = "button"
		s.IsRouter = SOURCE_ROUTER_NAN
	}
	s.Name = relativePath
	s.Desc = name
	s.Code = strings.Join(fullpath, ":")
	s.PCode = strings.Join(basePaths[:len(basePaths)-1], ":")
	s.Level = int8(len(basePaths))
	s.Api = fmt.Sprintf("/%s", path.Join(basePaths...))
	ps := resourcesMap[s.PCode]

	s.Path = fmt.Sprintf("/%s", relativePath)
	s.Component = relativePath
	if ps != nil {
		if ps.Path != "" {
			s.Path = fmt.Sprintf("%s/%s", ps.Path, relativePath)
		}
		if ps.Component != "" {
			s.Component = fmt.Sprintf("%s/%s", ps.Component, relativePath)

		}
	}
	s.Meta = map[string]any{
		"title": name,
		"icon":  "menu",
		"type":  metaType,
	}

	return s
}

const (
	jsonContent = "application/json"
	formContent = "application/x-www-form-urlencoded"

	// 标签相关
	tagJson    = "json"
	tagForm    = "form"
	tagDesc    = "desc"
	tagBinding = "binding"
	tagDefault = "default"
	tagGorm    = "gorm"
	tagView    = "views"

	// 参数相关
	paramsKey           = "key"
	paramsValue         = "value"
	paramsValueType     = "value_type"
	paramsValueInt      = "int"
	paramsValueNum      = "number"
	paramsValueBool     = "bool"
	paramsValueFloat    = "float"
	paramsValueDate     = "date"
	paramsValueDatetiem = "datetime"
)

// 根据参数结构体生成API参数,
//
// 当前仅支持 form 和 json 两种格式
//
// 根据参数第一个字段的标签判断content-type类型
//
// 如果参数结构体定义多个参数格式, 将其他参数无法正确使用
//
// TODO: 支持文件/多文件参数定义
func (r *router) genResourcesParams(method, pathName string, rt reflect.Type) (params string, contentType string) {
	if rt == nil {
		return
	}
	data := make([]any, 0)
	viewsList := make([]*Views, 0, 100)
	// viewsMap := make(map[string]*Views, 100)
	for i := 0; i < rt.NumField(); i++ {
		item := make(map[string]any, 4)
		field := rt.Field(i)

		// 获取前端参数名称
		key := field.Tag.Get(tagForm)
		contentTypeCustom := formContent
		if key == "" {
			contentTypeCustom = jsonContent
			key = field.Tag.Get(tagJson)
		}
		// TODO: 增加其他类型检查

		// 如果没有获取到key, 说明该参数无效, 跳过不在录入
		if key == "" {
			continue
		}
		// 通过第一个获取到参数的结构体的类型作为整个请求的content_type
		if contentType == "" {
			contentType = contentTypeCustom
		}
		view := &Views{
			ResourceCode:  r.resource.Code,
			ViewCode:      pathName,
			ViewType:      "form",
			ViewRole:      method,
			ValueType:     "string",
			Code:          key,
			FormatterType: "input",
		}

		item[paramsKey] = key

		// 前端参数的数据类型
		typeStr := field.Type.String()
		item[paramsValueType] = typeStr
		// 后端参数类型转换为前端的参数类型
		if strings.Contains(typeStr, paramsValueInt) {
			item[paramsValueType] = paramsValueNum
		} else if strings.Contains(typeStr, paramsValueFloat) {
			item[paramsValueType] = paramsValueNum
		}
		if strings.Contains(typeStr, "[]") {
			item[paramsValueType] = "[]:" + paramsValueNum
			view.DefaultValue = "[]"
		}

		// 用于前端API文档中的默认值
		item[paramsValue] = field.Tag.Get(tagDefault)
		// 用于字段描述
		item[tagDesc] = field.Tag.Get(tagDesc)
		// 用于校验参数信息
		item[tagBinding] = field.Tag.Get(tagBinding)
		views := map[string]any{}

		view.Name = field.Tag.Get(tagDesc)
		view.ValueType = item[paramsValueType].(string)
		view.Rules = field.Tag.Get(tagBinding)

		for _, val := range strings.Split(field.Tag.Get(tagView), ";") {
			kvs := strings.Split(val, ":")
			k := kvs[0]
			if len(k) == 0 {
				continue
			}
			var v any
			switch kvs[0] {
			case "select":
				v = key
				view.FormatterType = "select"
				view.Formatter = view.Code
				if len(kvs) > 1 {
					view.Formatter = kvs[1]
				}
			case "multiple":
				view.Multiple = "Y"
			case "disabled", "key":
				view.Disabled = "Y"
			case "hidden":
				view.Hidden = 1
			case "switch":
				view.FormatterType = "switch"
				if view.ValueType == "string" {
					view.CustomValue = []any{"Y", "N"}
					view.DefaultValue = "N"
				} else if view.ValueType == paramsValueNum {
					view.CustomValue = []any{1, -1}
					view.DefaultValue = "-1"
				}
			case "date":
				view.FormatterType = "date"
				view.Formatter = "YYYY-MM-DD"
				if len(kvs) > 1 {
					view.Formatter = kvs[1]
				}
			case "filter":
				view.Filter = kvs[1]
			case "span":
				span, _ := strconv.Atoi(kvs[1])
				view.Width = int16(span)
			case "depend":
				view.Depend = kvs[1]
			case "default":
				view.DefaultValue = kvs[1]
			}

			if len(kvs) > 1 {
				v = kvs[1]
			}
			views[k] = v
		}
		item[tagView] = views
		data = append(data, item)
		paramsB, _ := json.Marshal(data)
		params = string(paramsB)
		if view.Code == "page_num" || view.Code == "page_size" || view.Code == "orders" {
			continue
		}
		if view.Name == "" {
			sub := r.core.ViewsMap[fmt.Sprintf("%s:%s", view.ResourceCode, view.Code)]
			if sub != nil {
				view.Name = sub.Name
			}
		}
		if view.Name == "" {
			view.Name = strings.ToUpper(view.Code)
		}
		if view.FormatterType == "input" && view.ValueType == "number" {
			view.FormatterType = "number"
		}
		viewsList = append(viewsList, view)
	}

	r.core.Views = append(r.core.Views, viewsList...)
	return
}

// 用于设置某些路由不必写入资源库
func (r *router) NotWithSource() RouterGroup {
	delete(resourcesMap, r.resource.Code)
	r.resource = nil
	return r
}

// 用于设置设置模块/api权限
//
// 默认所有菜单,api都需要权限设置
//
// 对于通用模块如用户个人设置等信息, 可以设置为权限例外, 可以灵活的在初始化阶段完成权限配置
func (r *router) WithPermission(t int8) RouterGroup {
	r.resource.WithPermission(t)
	return r
}

// 去除前端菜单路由前缀,
//
// path: /ajax/user/list => /user/list
// Component : ajax/user/list => user/list
func (r *router) WithMenuNotPrefix(prefix string) RouterGroup {
	r.resource.WithMenuNotPrefix(prefix)
	return r
}

// 设置为前端路由或非路由
//
// 模块部分默认为前端路由(菜单), 默认值 SOURCE_ROUTER_IS = 1, api默认非路由, 为按钮权限, 默认值 SOURCE_ROUTER_NAN=0
//
// 通过该方法可以灵活的配置接口/模块的显示规则
func (r *router) WithRouter(t int8) RouterGroup {
	r.resource.WithRouter(t)
	return r
}

// 设置为前端组件路径
//
// 组件默认路径和后端路径一样, 通过该方法可以设置自定义前端组件路径
//
// 例如: /ajax/user/list => /user/list
func (r *router) WithComponent(component string) RouterGroup {
	r.resource.Component = component
	return r
}

// 获取前端组件路径
func (r *router) Component() string {
	return r.resource.Component

}

// 用于设置路由和资源的关系
type RouterResource interface {

	// 用于设置某些路由不必写入资源库
	NotWithSource() RouterGroup

	// 用于设置设置模块/api权限
	//
	// 默认所有菜单,api都需要权限设置
	//
	// 对于通用模块如用户个人设置等信息, 可以设置为权限例外, 可以灵活的在初始化阶段完成权限配置
	WithPermission(t int8) RouterGroup

	// 设置为前端路由或非路由
	//
	// 模块部分默认为前端路由(菜单), 默认值 SOURCE_ROUTER_IS = 1, api默认非路由, 为按钮权限, 默认值 SOURCE_ROUTER_NAN=0
	//
	// 通过该方法可以灵活的配置接口/模块的显示规则
	WithRouter(t int8) RouterGroup

	// 去除前端菜单路由前缀,
	//
	// 例如api接口转为前端路由: /ajax/user/list => /user/list
	WithMenuNotPrefix(prefix string) RouterGroup

	// 设置路由隐藏
	WithHidden() RouterGroup

	// 设置前端Meta信息
	WithMeta(key string, value any) RouterGroup

	// 设置为前端组件路径
	//
	// 组件默认路径和后端路径一样, 通过该方法可以设置自定义前端组件路径, 可以充分复用前端组件
	//
	// 例如: /ajax/user/list => /user/list
	WithComponent(component string) RouterGroup

	// 获取前端组件路径
	Component() string

	// 自定义前端组件字段展示
	WithViews(item any, fs ...FuncSetViews) RouterGroup
}

// 设置路由隐藏
func (r *router) WithHidden() RouterGroup {
	r.resource.WithHidden()
	return r
}

// 返回Core, 用于在应用模块快速使用Core资源
func (r *router) Core() Core {
	return r.core
}

// 设置前端Meta信息
func (r *router) WithMeta(key string, value any) RouterGroup {
	r.resource.Meta[key] = value
	return r
}

func (r *router) Use(middleware ...gin.HandlerFunc) RouterGroup {
	r.group.Use(middleware...)
	return r
}

// 设置接口的返回字段
//
// 仅支持struct	WithViews( struct ), 字段设置参数仅用于批量设置
//
// 标签支持 字段: json, key, 描述: gorm, desc,
func (r *router) WithViews(item any, fs ...FuncSetViews) RouterGroup {
	rt := reflect.TypeOf(item)
	options := &SetViewOptions{}
	for _, f := range fs {
		f(options)
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		view := r.genViewColumns(field, options)
		if view == nil {
			continue
		}

		key := view.ResourceCode + ":" + view.Code

		if r.core.ViewsMap[key] == nil {
			r.core.Views = append(r.core.Views, view)
			r.core.ViewsMap[key] = view
		}

	}

	return r
}

func (r *router) genViewColumns(field reflect.StructField, opt *SetViewOptions) (view *Views) {
	view = &Views{
		ResourceCode: r.resource.Code,
		ViewCode:     r.resource.Name,
		ViewType:     "table",
		ValueType:    "string",
		Hidden:       -1,
		Width:        120,
		Height:       0,
		IsOrder:      1,
		Filter:       "",
		Fixed:        "N",
		Account:      "system",
	}

	// 获取前端参数名称
	view.Code = field.Tag.Get(tagJson)
	// TODO: 增加其他类型检查

	// 处理前端参数名称, 从gorm标签或者自定义的desc标签获取
	view.Name = field.Tag.Get(tagDesc)
	if view.Name == "" {
		for _, tag := range strings.Split(field.Tag.Get(tagGorm), ";") {
			if strings.Contains(tag, "comment:") {
				view.Name = strings.Split(tag, ":")[1]
				break
			}
		}
	}
	typeStr := field.Type.String()

	// 后端参数类型转换为前端的参数类型
	if strings.Contains(typeStr, paramsValueInt) {
		view.ValueType = paramsValueNum
	} else if strings.Contains(typeStr, paramsValueFloat) {
		view.ValueType = paramsValueNum
	} else if strings.Contains(typeStr, paramsValueBool) {
		view.ValueType = paramsValueBool
	} else if strings.Contains(typeStr, paramsValueBool) {
		view.ValueType = paramsValueBool
	}
	if view.ValueType == paramsValueNum {
		view.FormatterType = "number"
	}
	// 对view标签进行处理
	genViewTag(view, field.Tag.Get(tagView))

	if view.Code == "" {
		return nil
	}
	if opt.ViewCode != "" {
		view.ViewCode = opt.ViewCode
	}

	if opt.ColumnWidth > 0 {
		view.Width = opt.ColumnWidth
	}

	if opt.ColumnHeight > 0 {
		view.Height = opt.ColumnHeight
	}

	if opt.ColumnIsHidden > 0 {
		view.Hidden = opt.ColumnIsHidden
	}
	if opt.ColumnIsOrder > 0 {
		view.IsOrder = opt.ColumnIsOrder
	}
	if opt.ColumnFilter != "" {
		view.Filter = opt.ColumnFilter
	}
	if opt.ColumnFixed != "" {
		view.Fixed = opt.ColumnFixed
	}
	if opt.Account != "" {
		view.Account = opt.Account
	}
	return
}

// view标签处理
func genViewTag(view *Views, viewTag string) {
	for _, tag := range strings.Split(viewTag, ";") {
		if strings.Contains(tag, "code:") {
			view.Code = strings.Split(tag, ":")[1]
		}
		if strings.Contains(tag, "name:") {
			view.Name = strings.Split(tag, ":")[1]
		}
		if strings.Contains(tag, "width:") {
			i, _ := strconv.Atoi(strings.Split(tag, ":")[1])
			if i > 0 {
				view.Width = int16(i)
			}
		}
		if strings.Contains(tag, "height:") {
			i, _ := strconv.Atoi(strings.Split(tag, ":")[1])
			if i > 0 {
				view.Height = int16(i)
			}
		}
		if strings.Contains(tag, "hidden") {
			view.Hidden = 1
		}
		if strings.Contains(tag, "isorder:") {
			switch strings.Split(tag, ":")[1] {
			case "true":
				view.IsOrder = 1
			case "false":
				view.IsOrder = -1
			default:
				view.IsOrder = 1
			}
		}
		if strings.Contains(tag, "filter:") {
			view.Filter = strings.Split(tag, ":")[1]
		}
		if strings.Contains(tag, "fixed:") {
			view.Fixed = strings.Split(tag, ":")[1]
		}
		if strings.Contains(tag, "formatter_type:") {
			view.FormatterType = strings.Split(tag, ":")[1]

		}
		if strings.Contains(tag, "formatter") {
			view.FormatterType = "text"
			kvs := strings.Split(tag, ":")
			k := kvs[0]
			if len(k) == 0 {
				continue
			}
			var v = view.Code

			if len(kvs) > 1 {
				v = kvs[1]
			}
			view.Formatter = v
		}

	}
}

func (r *router) SetViews(view *Views) RouterGroup {
	key := view.ResourceCode + ":" + view.Code

	r.core.ViewsMap[key].Code = view.Code
	r.core.ViewsMap[key].Name = view.Name
	r.core.ViewsMap[key].Width = view.Width
	r.core.ViewsMap[key].Height = view.Height
	r.core.ViewsMap[key].Hidden = view.Hidden
	r.core.ViewsMap[key].IsOrder = view.IsOrder
	r.core.ViewsMap[key].Filter = view.Filter
	r.core.ViewsMap[key].Fixed = view.Fixed
	r.core.ViewsMap[key].Account = view.Account

	return r
}

func (r *router) GetViews(ColumnCode string) *Views {
	key := r.resource.Code + ":" + ColumnCode
	return r.core.ViewsMap[key]
}

type SetViewOptions struct {
	ViewCode       string
	ColumnWidth    int16
	ColumnHeight   int16
	ColumnIsHidden int8
	ColumnIsOrder  int8
	ColumnFilter   string
	ColumnFixed    string
	Account        string
}

type FuncSetViews func(*SetViewOptions)

// 设置视图名称
func SetViewCode(viewCode string) FuncSetViews {
	return func(svo *SetViewOptions) {
		svo.ViewCode = viewCode
	}
}

// 设置字段宽度
func SetColumnWidth(columnWidth int16) FuncSetViews {
	return func(svo *SetViewOptions) {
		svo.ColumnWidth = columnWidth
	}
}

// 设置字段高度
func SetColumnHeight(columnHeight int16) FuncSetViews {
	return func(svo *SetViewOptions) {
		svo.ColumnHeight = columnHeight
	}
}

// 设置字段是否隐藏
func SetColumnIsHidden(columnIsHidden int8) FuncSetViews {
	return func(svo *SetViewOptions) {
		svo.ColumnIsHidden = columnIsHidden
	}
}

// 设置字段是否排序
func SetColumnIsOrder(columnIsOrder int8) FuncSetViews {
	return func(svo *SetViewOptions) {
		svo.ColumnIsOrder = columnIsOrder
	}
}

// 设置字段是否过滤
func SetColumnFilter(columnFilter string) FuncSetViews {
	return func(svo *SetViewOptions) {
		svo.ColumnFilter = columnFilter
	}
}

// 设置字段固定, left, right, none
func SetColumnFixed(columnFixed string) FuncSetViews {
	return func(svo *SetViewOptions) {
		svo.ColumnFixed = columnFixed
	}
}
