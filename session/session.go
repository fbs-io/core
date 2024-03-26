/*
 * @Author: reel
 * @Date: 2023-06-06 19:21:05
 * @LastEditors: reel
 * @LastEditTime: 2024-03-27 04:45:16
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
	"strings"
	"time"

	"github.com/fbs-io/core/pkg/env"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/cache"
	"github.com/fbs-io/core/store/dsn"
	"github.com/google/uuid"
)

var (
	ERROR_SESSION_LENGTH     = errorx.New("无法正确获取到session")
	ERROR_SESSION_NOT_LOGIN  = errorx.New("账户未登陆或登陆信息已生效，请重新登陆")
	ERROR_SESSION_ELSE_LOGIN = errorx.New("已在其他地方登陆, 请重新登陆")
)

type session struct {
	lifeTime   int // 秒
	cookieName string
	prefix     string
	singular   string // 校验是否单用户登陆
	store      cache.Store
}

type Session interface {
	sessionP()
	CookieName() string
	Singular() string // 是否单用户登陆, Y表示是, N或空表示否
	SetWithToken(sessionKey, sessionValue string)
	GetWithCookie(r *http.Request) (cookieValue, value string, err error)
	SetWithCookie(w http.ResponseWriter, cookieValue, internalValue string)
	GetWithToken(r *http.Request) (sessionKey, sessionValue string, err error)
	GetSessionWithCookie(r *http.Request, w http.ResponseWriter) (sessionValue string, err error)
	SetWithSid(w http.ResponseWriter, cookieValue, internalValue string)
	GetWithSid(r *http.Request) (sessionKey, sessionValue string, err error)

	// 服务端想客户端设置cookie, 使用请求头的X-CSRF-TOKEN字段
	SetWithCsrfToken(w http.ResponseWriter, cookieValue, internalValue string)

	// 客户端发送请求, 使用cookie传输
	GetWithCsrfToken(r *http.Request) (sessionKey, sessionValue string, err error)
}

func (s *session) sessionP() {}

var _ Session = (*session)(nil)

// 默认有效期30分钟
// 默认cookiename=sid, 单用户登陆
// 默认使用本地缓存, 底层为buntdb作为缓存支撑, 接口为store/cache.store
func New(funs ...optFunc) Session {
	opt := &option{
		lifeTime:   1800,
		cookieName: "sid",
		singular:   "Y",
	}
	for _, f := range funs {
		f(opt)
	}
	if opt.store == nil {
		opt.store = cache.New()
		opt.store.SetConfig(dsn.SetName("session"))
		opt.store.Start()

	}
	lifeTime := opt.lifeTime
	opt.lifeTime = opt.lifeTime * 48 * 30
	if env.Active().Value() != env.ENV_MODE_DEV {
		opt.prefix = fmt.Sprintf("%s::%d", opt.prefix, time.Now().UnixNano())
		opt.lifeTime = lifeTime
	}
	s := &session{
		lifeTime:   opt.lifeTime,
		cookieName: opt.cookieName,
		store:      opt.store,
		prefix:     opt.prefix,
		singular:   opt.singular,
	}

	return s
}

// 自动生成cookie的值, 36位长度
// 同时把cookie写入缓存中
// 如果没有设置缓存的值, 以cookie名称补充, 表示为未登陆用户
func (s *session) SetWithCookie(w http.ResponseWriter, sessionKey, sessionValue string) {
	if sessionKey == "" {
		sessionKey = GenSessionKey()
	}
	if sessionValue == "" {
		sessionValue = s.cookieName
	}
	s.setSession(sessionKey, sessionValue)

	cookie := &http.Cookie{
		Name:     s.cookieName,
		Value:    url.QueryEscape(sessionKey),
		MaxAge:   s.lifeTime,
		Path:     "/",
		Domain:   "",
		SameSite: 4,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

// 获取cookie
func (s *session) GetWithCookie(r *http.Request) (sessionKey, sessionValue string, err error) {
	cookie, err := r.Cookie(s.cookieName)
	if err != nil {
		return "", "", ERROR_SESSION_NOT_LOGIN
	}
	sessionKey, _ = url.QueryUnescape(cookie.Value)
	return s.getSession(sessionKey)
}

// 用于前端cookie存储的名字
func (s *session) CookieName() string {
	return s.cookieName
}

// 用于设置默认cookie的场景
// 如防止中间人攻击等情况
func (s *session) GetSessionWithCookie(r *http.Request, w http.ResponseWriter) (sessionValue string, err error) {
	cookieValue, sessionValue, err := s.GetWithCookie(r)
	s.SetWithCookie(w, cookieValue, sessionValue)
	return
}

// 设置token
func (s *session) SetWithToken(sessionKey, sessionValue string) {
	s.setSession(sessionKey, sessionValue)
}

// 获取token
func (s *session) GetWithToken(r *http.Request) (sessionKey, sessionValue string, err error) {
	token := r.Header.Get("Authorization")
	sessionKey, _ = url.QueryUnescape(token)
	return s.getSession(sessionKey)
}

func GenSessionKey() string {
	b := make([]byte, 36)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return base64.URLEncoding.EncodeToString([]byte(uuid.New().String()))
	}
	return base64.URLEncoding.EncodeToString(b)
}

// 自动生成session的值, 36位长度
// 同时把session写入缓存中请求头SID中
// 如果没有设置缓存的值, 以cookie名称补充, 表示为未登陆用户
func (s *session) SetWithSid(w http.ResponseWriter, sessionKey, sessionValue string) {
	if sessionKey == "" {
		sessionKey = GenSessionKey()
	}
	if sessionValue == "" {
		sessionValue = s.cookieName
	}
	s.setSession(sessionKey, sessionValue)
	w.Header().Set("SID", sessionKey)
}

// 通过sid获取session
func (s *session) GetWithSid(r *http.Request) (sessionKey, sessionValue string, err error) {
	cookie := r.Header.Get(s.cookieName)

	sessionKey, _ = url.QueryUnescape(cookie)
	return s.getSession(sessionKey)
}

// 自动生成session的值, 36位长度
// 同时把session写入缓存中请求头X-CSRF-TOKEN中
// 如果没有设置缓存的值, 以cookie名称补充, 表示为未登陆用户
func (s *session) SetWithCsrfToken(w http.ResponseWriter, sessionKey, sessionValue string) {
	if sessionKey == "" {
		sessionKey = GenSessionKey()
	}
	if sessionValue == "" {
		sessionValue = s.cookieName
	}
	s.setSession(sessionKey, sessionValue)
	w.Header().Set("X-CSRF-TOKEN", sessionKey)
}

// 通过Authorization获取session
func (s *session) GetWithCsrfToken(r *http.Request) (sessionKey, sessionValue string, err error) {
	token := r.Header.Get("Authorization")
	auth, _ := url.QueryUnescape(token)
	auths := strings.Split(auth, " ")
	sessionKey = auths[0]
	if len(auth) == 2 {
		sessionKey = auths[1]
	}
	return s.getSession(sessionKey)
}

func (s *session) GenStoreKey(sessionKey string) string {
	return fmt.Sprintf("%s::%s", s.prefix, sessionKey)
}

// 获取session是否单用户登陆
func (s *session) Singular() string {
	return s.singular
}

// 设置session
func (s *session) setSession(sessionKey, sessionValue string) {
	s.store.Set(s.GenStoreKey(sessionKey), sessionValue, cache.SetTTL(time.Duration(s.lifeTime)))
	if s.singular == "Y" {
		s.store.Set(sessionValue, sessionKey)
	}
}

// 获取session
func (s *session) getSession(sessionKey string) (sessionKey2, sessionValue string, err error) {

	if len(sessionKey) != 48 {
		return "", "", ERROR_SESSION_LENGTH
	}
	sessionValue = s.store.Get(s.GenStoreKey(sessionKey))
	if sessionValue == "" || sessionValue == s.cookieName {
		return "", "", ERROR_SESSION_NOT_LOGIN
	}
	if s.singular == "Y" {
		lastKey := s.store.Get(sessionValue)
		if lastKey == "" {
			s.store.Set(sessionValue, sessionKey, cache.SetTTL(1800))
		} else {
			if lastKey != sessionKey {
				return "", "", ERROR_SESSION_ELSE_LOGIN
			}
		}
	}
	return sessionKey, sessionValue, nil
}
