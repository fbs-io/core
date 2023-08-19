/*
 * @Author: reel
 * @Date: 2023-06-16 05:57:22
 * @LastEditors: reel
 * @LastEditTime: 2023-08-19 13:38:06
 * @Description: 系统资源model, 用于管理API及菜单
 */
package core

import (
	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

var (
	sources    = make([]*Sources, 0, 100)
	sourcesMap = make(map[string]*Sources, 100)

	// 只有当资源类型为菜单或按钮, 可以用作权限设置
	// 不受限的资源, 用户登陆后都可以访问
	SOURCE_TYPE_LIMITED    int8 = 0 // 受限, 前端不可访问
	SOURCE_TYPE_MENU       int8 = 1 // 菜单权限
	SOURCE_TYPE_PERMISSION int8 = 2 // 按钮/接口权限
	SOURCE_TYPE_UNLIMITED  int8 = 3 // 不受限菜单或按钮/接口
)

// 系统资源表
//
// 用于API文档, 菜单, 权限控制等
//
// 当使用core中的路由接口生成路由时, 系统资源会自动注册到这张表中
type SourcesBase struct {
	SourceCode    string `json:"source_code" gorm:"comment:资源代码;uniqueIndex"`           // 资源code
	SourceName    string `json:"source_name" gorm:"comment:资源名称;index"`                 // 资源名称
	SourceDesc    string `json:"source_desc" gorm:"comment:资源说明"`                       // 资源描述
	SourcePCode   string `json:"source_parent_code" gorm:"comment:上层资源code;index"`      // 父级code
	SourceDeep    int8   `json:"source_deep" gorm:"comment:资源层级;index"`                 // 层级, 方便定位数据
	SourceView    string `json:"source_view" gorm:"comment:资源视图"`                       // 前端路由视图
	SourceIcon    string `json:"source_icon" gorm:"comment:资源图标"`                       // 前端菜单图标
	SourceSort    string `json:"source_sort" gorm:"comment:资源排序"`                       // 前端菜单顺序
	SourcePath    string `json:"source_path" gorm:"comment:资源路径;index"`                 // 层级全路径
	SourceType    int8   `json:"source_type" gorm:"comment:资源类型,0表示都可以显示, 1表示受限;index"` // 用于区分资源类型, 可以设置那些是用做权限配置的
	ViewType      string `json:"view_type" gorm:"comment:显示类型"`                         // 定义前端显示的类型, 如menu, button等
	ApiMethod     string `json:"api_method" gorm:"comment:后台接口方法"`                      // api接口路径
	ApiName       string `json:"api_name" gorm:"comment:后台接口名称"`                        // 对应后台接口名称
	ApiDesc       string `json:"api_desc" gorm:"comment:后台接口说明"`                        // 对应后台接口路径
	ApiParams     string `json:"api_params" gorm:"comment:前端请求参数"`                      // db中存储的参数字符串
	ApiAcceptType string `json:"accept_type" gorm:"comment:前端请求参数类型"`                   // 约束接口传参方式
}

// 数据库字段
//
// 对SourcesBase进行的封装
type Sources struct {
	SourcesBase
	rdb.Model
}

func (a *Sources) TableName() string {
	return "e_sys_sources"
}

func (a *Sources) BeforeCreate(tx *gorm.DB) error {
	a.Model.BeforeCreate(tx)
	return nil
}

// 用于外部设置souces
func (a *Sources) WithPermission(t int8) *Sources {
	a.SourceType = t
	return a
}

// 用于外部设置souces
func (a *Sources) SetDescription(des string) *Sources {
	a.SourceDesc = des
	return a
}
