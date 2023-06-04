/*
 * @Author: reel
 * @Date: 2023-06-04 22:37:35
 * @LastEditors: reel
 * @LastEditTime: 2023-06-04 22:37:37
 * @Description: 请填写简介
 */
package logx

import (
    "os"
    "path"
    "time"

    "core/pkg/env"
    "core/pkg/errorx"
    "core/pkg/filex"

    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    "github.com/rifflock/lfshook"
    "github.com/sirupsen/logrus"
)

var _ Logger = (*logger)(nil)

type logger struct {
    log  *logrus.Logger
    file *os.File
}

type Logger interface {
    loggerP()
    Close()
    Debug(msg string, infoF ...infoFunc)
    Info(msg string, infoF ...infoFunc)
    Warn(msg string, infoF ...infoFunc)
    Error(msg string, infoF ...infoFunc)
    Fatal(msg string, infoF ...infoFunc)
}

func New(optF ...optFunc) (Logger, error) {

    var o = &opts{
        logName:         "sys.log",
        logPath:         "sys",
        timestampFormat: "2006-01-02 15:04:05.000",
        logMaxage:       365,
        logRotationTime: 7,
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

    log.Formatter = &logrus.JSONFormatter{TimestampFormat: o.timestampFormat}
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
    return &logger{
        log:  log,
        file: file,
    }, nil
}

func (l *logger) loggerP() {}

func (l *logger) Close() {
    l.file.Close()
}

func (log *logger) Debug(msg string, infoF ...infoFunc) {
    infos := mkInfo(msg, infoF...)
    log.log.WithFields(infos.fields).Debug(infos.msg)
}

func (log *logger) Info(msg string, infoF ...infoFunc) {
    infos := mkInfo(msg, infoF...)
    log.log.WithFields(infos.fields).Info(infos.msg)
}

func (log *logger) Warn(msg string, infoF ...infoFunc) {
    infos := mkInfo(msg, infoF...)
    log.log.WithFields(infos.fields).Warn(infos.msg)
}

func (log *logger) Error(msg string, infoF ...infoFunc) {
    infos := mkInfo(msg, infoF...)
    log.log.WithFields(infos.fields).Error(infos.msg)
}

func (log *logger) Fatal(msg string, infoF ...infoFunc) {
    infos := mkInfo(msg, infoF...)
    log.log.WithFields(infos.fields).Fatal(infos.msg)
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
    msg    string
    fields logrus.Fields
}

type infoFunc func(*info)

func F(key string, value interface{}) infoFunc {
    return func(i *info) {
        if value == nil {
            i.fields[key] = nil
        } else {
            i.fields[key] = value
        }
    }
}

func EV(err error) infoFunc {
    return func(i *info) {
        if err == nil {
            i.fields["error"] = ""
        } else {
            i.fields["error"] = err.Error()
        }

    }
}

func E(err error) infoFunc {
    return func(i *info) {
        i.msg = err.Error()
    }
}

func mkInfo(msg string, infoF ...infoFunc) *info {
    var infos = &info{
        msg:    msg,
        fields: make(logrus.Fields),
    }
    for _, infof := range infoF {
        infof(infos)
    }
    return infos
}
