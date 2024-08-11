/*
 * @Author: reel
 * @Date: 2023-06-19 23:18:19
 * @LastEditors: reel
 * @LastEditTime: 2024-07-07 20:49:00
 * @Description: 测试路由相关方法
 */

package core

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRouter(t *testing.T) {

	// api参数生成测试
	type OrgCreateRequest struct {
		OrgName       string `json:"org_name" default:"总经理室" binding:"required" desc:"组织名称, 可重复"`
		OrgCode       string `json:"org_code" default:"001" binding:"required" desc:"组织代码, 唯一不可重复"`
		OrgParentCode string `json:"org_parent_code" default:""  desc:"组织上级 code"`
	}
	rt := reflect.TypeOf(OrgCreateRequest{})
	fmt.Println(genResourcesParams(rt))
	fmt.Println(genResourcesParams(nil))

	// 生产资源数据测试
	rout := &router{
		group: gin.New().Group("api"),
	}
	source := rout.genResources("api", "api", "")
	fmt.Println(source)

	// 模拟使用时生成资源表测试
	c, _ := New()
	api := c.Group("api")
	org := api.Group("org", "组织管理")
	org.GET("list", "查询组织列表", OrgCreateRequest{}, func(ctx Context) {})
	org.POST("list", "查询组织列表", OrgCreateRequest{}, func(ctx Context) {})
	org.PUT("list", "查询组织列表", OrgCreateRequest{}, func(ctx Context) {})
	org.DELETE("list", "查询组织列表", OrgCreateRequest{}, func(ctx Context) {})

	for _, s := range resourcesMap {
		fmt.Println(s)

	}

	// 测试请求参数是正常
	for k, v := range requestParams {
		fmt.Println("请求路径: ", k, ", 请求参数: ", v)
	}
}
