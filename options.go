/*
 * @Author: reel
 * @Date: 2023-07-23 22:01:29
 * @LastEditors: reel
 * @LastEditTime: 2024-10-05 11:45:46
 * @Description: 初始化core配置
 */

package core

type options struct {
	limitSize     int    // 最多存储的令牌个数
	limitNumber   int    // 每秒生成的令牌个数
	appName       string // 设置应用名称
	appVersion    string // 设置应用版本
	shardingModel int8   // 分区模式
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

type operateOpt struct {
	content string      // 操作日志业务内容
	result  interface{} // 操作结果
	isSet   bool        // 设置操作日志

}

type FuncOperateOpt func(*operateOpt)

// 设置 业务内容
func SetContent(content string) FuncOperateOpt {
	return func(oo *operateOpt) {
		oo.content = content
	}
}

// 设置 返回结果
func SetResult(result interface{}) FuncOperateOpt {
	return func(oo *operateOpt) {
		oo.result = result
	}
}

// 设置 操作日志
func IsSetLog() FuncOperateOpt {
	return func(oo *operateOpt) {
		oo.isSet = true
	}
}
