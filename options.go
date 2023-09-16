/*
 * @Author: reel
 * @Date: 2023-07-23 22:01:29
 * @LastEditors: reel
 * @LastEditTime: 2023-07-23 22:07:40
 * @Description: 初始化core配置
 */

package core

type options struct {
	limitSize   int // 最多存储的令牌个数
	limitNumber int // 每秒生成的令牌个数
}

type FuncCores func(*options)

// 设置令牌桶大小
func SetLimitSize(limitSize int) FuncCores {
	return func(opt *options) {
		opt.limitSize = limitSize
	}
}

// 设置每秒令牌生成个数
func SetLimitNumber(limitNumber int) FuncCores {
	return func(opt *options) {
		opt.limitNumber = limitNumber
	}
}
