/*
 * @Author: reel
 * @Date: 2023-06-04 22:37:35
 * @LastEditors: reel
 * @LastEditTime: 2024-10-07 01:31:34
 * @Description: 请填写简介
 */
package logx

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/fbs-io/core/pkg/consts"
	"github.com/fbs-io/core/pkg/env"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/pkg/filex"
	"github.com/fbs-io/core/pkg/trace"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var _ Logger = (*logger)(nil)

type logger struct {
	log           *logrus.Logger
	file          *os.File
	slowThreshold time.Duration
}

type Logger interface {
	loggerP()
	Close()
	Debug(msg string, infoF ...EntityFunc)
	Info(msg string, infoF ...EntityFunc)
	Warn(msg string, infoF ...EntityFunc)
	Error(msg string, infoF ...EntityFunc)
	Fatal(msg string, infoF ...EntityFunc)
}

func New(optF ...optFunc) (Logger, error) {

	var o = &opts{
		logName:         "sys.log",
		logPath:         "sys",
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
	lfHook := lfshook.NewHook(writerMap, &logrus.TextFormatter{
		TimestampFormat: o.timestampFormat,
	})
	log.AddHook(lfHook)
	return &logger{
		log:           log,
		file:          file,
		slowThreshold: o.slowThreshold,
	}, nil
}

func (l *logger) loggerP() {}

func (l *logger) Close() {
	l.file.Close()
}

func (log *logger) Debug(msg string, infoF ...EntityFunc) {
	entity := log.Entity(infoF...)
	entity.Debug(msg)
}

func (log *logger) Info(msg string, infoF ...EntityFunc) {
	entity := log.Entity(infoF...)
	entity.Info(msg)
}

func (log *logger) Warn(msg string, infoF ...EntityFunc) {
	entity := log.Entity(infoF...)
	entity.Warn(msg)
}

func (log *logger) Error(msg string, infoF ...EntityFunc) {
	entity := log.Entity(infoF...)
	entity.Error(msg)
}

func (log *logger) Fatal(msg string, infoF ...EntityFunc) {
	entity := log.Entity(infoF...)
	entity.Fatal(msg)
}

// gorm 日志使用
func (log *logger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		log.log.WithFields(logrus.Fields{
			"module":        "gorm",
			"type":          "sql",
			"src":           v[1],
			"duration":      v[2],
			"sql":           v[3],
			"values":        v[4],
			"rows_returned": v[5],
		}).Error("sql")
	case "log":
		log.log.WithFields(logrus.Fields{"gorm": v[2]}).Error("log")
	}
}

// 日志写入参数
type info struct {
	msg       string
	fields    logrus.Fields
	context   context.Context
	startTime time.Time
}

type EntityFunc func(*info)

func F(key string, value interface{}) EntityFunc {
	return func(i *info) {
		if value == nil {
			i.fields[key] = nil
		} else {
			i.fields[key] = value
		}
	}
}

// 错误信息
//
// 作为其中一个字段进行打印
func EV(err error) EntityFunc {
	return func(i *info) {
		if err != nil {
			i.fields["error"] = err.Error()
		}
	}
}

// 错误信息
//
// 打印在msg
func E(err error) EntityFunc {
	return func(i *info) {
		i.msg = err.Error()
	}
}

// 添加上下文
func Context(ctx context.Context) EntityFunc {
	return func(i *info) {
		i.context = ctx
	}
}

// 计算执行时间
func DiffTime(startTime time.Time) EntityFunc {
	return func(i *info) {
		i.startTime = startTime
	}
}

func (log *logger) Entity(infoF ...EntityFunc) (entity *logrus.Entry) {
	var infos = &info{
		fields:  make(logrus.Fields, 10),
		context: nil,
	}
	for _, infof := range infoF {
		infof(infos)
	}
	if infos.startTime.Unix() > 0 {
		infos.fields["diff_time"] = time.Since(infos.startTime)
	}
	if infos.context != nil {
		ti := infos.context.Value(consts.CTX_TRACE_ID)
		var tv *trace.Trace
		if ti != nil {
			tv = ti.(*trace.Trace)
			if tv.SpanID != "" {
				entity = log.log.
					WithFields(logrus.Fields{consts.TRACE_ID: tv.TraceID}).
					WithFields(logrus.Fields{consts.SPAN_ID: tv.SpanID}).
					WithFields(infos.fields)
				return
			}
			entity = log.log.WithFields(logrus.Fields{consts.TRACE_ID: tv.TraceID}).
				WithFields(infos.fields)
			return
		}
	}
	entity = log.log.WithFields(infos.fields)
	return
}

func Details(mags map[string]interface{}) EntityFunc {
	return func(i *info) {
		for k, v := range mags {
			i.fields[k] = v
		}
	}
}
