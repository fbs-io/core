/*
 * @Author: reel
 * @Date: 2023-12-30 17:09:47
 * @LastEditors: reel
 * @LastEditTime: 2023-12-30 17:36:08
 * @Description: 请填写简介
 */

package mscui

import "embed"

var (

	//go:embed website/*
	Static embed.FS

	//go:embed index.html
	Index []byte
)
