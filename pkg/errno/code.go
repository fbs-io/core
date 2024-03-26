/*
 * @Author: reel
 * @Date: 2023-05-11 22:19:24
 * @LastEditors: reel
 * @LastEditTime: 2024-03-26 06:24:01
 * @Description: 定义常用的错误代码
 */
package errno

// 定义常用的错误代码

var (
	ERRNO_OK Errno = New(200, 0, "OK")

	// 系统相关代码
	// 系统内部通用错误代码
	ERRNO_SYSTEM            Errno = New(200, 50000, "系统内部错误")
	ERRNO_TOO_MANY_REQUESTS Errno = New(200, 50001, "请求次数过多,请稍后再试")
	// 系统初始化及授权相关
	ERRNO_INIT        Errno = New(200, 50001, "请先进行项目初始化")
	ERRNO_DISCLAIMERS Errno = New(200, 50002, "请先同意软件服务协议")
	ERRNO_IS_INSTALL  Errno = New(200, 50007, "已完成项目初始化, 请删除配置文件重新更新")
	// ERRNO_GEN_MACHINE_CODE  Errno = New(200, 50003, "无法生成机器码, 请联系管理员")
	// ERRNO_DEACTIVATE        Errno = New(200, 50004, "授权已过期, 请联系管理员")
	// ERRNO_EXCEEDING_ACCOUNT Errno = New(200, 50005, "账户已超过最大限制, 请联系管理员")
	// ERRNO_ACTIVe            Errno = New(200, 50006, "验证 license 错误")

	// 数据库错误相关
	ERRNO_RDB                Errno = New(200, 50100, "数据库通用错误")
	ERRNO_RDB_QUERY          Errno = New(200, 50101, "数据查询错误") //数据查询错误
	ERRNO_RDB_DELETE         Errno = New(200, 50102, "数据删除错误") //数据删除错误
	ERRNO_RDB_CREATE         Errno = New(200, 50103, "数据创建错误") //数据创建错误
	ERRNO_RDB_UPDATE         Errno = New(200, 50104, "数据更新错误") //数据更新错误
	ERRNO_RDB_DUPLICATED_KEY       = New(200, 50105, "重复主键错误") //数据重复错误
	// 缓存错误相关
	ERRNO_CACHE       Errno = New(200, 50201, "缓存错误")
	ERRNO_CACHE_QUERY Errno = New(200, 50201, "缓存查询错误")

	// 业务错误
	// 请求参数
	ERRNO_PARAMS_BIND    Errno = New(200, 10101, "参数绑定错误")
	ERRNO_PARAMS_INVALID Errno = New(200, 10102, "参数校验错误")
	// ERRNO_PARAMS_CREATE  Errno = New(200, 10103, "参数创建错误")

	// websocket 错误相关
	ERRNO_WS         Errno = New(200, 10200, "ws 错误")
	ERRNO_WS_REQUEST Errno = New(200, 10201, "创建 ws 连接失败")
	// ERRNO_PARAMS_BIND Errno = New(200, 10102, "参数错误")

	ERRNO_LOGIN_WITH_BROWSER Errno = New(200, 30101, "扫码登录失败")

	// 用户操作错误
	// 请求路径错误
	ERRNO_AUTH             Errno = New(401, 40000, "用户权限错误")
	ERRNO_AUTH_NOT_LOGIN   Errno = New(401, 40001, "用户未登陆")
	ERRNO_AUTH_PERMISSION  Errno = New(401, 40002, "用户无访问权限")
	ERRNO_AUTH_USER_OR_PWD Errno = New(200, 40004, "账号或密码错误")
	ERRNO_AUTH_ELSE_LOGIN  Errno = New(200, 40005, "已在其他地方登陆")

	// 权限错误
)
