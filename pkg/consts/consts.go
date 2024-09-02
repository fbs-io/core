/*
 * @Author: reel
 * @Date: 2023-06-04 22:29:07
 * @LastEditors: reel
 * @LastEditTime: 2024-08-21 07:32:43
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
	CTX_AUTH = "ctx_auth"
	// 上下文中分区标识
	CTX_SHARDING_KEY = "ctx_sharding_key"
	// 上下文中权限标识
	CTX_DATA_PERMISSION_KEY = "ctx_data_permission"
	// 上下文中分区DB标识
	CTX_SHARDING_DB = "ctx_sharding_db"
	// 上下文中的请求id
	CTX_TRACE_ID = "ctx_trace_id"
	// 链路追踪key
	TRACE_ID = "trace_id"
	// 链路追踪spankey
	SPAN_ID = "span_id"
	// 请求头的请求id
	REQUEST_HEADER_TRACE_ID = "TRACE-ID"

	// 日期格式化

	// 完整的年月
	DATE_FORMAT_Y      = "YYYY"
	DATE_FORMAT_YM     = "YYYYMM"
	DATE_FORMAT_YMD    = "YYYYMMDD"
	DATE_FORMAT_YMDH   = "YYYYMMDDhh"
	DATE_FORMAT_YMDHM  = "YYYYMMDDhhmm"
	DATE_FORMAT_YMDHMS = "YYYYMMDDHHmmss"

	// 简写年月
	DATE_FORMAT_SHORT_Y      = "YY"
	DATE_FORMAT_SHORT_YM     = "YYMM"
	DATE_FORMAT_SHORT_YMD    = "YYMMDD"
	DATE_FORMAT_SHORT_YMDH   = "YYMMDDhh"
	DATE_FORMAT_SHORT_YMDHM  = "YYMMDDhhmm"
	DATE_FORMAT_SHORT_YMDHMS = "YYMMDDHHmmss"
)
