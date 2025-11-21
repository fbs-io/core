/*
 * @Author: reel
 * @Date: 2023-07-23 22:01:29
 * @LastEditors: reel
 * @LastEditTime: 2025-11-08 18:47:58
 * @Description: 初始化core配置
 */

package core

import (
	"encoding/json"
	"fmt"
)

type options struct {
	limitSize     int    // 最多存储的令牌个数
	limitNumber   int    // 每秒生成的令牌个数
	appName       string // 设置应用名称
	appVersion    string // 设置应用版本
	shardingModel int8   // 分区模式
	sessionTTL    int    // session过期时间
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

func SetAppName(appName string) FuncCores {
	return func(opt *options) {
		opt.appName = appName
	}
}

// 设置应用版本信息
func SetAppVersion(appVersion string) FuncCores {
	return func(opt *options) {
		opt.appVersion = appVersion
	}
}

// 设置分区模式
func SetShardingModel(model int8) FuncCores {
	return func(opt *options) {
		opt.shardingModel = model
	}
}

// 设置 session过期时间
func SetSessionTTL(ttl int) FuncCores {
	return func(opt *options) {
		opt.sessionTTL = ttl
	}
}

type operateOpt struct {
	content string // 操作日志业务内容
	result  any    // 操作结果
	isSet   bool   // 设置操作日志
	params  string // 参数
	user    string // 用户
}

type FuncOperateOpt func(*operateOpt)

// 设置 业务内容
func SetContent(content string) FuncOperateOpt {
	return func(oo *operateOpt) {
		oo.content = content
	}
}

// 设置 返回结果
func SetResult(result any) FuncOperateOpt {
	return func(oo *operateOpt) {
		oo.result = result
	}
}
func SetResultByID(result any) FuncOperateOpt {
	return func(oo *operateOpt) {
		oo.result = fmt.Sprintf("id: %v", result)
	}
}

func SetResultByCode(result any) FuncOperateOpt {
	return func(oo *operateOpt) {
		oo.result = fmt.Sprintf("code: %v", result)
	}
}
func SetUser(user string) FuncOperateOpt {
	return func(oo *operateOpt) {
		oo.user = user
	}
}
func SetParams(params any) FuncOperateOpt {
	return func(oo *operateOpt) {
		if params == nil {
			return
		}
		paramsB, err := json.Marshal(params)
		if err != nil {
			oo.params = fmt.Sprintf("params: %v", params)
			return
		}
		oo.params = string(paramsB)
	}
}

// 设置 操作日志
func IsSetLog() FuncOperateOpt {
	return func(oo *operateOpt) {
		oo.isSet = true
	}
}
