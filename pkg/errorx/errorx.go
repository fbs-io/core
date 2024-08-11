/*
 * @Author: reel
 * @Date: 2023-06-04 22:09:34
 * @LastEditors: reel
 * @LastEditTime: 2024-07-07 18:22:36
 * @Description: 请填写简介
 */
package errorx

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

func callers() []uintptr {
	var pcs [32]uintptr
	l := runtime.Callers(3, pcs[:])
	return pcs[:l]
}

type Error interface {
	error
	ex()
	Format() string
}

var _ Error = (*errx)(nil)

type errx struct {
	msg   string
	stack []uintptr
}

func (i *errx) Error() string {
	return i.msg
}
func (i *errx) ex() {}

func (i *errx) Format() string {
	var msg string
	for _, pc := range i.stack {
		msg += fmt.Sprintf("%+v\n", errors.Frame(pc))
	}
	return msg
}

// 生成一个错误
func New(msg string) Error {
	return &errx{msg: msg}
}

func NewWithStack(msg string) Error {
	return &errx{msg: msg, stack: callers()}
}

// 格式化生成一个错误
func Errorf(format string, args ...interface{}) Error {
	return &errx{msg: fmt.Sprintf(format, args...)}
}

// func ErrorfWithStack(format string, args ...interface{}) Error {
// 	return &errx{msg: fmt.Sprintf(format, args...), stack: callers()}
// }

// 给错误添加额外信息
func Wrap(err error, msg string) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(*errx)
	if !ok {
		return &errx{msg: msg + "; " + err.Error()}
	}
	e.msg = msg + "; " + e.msg
	return e
}
func WrapWithStack(err error, msg string) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(*errx)
	if !ok {
		return &errx{msg: msg + "; " + err.Error(), stack: callers()}
	}
	e.msg = msg + "; " + e.msg
	return e
}

// 通过格式化, 给错误添加额外信息
// func WrapF(err error, format string, args ...interface{}) Error {
// 	if err == nil {
// 		return nil
// 	}
// 	msg := fmt.Sprintf(format, args...)
// 	e, ok := err.(*errx)
// 	if !ok {
// 		return &errx{msg: msg + "; " + err.Error()}
// 	}
// 	e.msg = msg + "; " + e.msg
// 	return e
// }

// func WrapFWithStack(err error, format string, args ...interface{}) Error {
// 	if err == nil {
// 		return nil
// 	}
// 	msg := fmt.Sprintf(format, args...)
// 	e, ok := err.(*errx)
// 	if !ok {
// 		return &errx{msg: msg + "; " + err.Error(), stack: callers()}
// 	}
// 	e.msg = msg + "; " + e.msg
// 	return e
// }

// 给错误添加堆栈信息
func WithStack(err error) Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*errx); ok {
		return e
	}
	return &errx{msg: err.Error(), stack: callers()}
}
