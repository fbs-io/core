/*
 * @Author: reel
 * @Date: 2023-09-08 06:52:29
 * @LastEditors: reel
 * @LastEditTime: 2023-09-08 07:14:30
 * @Description: 请填写简介
 */
package mscui

import "embed"

var (
	//go:embed mscui
	Mscui embed.FS

	//go:embed mscui/index.html
	Index []byte
)
