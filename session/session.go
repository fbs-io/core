/*
 * @Author: reel
 * @Date: 2023-06-06 19:21:05
 * @LastEditors: reel
 * @LastEditTime: 2023-06-07 07:32:59
 * @Description: session 模块
 */
package session

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"

    "github.com/fbs-io/core/store/cache"
    "github.com/fbs-io/core/store/dsn"
    "github.com/google/uuid"
)

// type Values sync.Map

type session struct {
    lifeTime   int // 秒
    cookieName string
    store      cache.Store
}

type Session interface {
    sessionP()
    CookieName() string
    Get(r *http.Request) (cookieValue, value string, err error)
    Set(w http.ResponseWriter, cookieValue, internalValue string)
    Check(r *http.Request, w http.ResponseWriter) (sessionValue string, err error)
}

func (s *session) sessionP() {}

var _ Session = (*session)(nil)

// 默认有效期30分钟
// 默认cookiename=sid
// 默认使用本地缓存, 底层为buntdb作为缓存支撑, 接口为store/cache.store
func New(funs ...optFunc) Session {
    opt := &option{
        lifeTime:   1800,
        cookieName: "sid",
    }
    for _, f := range funs {
        f(opt)
    }
    if opt.store == nil {
        opt.store = cache.New()
        opt.store.SetConfig(dsn.SetName("session"))
        opt.store.Start()

    }

    s := &session{
        lifeTime:   opt.lifeTime,
        cookieName: opt.cookieName,
        store:      opt.store,
    }

    return s
}

// 自动生成cookie的值, 36位长度
// 同时把cookie写入缓存中
// 如果没有设置缓存的值, 以cookie名称补充, 表示为未登陆用户
func (s *session) Set(w http.ResponseWriter, cookieValue, internalValue string) {
    if cookieValue == "" {
        cookieValue = genCookieValue()
    }
    if internalValue == "" {
        internalValue = s.cookieName
    }
    s.store.Set(cookieValue, internalValue, cache.SetTTL(time.Duration(s.lifeTime)))
    cookie := &http.Cookie{
        Name:     s.cookieName,
        Value:    url.QueryEscape(cookieValue),
        MaxAge:   s.lifeTime,
        Path:     "/",
        Domain:   "",
        SameSite: 0,
        Secure:   false,
        HttpOnly: true,
    }
    http.SetCookie(w, cookie)
    return
}

// 获取cookie
func (s *session) Get(r *http.Request) (cookieValue, value string, err error) {
    cookie, err := r.Cookie(s.cookieName)
    if err != nil {
        return "", "", err
    }
    val, _ := url.QueryUnescape(cookie.Value)
    sessionValue := s.store.Get(val)
    if len(val) != 48 {
        return "", "", fmt.Errorf("无法正确获取到session, session长度:%d", len(sessionValue))
    }
    return val, sessionValue, nil
}

// 用于前端cookie存储的名字
func (s *session) CookieName() string {
    return s.cookieName
}

func (s *session) Check(r *http.Request, w http.ResponseWriter) (sessionValue string, err error) {
    cookieValue, sessionValue, err := s.Get(r)
    s.Set(w, cookieValue, sessionValue)
    return
}

func genCookieValue() string {
    b := make([]byte, 36)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return base64.URLEncoding.EncodeToString([]byte(uuid.New().String()))
    }
    return base64.URLEncoding.EncodeToString(b)
}
