/*
 * @Author: reel
 * @Date: 2023-06-18 10:32:27
 * @LastEditors: reel
 * @LastEditTime: 2023-09-09 19:36:56
 * @Description: 入口函数, 用于测试和示例, 不作为项目使用
 */
package main

import (
	"fmt"
	"os"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/gin-gonic/gin"
)

type params struct {
	ID string `json:"id" binding:"required"`
}

func main() {
	c, err := core.New()
	if err != nil {
		fmt.Println("初始化失败, 错误:", err)
		os.Exit(2)
	}

	ajax := c.Group("ajax")

	dim := ajax.Group("dim", "字典数据")

	pkl := dim.Group("picklist", "码值表")

	pkl.GET("list", "获取码值列表", params{}, func(ctx core.Context) {
		ctx.JSON(errno.ERRNO_OK.ToMapWithData("请求成功"))

	})
	c.Engine().GET("/", func(ctx *gin.Context) { ctx.JSON(200, errno.ERRNO_OK.ToMapWithData("请求成功")) })
	c.Run()

}
