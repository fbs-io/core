/*
 * @Author: reel
 * @Date: 2023-05-11 23:32:22
 * @LastEditors: reel
 * @LastEditTime: 2023-06-05 07:56:02
 * @Description: 管理接口
 */
package core

import (
    "fmt"
    "os"
    "os/signal"

    "github.com/fbs-io/core/internal/pem"
    "github.com/fbs-io/core/logx"
    "github.com/fbs-io/core/service"
)

// 关闭整个服务
func (c *core) Shutdown() {
    service.Stop()
    c.msc.Stop()
    if c.config.IsLoad {
        logx.Sys.Info(fmt.Sprintf("[%s]服务关闭", c.msc.Name()))
        logx.Close()
    }
}

func (c *core) Run() {
    // 加载系统启动配置
    if _, err := pem.GetPems(); err == nil {
        err = c.config.Load()
        if err != nil {
            fmt.Println("系统加载配置失败", err)
            os.Exit(2)
        }

        // 配置服务
        err := c.install()
        if err != nil {
            fmt.Println("系统初始化服务失败", err)
            os.Exit(2)
        }

        // 服务启动
        err = service.Start()
        if err != nil {
            fmt.Println("系统启动失败", err)
            os.Exit(2)
        }
    }

    // 等待指令用于关机
    quit := make(chan os.Signal, 1)

    signal.Notify(quit, os.Interrupt)
    <-quit

    c.Shutdown()
}
