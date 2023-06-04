/*
 * @Author: reel
 * @Date: 2023-05-26 07:44:09
 * @LastEditors: reel
 * @LastEditTime: 2023-05-28 14:42:18
 * @Description: 静态资源
 */
package views

import "embed"

var (
	//go:embed static
	Static embed.FS
)
