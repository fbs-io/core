/*
 * @Author: reel
 * @Date: 2023-05-11 22:19:24
 * @LastEditors: reel
 * @LastEditTime: 2023-06-06 06:15:20
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
}

type code struct {
    httpCode   int    // http 状态码
    errno      int    // 业务状态码
    message    string // 业务状态描述
    stackError error  // 含有堆栈的状态描述
}

// 生成新的 code 暂时不需要调用错误堆栈错误, 需要时再调用
func New(httpCode, errno int, msg string) Errno {
    return &code{
        httpCode: httpCode,
        errno:    errno,
        message:  msg,
    }
}

func (e *code) en() {}

func (e *code) WithStack(err error) Errno {
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

func (e *code) ToMap() map[string]interface{} {
    return map[string]interface{}{
        "errno":   e.errno,
        "message": e.message,
        // "error":   e.stackError.Error(),
    }
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
