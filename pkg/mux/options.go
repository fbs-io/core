/*
 * @Author: reel
 * @Date: 2023-05-12 06:40:04
 * @LastEditors: reel
 * @LastEditTime: 2023-05-23 07:36:45
 * @Description: http服务配置选项
 */
package mux

import (
    "time"
)

type opts struct {
    addr          string
    mode          string
    name          string
    maxHeaderSize int
    maxReadTime   time.Duration
    maxWriteTIme  time.Duration
    timeout       time.Duration
}

type OptFunc func(*opts)

// 设置监听地址
func SetHost(addr string) OptFunc {
    return func(o *opts) {
        o.addr = addr
    }
}

// 设置服务名称
func SetName(name string) OptFunc {
    return func(o *opts) {
        o.name = name
    }
}

// 启动模式
func SetMode(mode string) OptFunc {
    return func(o *opts) {
        o.mode = mode
    }
}

// 设置最大访问时间
func SetMaxReadTime(maxReadTime time.Duration) OptFunc {
    return func(o *opts) {
        o.maxReadTime = maxReadTime
    }
}

// 设置最大写入时间
func SetMaxWriteTIme(maxWriteTIme time.Duration) OptFunc {
    return func(o *opts) {
        o.maxWriteTIme = maxWriteTIme
    }
}

// 设置请求头最大字节
func SetMaxHeaderSize(maxHeaderSize int) OptFunc {
    return func(o *opts) {
        o.maxHeaderSize = maxHeaderSize
    }
}

// 设置关机最大停留时间
func SetTimeout(timeout int) OptFunc {
    return func(o *opts) {
        o.timeout = time.Second * time.Duration(timeout)
    }
}
