/*
 * @Author: reel
 * @Date: 2023-05-11 23:04:14
 * @LastEditors: reel
 * @LastEditTime: 2023-08-15 23:23:39
 * @Description: 用于管理多个日志
 */

package logx

import "github.com/fbs-io/core/pkg/errorx"

var (
	Sys Logger
	DB  *gormLogger
	APP Logger
)

// 仅用于日志路径的设置, 请不要设置日志名称
func Init(optF ...optFunc) (err error) {
	sysOptfs := append(optF, SetLogPath("sys"), SetLogName("sys.log"))
	Sys, err = New(sysOptfs...)
	if err != nil {
		return errorx.Wrap(err, "系统日志配置发生错误")
	}
	DBOptfs := append(optF, SetLogPath("db"), SetLogName("db.log"))
	DB, err = NewGormLogger(DBOptfs...)
	if err != nil {
		return errorx.Wrap(err, "DB日志配置发生错误")
	}
	APPOptfs := append(optF, SetLogPath("app"), SetLogName("app.log"))
	APP, err = New(APPOptfs...)
	if err != nil {
		return errorx.Wrap(err, "应用日志配置发生错误")
	}
	return nil
}

func Close() {
	Sys.Close()
	DB.Close()
	APP.Close()
}
