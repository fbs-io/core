/*
 * @Author: reel
 * @Date: 2023-05-26 07:44:09
 * @LastEditors: reel
 * @LastEditTime: 2023-06-23 11:16:53
 * @Description: 静态资源
 */
package views

import "embed"

var (
    //go:embed install
    Install embed.FS

    //go:embed mscui
    Mscui embed.FS
)
