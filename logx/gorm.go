/*
 * @Author: reel
 * @Date: 2023-08-14 20:39:55
 * @LastEditors: reel
<<<<<<< HEAD
 * @LastEditTime: 2023-08-16 21:28:17
=======
 * @LastEditTime: 2024-01-14 17:38:07
>>>>>>> develop
 * @Description: 适用于gorm的logger
 */
package logx

import (
	"context"
	"errors"
	"os"
	"path"
	"time"

	"github.com/fbs-io/core/pkg/env"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/pkg/filex"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type gormLogger struct {
	log           *logrus.Logger
	file          *os.File
	slowThreshold time.Duration
}

func NewGormLogger(optF ...optFunc) (*gormLogger, error) {

	var o = &opts{
		logName:         "db.log",
		logPath:         "db",
		timestampFormat: "2006-01-02 15:04:05.000",
		logMaxage:       365,
		logRotationTime: 7,
		slowThreshold:   5 * time.Second,
	}
	if env.Active() == nil {
		o.level = logrus.InfoLevel
	} else {
		o.level = env.Active().Level()
	}
	for _, opt := range optF {
		opt(o)
	}
	o.logPath = path.Join(o.dataPath, "logs", o.logPath)
	err := filex.CreatDir(o.logPath)
	if err != nil {
		return nil, errorx.Wrap(err, "日志路径创建/打开错误")
	}
	logFullPath := path.Join(o.logPath, o.logName)
	file, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errorx.Wrap(err, "日志文件创建/打开错误")
	}

	log := logrus.New()
	// log.Formatter = &logrus.JSONFormatter{TimestampFormat: o.timestampFormat}
	log.SetLevel(o.level)

	log.Out = os.Stdout
	if file != nil {
		log.Out = file
	}
	// 分割文件
	logWriter, _ := rotatelogs.New(
		// 分割后的文件名称
		logFullPath+".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(logFullPath),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(o.logMaxage*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(o.logRotationTime*24*time.Hour),
	)

	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: o.timestampFormat,
	})
	log.AddHook(lfHook)
	return &gormLogger{
		log:           log,
		file:          file,
		slowThreshold: o.slowThreshold,
	}, nil
}

func (l *gormLogger) LogMode(levle gormlogger.LogLevel) gormlogger.Interface {
	var logLevle logrus.Level

	switch levle {
	case gormlogger.Info:
		logLevle = logrus.InfoLevel
	case gormlogger.Error:
		logLevle = logrus.ErrorLevel
	case gormlogger.Warn:
		logLevle = logrus.ErrorLevel
	}
	if env.Active().Value() == env.ENV_MODE_DEV {
		logLevle = logrus.DebugLevel
	}
	l.log.SetLevel(logLevle)
	return l
}

func (l *gormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Infof(s, args...)
}

func (l *gormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Warnf(s, args...)
}

func (l *gormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Errorf(s, args...)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	// if l.SourceField != "" {
	// 	fields[l.SourceField] = utils.FileWithLineNum()
	// }
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		fields[logrus.ErrorKey] = err
		l.log.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		l.log.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	if env.Active().Value() == env.ENV_MODE_DEV {
		l.log.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}

func (l *gormLogger) Close() {
	l.file.Close()
}
