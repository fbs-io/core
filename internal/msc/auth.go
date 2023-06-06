/*
 * @Author: reel
 * @Date: 2023-06-04 18:45:20
 * @LastEditors: reel
 * @LastEditTime: 2023-06-04 20:00:00
 * @Description: 登陆相关
 */
package msc

import "github.com/gin-gonic/gin"

func (m *handler) login() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        ctx.Header("Content-Type", "text/html; charset=utf-8")
        index := indexIndex
        if !m.config.IsLoad {
            index = indexInstall
        }
        ctx.String(200, index)

    }
}

// func (m *handler) resetpwd() gin.HandlerFunc {
//     return func(ctx *gin.Context) {
//         ctx.Header("Content-Type", "text/html; charset=utf-8")
//         index := indexIndex
//         if !m.config.IsLoad {
//             index = indexInstall
//         }
//         ctx.String(200, index)
//     }
// }
