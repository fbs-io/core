/*
 * @Author: reel
 * @Date: 2023-06-13 21:39:44
 * @LastEditors: reel
 * @LastEditTime: 2023-06-14 21:12:06
 * @Description: 内置相关服务信息
 */
package msc

import (
    "github.com/fbs-io/core/pkg/errno"
    "github.com/fbs-io/core/service"
    "github.com/gin-gonic/gin"
)

func (h *handler) getSrvStatus() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        data := append(h.srvStatus, map[string]interface{}{
            "service": "MSC",
            "node1":   1,
        })
        ctx.JSON(200, errno.ERRNO_OK.ToMapWithData(data))
    }
}

func (h *handler) setSrvRestart() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        go service.Restart()
        h.srvStatus = service.Status()
        ctx.JSON(200, errno.ERRNO_OK.ToMap())
    }
}
