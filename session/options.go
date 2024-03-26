/*
 * @Author: reel
 * @Date: 2023-06-06 22:40:43
 * @LastEditors: reel
 * @LastEditTime: 2024-03-26 07:01:09
 * @Description: session初始化相关设置
 */
package session

import "github.com/fbs-io/core/store/cache"

type option struct {
	lifeTime   int
	prefix     string
	cookieName string
	singular   string
	store      cache.Store
}

type optFunc func(*option)

// 设置过期时间
// 默认1800秒, 30分钟
func Lifetime(lifeTime int) optFunc {
	return func(opt *option) {
		opt.lifeTime = lifeTime
	}
}

// 设置客户端存储的cookie名称
// 默认名称sid
func Cookiename(cookieName string) optFunc {
	return func(opt *option) {
		opt.cookieName = cookieName
	}
}

// 设置存储,
// 默认使用本地缓存, 底层为buntdb作为缓存支撑, 接口为store/cache.store
func Store(cache cache.Store) optFunc {
	return func(opt *option) {
		opt.store = cache
	}
}

func Prefix(prefix string) optFunc {
	return func(opt *option) {
		opt.prefix = prefix
	}
}

// 设置是否单用户登陆
func Singular(singular string) optFunc {
	return func(opt *option) {
		opt.singular = singular
	}
}
