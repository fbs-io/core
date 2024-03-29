/*
 * @Author: reel
 * @Date: 2023-06-04 22:29:07
 * @LastEditors: reel
 * @LastEditTime: 2023-12-30 11:21:29
 * @Description: 请填写简介
 */
package consts

const (
	// 服务相关
	SERVER_IS_RUN  = 1  // 服务正常运行
	SERVER_IS_DOWN = 0  // 服务关闭
	SERVER_IS_NULL = -1 // 服务不可用

	// 上下文相关
	// 上下文的用户
	CTX_AUTH                = "ctx_auth"
	CTX_SHARDING_KEY        = "ctx_sharding_key"
	CTX_DATA_PERMISSION_KEY = "ctx_data_permission"
)
