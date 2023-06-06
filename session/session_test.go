/*
 * @Author: reel
 * @Date: 2023-06-06 21:59:57
 * @LastEditors: reel
 * @LastEditTime: 2023-06-06 23:10:34
 * @Description: 请填写简介
 */
package session

import (
    "encoding/base64"
    "fmt"
    "testing"

    "github.com/fbs-io/core/pkg/env"
    "github.com/google/uuid"
)

func TestSession(t *testing.T) {
    env.Init()
    fmt.Println(genCookieValue())
    fmt.Println(base64.URLEncoding.EncodeToString([]byte(uuid.New().String())))
    s := New()
    // s.Get("123")
    s.Set(nil, "")
}
