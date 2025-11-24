/*
 * @Author: reel
 * @Date: 2023-06-16 05:57:22
 * @LastEditors: reel
 * @LastEditTime: 2025-11-24 22:06:08
 * @Description: 系统资源model, 用于管理API及菜单
 */
package core

import (
	"fmt"
	"strings"

	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

var (
	resources    = make([]*Resources, 0, 100)
	resourcesMap = make(map[string]*Resources, 100)

	// 只有当资源类型为菜单或按钮, 可以用作权限设置
	// 不受限的资源, 用户登陆后都可以访问

)

const (
	SOURCE_TYPE_LIMITED      int8 = iota // 受限, 前端不可访问
	SOURCE_TYPE_UNLIMITED                // 不受限菜单或按钮/接口
	SOURCE_TYPE_MENU                     // 菜单权限
	SOURCE_TYPE_UNMENU                   // 不受限菜单
	SOURCE_TYPE_PERMISSION               // 按钮/接口权限
	SOURCE_TYPE_UNPERMISSION             // 不受限按钮/接口

	CONST_MENU   = "menu"
	CONST_BUTTON = "button"

	// 判断是否时路由
	SOURCE_ROUTER_NAN int8 = 0 // 不返回路由
	SOURCE_ROUTER_IS  int8 = 1 // 返回路由, 默认菜单返回路由, 也可以单独设置按钮作为路由返回

	// 表名
	// 资源表
	TABLE_SYSTEM_CORE_RESOURCE = "e_sys_core_resources"
	// 操作日志表
	TABLE_SYSTEM_CORE_OPERATELOG = "e_sys_core_operatelog"
	// 前端视图字段表
	TABLE_SYSTEM_CORE_VIEWS = "e_sys_core_views"
)

// 系统资源表
//
// 用于API文档, 菜单, 权限控制等
//
// 当使用core中的路由接口生成路由时, 系统资源会自动注册到这张表中
type ResourcesBase struct {
	Code  string `json:"code" gorm:"column:resource_code;comment:资源代码;index"`                 // 资源code
	Name  string `json:"name" gorm:"column:resource_name;comment:资源名称;index"`                 // 资源名称,
	Desc  string `json:"desc" gorm:"column:resource_desc;comment:资源说明"`                       // 资源描述,可用作title
	PCode string `json:"pcode" gorm:"column:resource_pcode;comment:上层资源code;index"`           // 父级code
	Level int8   `json:"level" gorm:"column:resource_level;comment:资源层级;index"`               // 层级, 方便定位数据
	Api   string `json:"api" gorm:"column:resource_api;comment:资源路径;index"`                   // 资源访问api
	Type  int8   `json:"type" gorm:"column:resource_type;comment:资源类型,0表示都可以显示, 1表示受限;index"` // 用于区分资源类型, 可以设置那些是用做权限配置的
	Sort  string `json:"sort" gorm:"column:resource_sort;comment:资源排序"`                       // 前端菜单顺序
	// API文档用, 请求方法
	Method     string `json:"method" gorm:"column:resource_method;comment:后台接口方法"`             // api接口路径
	Params     string `json:"params" gorm:"column:resource_params;comment:前端请求参数"`             // db中存储的参数字符串
	AcceptType string `json:"accept_type" gorm:"column:resource_accept_type;comment:前端请求参数类型"` // 约束接口传参方式
	// 前端路由菜单用
	IsRouter  int8            `json:"is_router" gorm:"column:resource_is_router;comment:前端用路由判断;index"`        // 主要用于某些button需要展示路由上
	Path      string          `json:"path" gorm:"column:resource_path;comment:前端用路径;index"`                    // 前端用组件方法
	Component string          `json:"component" gorm:"column:resource_component;comment:组件名称"`                 // 前端组件名称
	Meta      rdb.ModeMapJson `json:"meta" gorm:"column:resource_meta;type:varchar(10000);comment:前端用路由参数元信息"` // 前端组件原信息
	PageView  rdb.ModeMapJson `json:"page_view" gorm:"column:resource_views;type:varchar(10000);comment:前端视图配置信息"`
}

// 数据库字段
//
// 对SourcesBase进行的封装
type Resources struct {
	ResourcesBase
	rdb.Model
	rdb.ShardingModel
	Children []*Resources `json:"children" gorm:"-"`
}

func (s *Resources) TableName() string {
	return TABLE_SYSTEM_CORE_RESOURCE
}

func (s *Resources) BeforeCreate(tx *gorm.DB) error {
	s.Model.BeforeCreate(tx)
	return nil
}

// 用于外部设置souces, 请通过 core.SOURCE_TYPE_* 进行设置
//
// 0: 受限资源, 无法访问, 该资源下在子集, 自动去除菜单和组件的前缀; 1:不受限, api和菜单都可访问
//
// 2: 受限菜单, 可通过权限设置访问; 3: 不受限菜单, 登陆用户均可访问
//
// 4: 受限api, 可通过权限设置访问; 5: 不受限api, 登陆用户均可访问
func (s *Resources) WithPermission(t int8) *Resources {
	s.Type = t
	switch t {
	case SOURCE_TYPE_MENU, SOURCE_TYPE_UNMENU:
		if s.Meta != nil {
			s.Meta["type"] = CONST_MENU
		}
	case SOURCE_TYPE_PERMISSION, SOURCE_TYPE_UNPERMISSION:
		if s.Meta != nil {
			s.Meta["type"] = CONST_BUTTON
		}
		s.IsRouter = SOURCE_ROUTER_NAN
	case SOURCE_TYPE_LIMITED:
		s.IsRouter = SOURCE_ROUTER_NAN
		s.Path = ""
		s.Component = ""
	case SOURCE_TYPE_UNLIMITED:
		if s.Method == "" {
			s.IsRouter = SOURCE_ROUTER_IS
		}
	}
	return s
}

// 用于外部设置路由
func (s *Resources) WithRouter(t int8) *Resources {
	s.IsRouter = t
	return s
}

// 用于外部设置souces
func (s *Resources) SetDescription(des string) *Resources {
	s.Desc = des
	return s
}

// 拼接请求参数和路由
//
// 主要用于权限校验
func (s *Resources) GenRequestKey() string {
	return fmt.Sprintf("%s:%s", strings.ToUpper(s.Method), s.Api)
}

// 设置允许通过的登陆签名校验的接口
//
// 默认所有接口需要签名校验
//
// 通过该方法可以设置例外接口, 如登陆接口
func (s *Resources) WithAllowSignature() *Resources {
	AddAllowSource(s.GenRequestKey())
	return s
}

// 去除前端菜单路由前缀
//
// 例如api接口转为前端路由: /ajax/user/list => /user/list
func (s *Resources) WithMenuNotPrefix(prefix string) *Resources {
	s.Path = strings.Replace(s.Path, "/"+prefix, "", -1)
	s.Component = strings.Replace(s.Component, prefix+"/", "", -1)
	return s
}

// 设置路由隐藏
func (s *Resources) WithHidden() *Resources {
	s.Meta["hidden"] = true
	return s
}

// 设置前端Meta信息
func (s *Resources) WithMeta(key string, value interface{}) *Resources {
	s.Meta[key] = value
	return s
}

func (e *Resources) ColumnNameWithCode() string {
	return "resource_code"
}

func (e *Resources) ParentCode() string {
	return e.PCode
}

// 设置前端页面信息
//
// 整个页面级别的设置, 如宽度等

type FuncSetPageViews func(options rdb.ModeMapJson)

// 设置视图元素宽度
func SetPageViewWidth(width int16) FuncSetPageViews {
	return func(options rdb.ModeMapJson) {
		options["ViewWidth"] = width
	}
}

// 设置table用于选择的key, 前端默认id
func SetPageTableKey(key string) FuncSetPageViews {
	return func(options rdb.ModeMapJson) {
		options["TableKey"] = key
	}
}

// 用于控制显示的列
func SetColumnsShow(columns []string) FuncSetPageViews {
	return func(options rdb.ModeMapJson) {
		options["ColumnsShow"] = columns
	}
}

// 用于控制隐藏的列
func SetColumnHidden(columns []string) FuncSetPageViews {
	return func(options rdb.ModeMapJson) {
		options["ColumnsHidden"] = columns
	}
}

func (e *Resources) SetPageViews(fns ...FuncSetPageViews) *Resources {
	view := make(rdb.ModeMapJson, 100)
	for _, fn := range fns {
		fn(view)
	}
	e.PageView = view
	return e
}

type OperateLog struct {
	IP        string `json:"ip" gorm:"操作ip;index"`
	User      string `json:"oper" gorm:"comment:操作用户;index"`
	Content   string `json:"content" gorm:"comment:业务操作内容"`
	Result    string `json:"result" gorm:"comment:结果"`
	Params    string `json:"params" gorm:"comment:请求参数"`
	Method    string `json:"method" gorm:"comment:请求方法;index"`
	Api       string `json:"api" gorm:"comment:操作接口;index"`
	ApiName   string `json:"api_name" gorm:"comment:接口名称"`
	TraceID   string `json:"trace_id" gorm:"comment:链路id;index"`
	OperateID string `json:"operate_id" gorm:"comment:操作id;index"`   // 部分页面会增加重复提交id
	CreatedAT uint   `json:"created_at" gorm:"autoCreateTime:milli"` // 创建时间
	rdb.ShardingModel
}

func (o *OperateLog) TableName() string {
	return TABLE_SYSTEM_CORE_OPERATELOG
}

// 用于前端展示的视图
type Views struct {
	ResourceCode string           `json:"resource_code" gorm:"column:resource_code;comment:资源code;index"`
	ViewCode     string           `json:"view_code" gorm:"column:view_code;comment:视图code;index"`
	ViewType     string           `json:"view_type" gorm:"column:view_type;comment:视图类型,如表格, 表单"`
	ViewItemCode string           `json:"view_item_code" gorm:"column:view_item_code;comment:视图项code;index"`
	ViewRole     string           `json:"view_role" gorm:"column:view_role;comment:视图功能,如新增, 编辑, 删除"`
	Code         string           `json:"code" gorm:"column:code;comment:字段code;index"`
	Name         string           `json:"name" gorm:"column:name;comment:字段名称"`
	Key          string           `json:"key" gorm:"column:key;comment:主键"`
	Calc         string           `json:"calc" gorm:"column:calc;comment:计算字段,用于数字计算"`
	ValueType    string           `json:"value_type" gorm:"column:value_type;comment:字段值类型,如string, int, float"`
	Width        int16            `json:"width" gorm:"column:width;comment:宽度"`
	Height       int16            `json:"height" gorm:"column:height;comment:高度"`
	Hidden       int8             `json:"hidden" gorm:"column:hidden;comment:是否隐藏"`
	IsOrder      int8             `json:"order" gorm:"column:isorder;comment:是否排序"`
	Filter       string           `json:"filter" gorm:"column:filter;comment:过滤字段"`
	Fixed        string           `json:"fixed" gorm:"column:fixed;comment:固定列"`
	Align        string           `json:"align" gorm:"column:align;comment:对齐方式,如左:left, 右:right, 居中:center"`
	FormatType   string           `json:"format_type" gorm:"column:format_type;comment:格式化类型,如时间:time, 金额:money,下拉选项等"`
	Format       string           `json:"format" gorm:"column:format;comment:格式化字段名"`
	Disabled     string           `json:"disabled" gorm:"column:disabled;comment:表单是否禁用编辑"`
	Multiple     string           `json:"multiple" gorm:"column:props;comment:是否多选"`
	Rules        string           `json:"rules" gorm:"column:rules;comment:表单组件校验规则"`
	Placeholder  string           `json:"placeholder" gorm:"column:placeholder;comment:表单组件提示信息"`
	Default      string           `json:"default" gorm:"column:default;comment:表单组件默认值"`
	Related      string           `json:"related" gorm:"column:related;comment:组件关联其他字段显示的默认值"`
	CustomValue  rdb.ModeListJson `json:"customValue" gorm:"type:string;column:custom_value;comment:表单组件自定义值"`
	Depend       string           `json:"depend" gorm:"column:depend;comment:依赖字段"`
	Account      string           `json:"account" gorm:"column:account;comment:账号"`
	rdb.Model
	rdb.ShardingModel
}

func (o *Views) TableName() string {
	return TABLE_SYSTEM_CORE_VIEWS
}

func (e *Views) ColumnNameWithCode() string {
	return "resource_code"
}
