/*
 * @Author: reel
 * @Date: 2023-05-11 22:27:38
 * @LastEditors: reel
 * @LastEditTime: 2023-08-14 21:13:47
 * @Description: 请填写简介
 */
package logx

import (
	"time"

	"github.com/sirupsen/logrus"
)

type opts struct {
	dataPath        string
	logName         string
	logPath         string
	timestampFormat string
	level           logrus.Level
	logMaxage       time.Duration
	logRotationTime time.Duration
	slowThreshold   time.Duration
}

type optFunc func(*opts)

// 设置日志名称
func SetLogName(logName string) optFunc {
	return func(opt *opts) {
		opt.logName = logName
	}
}

// 设置日志存储目录
func SetDataPath(dataPath string) optFunc {
	return func(opt *opts) {
		opt.dataPath = dataPath
	}
}

// 设置日志路径
func SetLogPath(logPath string) optFunc {
	return func(o *opts) {
		o.logPath = logPath
	}
}

// 设置时间格式
func SetTimestampFormat(timestampFormat string) optFunc {
	return func(o *opts) {
		o.timestampFormat = timestampFormat
	}
}

// 设置 日志级别
func SetLevel(level logrus.Level) optFunc {
	return func(o *opts) {
		o.level = level
	}
}

// 设置最大保存时间
func SetLogMaxage(logMaxage time.Duration) optFunc {
	return func(o *opts) {
		o.logMaxage = logMaxage
	}
}

// 设置覆盖时间
func SetLogRotationTime(logRotationTime time.Duration) optFunc {
	return func(o *opts) {
		o.logRotationTime = logRotationTime
	}
}

// 设置慢查询时间
func SetLogSlowThreshold(t time.Duration) optFunc {
	return func(o *opts) {
		o.slowThreshold = t
	}
}
