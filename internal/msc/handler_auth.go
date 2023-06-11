/*
 * @Author: reel
 * @Date: 2023-06-04 18:45:20
 * @LastEditors: reel
 * @LastEditTime: 2023-06-10 07:42:57
 * @Description: 登陆相关
 */
package msc

import (
    "github.com/fbs-io/core/pkg/errno"
    "github.com/fbs-io/core/session"
    "github.com/gin-gonic/gin"
)

type loginParams struct {
    User string `json:"username" binding:"required"`
    Pwd  string `json:"password" binding:"required"`
}

// 登录时, 返回用户相关信息
// 如操作菜单, 权限等
func (m *handler) login() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        var p = &loginParams{}
        err := ctx.ShouldBindJSON(&p)
        if err != nil {
            ctx.JSON(200, errno.ERRNO_PARAMS_BIND.ToMapWithError(err))
            return
        }
        if m.config.User != p.User || m.config.Pwd != p.Pwd {
            ctx.JSON(200, errno.ERRNO_AUTH_USER_OR_PWD.ToMap())
            return
        }
        sessionKey := session.GenSessionKey()
        m.session.SetWithToken(sessionKey, m.config.User)
        result := map[string]interface{}{
            "token": sessionKey,
            "userInfo": map[string]interface{}{
                "dashboard": "0",
                "role":      []string{"admin"},
                "userId":    1,
                "userName":  "root",
            },
            "menu": menu,
            "permissions": []string{
                "list.add",
                "list.edit",
                "list.delete",
                "user.add",
                "user.edit",
                "user.delete",
            },
        }

        ctx.JSON(200, errno.ERRNO_OK.ToMapWithData(result))
    }
}
