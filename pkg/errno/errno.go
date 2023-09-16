/*
 * @Author: reel
 * @Date: 2023-05-11 22:19:24
 * @LastEditors: reel
 * @LastEditTime: 2023-08-26 21:53:28
 * @Description: API 常用状态码
 */
package errno

import "github.com/fbs-io/core/pkg/errorx"

var _ Errno = (*code)(nil)

type Errno interface {
	en()
	WithStack(err error) Errno
	HTTPCode() int
	code() int
	Message() string
	Stack() error
	ToMap() map[string]interface{}
	ToMapWithData(data interface{}) map[string]interface{}
	ToMapWithError(err error) map[string]interface{}
	WrapData(details interface{}) Errno
	WrapError(err error) Errno
	Notify() Errno
}

type code struct {
	httpCode   int         // http 状态码
	errno      int         // 业务状态码
	message    string      // 业务状态描述
	stackError error       // 含有堆栈的状态描述
	details    interface{} // 业务数据
	isStack    bool        // 是否加载堆栈信息
	isNotify   bool        // 是否通知
	isConst    bool        // 是否是常量
}

// 生成新的 code 暂时不需要调用错误堆栈错误, 需要时再调用
func New(httpCode, errno int, msg string) Errno {
	return &code{
		httpCode: httpCode,
		errno:    errno,
		message:  msg,
		isConst:  true,
	}
}

func (e *code) en() {}

func (e *code) WithStack(err error) Errno {
	e.isStack = true
	e.stackError = errorx.WithStack(err)
	return e
}

func (e *code) HTTPCode() int {
	return e.httpCode
}

func (e *code) code() int {
	return e.errno
}

func (e *code) Message() string {
	return e.message
}

func (e *code) Stack() error {
	return e.stackError
}

func (e *code) ToMap() (res map[string]interface{}) {
	res = map[string]interface{}{
		"errno":   e.errno,
		"message": e.message,
	}

	if e.details != nil {
		res["details"] = e.details
	}
	if e.isStack {
		res["stack_error"] = e.stackError
	}
	if e.errno != 0 {
		e.isNotify = true
	}
	if e.isNotify {
		res["notify"] = e.isNotify
	}
	return res
}

func (e *code) ToMapWithData(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"errno":   e.errno,
		"message": e.message,
		"details": data,
	}
}

func (e *code) ToMapWithError(err error) map[string]interface{} {
	var errStr string
	if err != nil {
		errStr = err.Error()

	}
	return map[string]interface{}{
		"errno":   e.errno,
		"message": e.message + ": " + errStr,
	}
}

func (e *code) WrapData(details interface{}) Errno {
	if e.isConst {
		en := &code{
			httpCode:   e.httpCode,
			errno:      e.errno,
			message:    e.message,
			stackError: e.stackError,
			details:    details,
			isStack:    e.isStack,
			isNotify:   e.isNotify,
		}
		return en

	}
	e.details = details
	return e
}

func (e *code) WrapError(err error) Errno {
	if e.isConst {
		en := &code{
			httpCode:   e.httpCode,
			errno:      e.errno,
			message:    e.message,
			stackError: e.stackError,
			details:    err.Error(),
			isStack:    e.isStack,
			isNotify:   e.isNotify,
		}
		return en
	}
	e.details = err.Error()
	return e
}
func (e *code) Notify() Errno {
	if e.isConst {
		en := &code{
			httpCode:   e.httpCode,
			errno:      e.errno,
			message:    e.message,
			stackError: e.stackError,
			details:    e.details,
			isStack:    e.isStack,
			isNotify:   true,
		}
		return en
	}
	e.isNotify = true
	return e
}
