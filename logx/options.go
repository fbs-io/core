/*
 * @Author: reel
 * @Date: 2023-05-11 22:27:38
 * @LastEditors: reel
 * @LastEditTime: 2023-05-16 20:38:01
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
}

type optFunc func(*opts)

func SetLogName(logName string) optFunc {
	return func(opt *opts) {
		opt.logName = logName
	}
}
func SetDataPath(dataPath string) optFunc {
	return func(opt *opts) {
		opt.dataPath = dataPath
	}
}

func SetLogPath(logPath string) optFunc {
	return func(o *opts) {
		o.logPath = logPath
	}
}

func SetTimestampFormat(timestampFormat string) optFunc {
	return func(o *opts) {
		o.timestampFormat = timestampFormat
	}
}
func SetLevel(level logrus.Level) optFunc {
	return func(o *opts) {
		o.level = level
	}
}

func SetLogMaxage(logMaxage time.Duration) optFunc {
	return func(o *opts) {
		o.logMaxage = logMaxage
	}
}

func SetLogRotationTime(logRotationTime time.Duration) optFunc {
	return func(o *opts) {
		o.logRotationTime = logRotationTime
	}
}
