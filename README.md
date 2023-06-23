# FBS 系统内核

    基于Gin的二次封装, 对常用模块进行封装, 开发仅着重于业务层开发, 从而实现项目的快速开发和部署.

### 模块

* [X] 日志模块
* [X] 缓存模块
* [X] 数据库模块
* [X] 定时作业模块
* [X] HTTP服务模块
* [X] 系统管理中心模块
* [X] 系统资源管理模块
* [ ] 接口调试模块

### 特性

* 支持从Web页面进行系统配置, 包括数据库, 缓存, 端口, 超管用户等.
* 系统配置中心模块集成了服务器资源一栏, 后台服务管理, 方便管理员对系统进行操作
* 基于gin的路由封装, 可以自动完成简单的查询构造
* 自动生成接口一览, 及调试画面(开发中)
* 自动生成资源一览, 用于权限分配(开发中)
* 支持本地缓存和redis(待开发)
* 支持sqlite, mysql, postgresql 三款数据库

### 启动

```go
package main

import (
	"fmt"
	"os"

	"github.com/fbs-io/core"
)

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

    c.Run()

}
```
