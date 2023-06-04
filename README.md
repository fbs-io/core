# FBS 系统内核

    基于Gin的二次封装, 对常用模块进行封装, 开发仅着重于业务层开发, 从而实现项目的快速开发和部署.

### 集成模块

* [X] 日志模块
* [X] 缓存模块
* [X] 数据库模块
* [X] 定时作业模块
* [X] HTTPS服务模块
* [ ] 系统管理中心模块
* [ ] 系统资源管理模块
* [ ] 接口调试模块

### 内核特性

* 支持从Web页面进行系统配置, 包括数据库, 缓存, 端口, 超管用户等.
* 系统配置中心模块集成了服务器资源一栏, 后台服务管理, 方便管理员对系统进行操作
* 自动生成接口一览, 及调试画面
* 自动生成资源一览, 用于权限分配
* 支持本地缓存和redis(未来)
* 支持sqlite, mysql, postgresql 三款数据库, 计划支持读写分离, 多环境配置

### 快速启动

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
		fmt.Println("系统初始化失败", err)
		os.Exit(2)
	}

	err = c.Start()
	if err != nil {
		fmt.Println("系统启动失败", err)
		os.Exit(2)
	}

	c.Shutdown()

}
```
