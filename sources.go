/*
 * @Author: reel
 * @Date: 2023-06-16 05:57:22
 * @LastEditors: reel
 * @LastEditTime: 2023-09-12 06:52:26
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
	sources    = make([]*Sources, 0, 100)
	sourcesMap = make(map[string]*Sources, 100)

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
)

// 系统资源表
//
// 用于API文档, 菜单, 权限控制等
//
// 当使用core中的路由接口生成路由时, 系统资源会自动注册到这张表中
type SourcesBase struct {
	Code  string `gorm:"comment:资源代码;uniqueIndex"`           // 资源code
	Name  string `gorm:"comment:资源名称;index"`                 // 资源名称,
	Desc  string `gorm:"comment:资源说明"`                       // 资源描述,可用作title
	PCode string `gorm:"comment:上层资源code;index"`             // 父级code
	Level int8   `gorm:"comment:资源层级;index"`                 // 层级, 方便定位数据
	Api   string `gorm:"comment:资源路径;index"`                 // 资源访问api
	Type  int8   `gorm:"comment:资源类型,0表示都可以显示, 1表示受限;index"` // 用于区分资源类型, 可以设置那些是用做权限配置的
	Sort  string `gorm:"comment:资源排序"`                       // 前端菜单顺序
	// API文档用, 请求方法
	Method     string `gorm:"comment:后台接口方法"`   // api接口路径
	Params     string `gorm:"comment:前端请求参数"`   // db中存储的参数字符串
	AcceptType string `gorm:"comment:前端请求参数类型"` // 约束接口传参方式
	// 前端路由菜单用
	IsRouter  int8            `gorm:"comment:前端用路由判断;index"`                 // 主要用于某些button需要展示路由上
	Path      string          `gorm:"comment:前端用路径;index"`                   // 前端用组件方法
	Component string          `gorm:"comment:组件名称"`                          // 前端组件名称
	Meta      rdb.ModeMapJson `gorm:"type:varchar(1000);comment:前端用路由参数元信息"` // 前端组件原信息
}

// 数据库字段
//
// 对SourcesBase进行的封装
type Sources struct {
	SourcesBase
	rdb.Model
}

func (s *Sources) TableName() string {
	return "e_sys_sources"
}

func (s *Sources) BeforeCreate(tx *gorm.DB) error {
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
func (s *Sources) WithPermission(t int8) *Sources {
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
		s.IsRouter = SOURCE_ROUTER_IS
	}
	return s
}

// 用于外部设置路由
func (s *Sources) WithRouter(t int8) *Sources {
	s.IsRouter = t
	return s
}

// 用于外部设置souces
func (s *Sources) SetDescription(des string) *Sources {
	s.Desc = des
	return s
}

// 拼接请求参数和路由
//
// 主要用于权限校验
func (s *Sources) GenRequestKey() string {
	return fmt.Sprintf("%s:%s", strings.ToUpper(s.Method), s.Api)
}

// 设置允许通过的登陆签名校验的接口
//
// 默认所有接口需要签名校验
//
// 通过该方法可以设置例外接口, 如登陆接口
func (s *Sources) WithAllowSignature() *Sources {
	AddAllowSource(s.GenRequestKey())
	return s
}

// 去除前端菜单路由前缀
//
// 例如api接口转为前端路由: /ajax/user/list => /user/list
func (s *Sources) WithMenuNotPrefix(prefix string) *Sources {
	s.Path = strings.Replace(s.Path, "/"+prefix, "", -1)
	s.Component = strings.Replace(s.Component, prefix+"/", "", -1)
	return s
}

// 设置路由隐藏
func (s *Sources) WithHidden() *Sources {
	s.Meta["hidden"] = true
	return s
}

// 设置前端Meta信息
func (s *Sources) WithMeta(key string, value interface{}) *Sources {
	s.Meta[key] = value
	return s
}
