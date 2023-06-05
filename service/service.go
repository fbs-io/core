/*
 * @Author: reel
 * @Date: 2023-05-16 21:09:14
 * @LastEditors: reel
 * @LastEditTime: 2023-06-05 07:55:21
 * @Description: 请填写简介
 */
package service

import (
    "fmt"
    "time"

    "github.com/fbs-io/core/pkg/errorx"
)

var services = make([]Service, 0, 20)

type Service interface {
    Start() error
    Stop() error
    Name() string
    Status() int8
}

func Assertion(p interface{}) (srv Service, err error) {
    defer func() {
        e := recover()
        if e != nil {
            srv = nil
            err = errorx.New(fmt.Sprintf("%v", e))
        }
    }()
    srv = p.(Service)
    return
}

func NewServices() []Service {
    return services
}

func Append(service Service) {
    services = append(services, service)
}

func Start() error {
    for _, s := range services {
        err := s.Start()
        if err != nil {
            return errorx.Wrap(err, fmt.Sprintf("[%s]服务启动失败", s.Name()))
        }
    }
    return nil
}

func Stop() error {
    for _, s := range services {
        err := s.Stop()
        if err != nil {

            return errorx.Wrap(err, fmt.Sprintf("[%s]服务关闭失败", s.Name()))
        }
    }
    return nil
}

func Restart() {
    for _, srv := range services {
        srv.Stop()
    }
    time.Sleep(30 * time.Second)
    for _, srv := range services {
        srv.Start()
    }
}
