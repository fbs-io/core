package core

import (
    // "fbs/internal/core/apidoc"
    // "fbs/internal/core/means"
    // "fbs/pkg/convx"
    "encoding/json"
    "fmt"
    "path"
    "reflect"
    "strings"
    "sync"

    "github.com/gin-gonic/gin"
)

type HandlerFunc func(Context)

func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
    funcs := make([]gin.HandlerFunc, len(handlers))
    for i, h := range handlers {
        funcs[i] = func(c *gin.Context) {
            ctx := newCtx(c)
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
    Group(api, apiName string, handlers ...HandlerFunc) RouterGroup
    IRoutes
}

var _ IRoutes = (*router)(nil)

// IRoutes 包装gin的IRoutes
type IRoutes interface {
    // Any(string, ...HandlerFunc)
    //需要填写相对路由路径, 名称, 参数, 及中间件, 用于在 api 文档和菜单中注册
    //参数为如果为空, 该方法不会在 api 文档中进行注册
    GET(relativePath, pathName string, params interface{}, handlers ...HandlerFunc)
    //需要填写相对路由路径, 名称, 参数, 及中间件, 用于在 api 文档和菜单中注册
    //参数为如果为空, 该方法不会在 api 文档中进行注册
    PUT(relativePath, pathName string, params interface{}, handlers ...HandlerFunc)
    //需要填写相对路由路径, 名称, 参数, 及中间件, 用于在 api 文档和菜单中注册
    //参数为如果为空, 该方法不会在 api 文档中进行注册
    POST(relativePath, pathName string, params interface{}, handlers ...HandlerFunc)
    //需要填写相对路由路径, 名称, 参数, 及中间件, 用于在 api 文档和菜单中注册
    //参数为如果为空, 该方法不会在 api 文档中进行注册
    DELETE(relativePath, pathName string, params interface{}, handlers ...HandlerFunc)
    // TODO: 以后根据业务进行扩展
    // PATCH(string, ...HandlerFunc)
    // OPTIONS(string, ...HandlerFunc)
    // HEAD(string, ...HandlerFunc)
}

type router struct {
    group *gin.RouterGroup
}

var (
    routers = make(map[string]*router, 100)
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
    lock := &sync.Mutex{}
    lock.Lock()
    defer lock.Unlock()

    if routers[relativePath] != nil {
        return
    }
    if pathName == "" {
        pathName = relativePath
    }
    sources = append(sources, r.genSources(relativePath, pathName, "", nil))

    routers[relativePath] = r

}

// 顶层路由分组, 使用路由的入口
func (c *core) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
    r := getRouter(relativePath)
    if r == nil {
        r = &router{group: c.Engine().Group(relativePath, wrapHandlers(handlers...)...)}
    }
    setRouter(relativePath, relativePath, r)
    return r
}

// 对gin路由分组的封装
func (r *router) Group(relativePath, pathName string, handlers ...HandlerFunc) RouterGroup {
    group := r.group.Group(relativePath, wrapHandlers(handlers...)...)

    rout := getRouter(relativePath)
    if rout == nil {
        rout = &router{group: group}
    }
    setRouter(relativePath, pathName, rout)

    return rout
}

// Get请求方式封装
//
// 参数如果为空, 该方法不会被记录在资源表中
func (r *router) GET(relativePath, pathName string, params interface{}, handlers ...HandlerFunc) {
    // handlers = append([]HandlerFunc{r.validParams()}, handlers...)
    r.operation("GET", relativePath, pathName, params)
    r.group.GET(relativePath, wrapHandlers(handlers...)...)
}

// Post请求方式封装
//
// 参数如果为空, 该方法不会被记录在资源表中
func (r *router) POST(relativePath, pathName string, params interface{}, handlers ...HandlerFunc) {
    // handlers = append([]HandlerFunc{r.validParams()}, handlers...)
    r.operation("POST", relativePath, pathName, params)
    r.group.POST(relativePath, wrapHandlers(handlers...)...)
}

// Delete请求方式封装
//
// 参数如果为空, 该方法不会被记录在资源表中
func (r *router) DELETE(relativePath, pathName string, params interface{}, handlers ...HandlerFunc) {
    // handlers = append([]HandlerFunc{r.validParams()}, handlers...)
    r.operation("DELETE", relativePath, pathName, params)
    r.group.DELETE(relativePath, wrapHandlers(handlers...)...)
}

// Put请求方式封装
//
// 参数如果为空, 该方法不会被记录在资源表中
func (r *router) PUT(relativePath, pathName string, params interface{}, handlers ...HandlerFunc) {
    // handlers = append([]HandlerFunc{r.validParams()}, handlers...)
    r.operation("PUT", relativePath, pathName, params)
    r.group.PUT(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) PATCH(relativePath string, handlers ...HandlerFunc) {
    r.group.PATCH(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) OPTIONS(relativePath string, handlers ...HandlerFunc) {
    r.group.OPTIONS(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) HEAD(relativePath string, handlers ...HandlerFunc) {
    r.group.HEAD(relativePath, wrapHandlers(handlers...)...)
}

// 处理参数生成逻辑
func (r *router) operation(method, relativePath, pathName string, params interface{}) {
    if relativePath == "" {
        relativePath = "/"
    }
    if params != nil {
        rt := reflect.TypeOf(params)
        // 每个接口的参数存放在变量中便于后面查询使用
        requestParams[fmt.Sprintf("%s:%s/%s", method, r.group.BasePath(), relativePath)] = rt
        sources = append(sources, r.genSources(relativePath, pathName, method, rt))
    }
}

// 用于生成系统资源结构体
//
// 用于API文档和前端菜单
//
// 同时可用于权限及自动生成gorm查询参数
func (r *router) genSources(relativePath, name, method string, params reflect.Type) *Sources {

    basePaths := strings.Split(r.group.BasePath(), "/")[1:]
    method = strings.ToLower(method)
    s := &Sources{}
    fullpath := basePaths
    if method != "" {
        basePaths = append(basePaths, relativePath)
        fullpath = append([]string{method}, basePaths...)
        s.ApiMethod = method
        s.ApiName = name
        s.ApiParams, s.ApiAcceptType = genSourcesParams(params)
    }
    s.SourceName = name
    s.SourceCode = strings.Join(fullpath, ":")
    s.SourcePCode = strings.Join(basePaths[:len(basePaths)-1], ":")
    s.SourceDeep = int8(len(basePaths))
    s.ApiPath = path.Join(basePaths...)

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

    // 参数相关
    paramsKey        = "key"
    paramsValue      = "value"
    paramsValueType  = "value_type"
    paramsValueInt   = "int"
    paramsValueNum   = "number"
    paramsValueBool  = "bool"
    paramsValueFloat = "float"
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
func genSourcesParams(rt reflect.Type) (params string, contentType string) {
    if rt == nil {
        return
    }
    data := make([]interface{}, 0)

    for i := 0; i < rt.NumField(); i++ {
        item := make(map[string]interface{}, 4)
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

        // 用于前端API文档中的默认值
        item[paramsValue] = field.Tag.Get(tagDefault)
        // 用于字段描述
        item[tagDesc] = field.Tag.Get(tagDesc)
        // 用于校验参数信息
        item[tagBinding] = field.Tag.Get(tagBinding)

        data = append(data, item)
        paramsB, _ := json.Marshal(data)
        params = string(paramsB)
    }
    return
}
