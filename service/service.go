/*
 * @Author: reel
 * @Date: 2023-05-16 21:09:14
 * @LastEditors: reel
 * @LastEditTime: 2024-03-30 09:40:07
 * @Description: 请填写简介
 */
package service

import (
	"fmt"
	"time"

	"github.com/fbs-io/core/logx"
	"github.com/fbs-io/core/pkg/errorx"
)

var (
	services = make([]Service, 0, 20)
)

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
	serviceName := ""
	defer func() {
		err := recover()
		if err != nil {
			logx.Sys.Info(fmt.Sprintf("[%s]服务启动过程中发生位置错误: %v", serviceName, err))
		}
	}()
	for _, s := range services {
		time.Sleep(1 * time.Second)
		serviceName = s.Name()
		logx.Sys.Info(fmt.Sprintf("[%s]服务开始启动: ", serviceName))
		err := s.Start()
		if err != nil {
			logx.Sys.Error(fmt.Sprintf("[%s]服务启动发生错误: ", serviceName), logx.E(err))
			return errorx.Wrap(err, fmt.Sprintf("[%s]服务启动失败", serviceName))
		}
		logx.Sys.Info(fmt.Sprintf("[%s]服务完成启动: ", serviceName))
	}
	return nil
}

func Stop() error {
	serviceName := ""
	defer func() {
		err := recover()
		if err != nil {
			logx.Sys.Info(fmt.Sprintf("[%s]服务关闭过程中发生位置错误: %v", serviceName, err))
		}
	}()
	for _, s := range services {
		time.Sleep(1 * time.Second)
		err := s.Stop()
		serviceName = s.Name()
		logx.Sys.Info(fmt.Sprintf("[%s]服务开始关闭: ", serviceName))
		if err != nil {
			logx.Sys.Info(fmt.Sprintf("[%s]服务关闭发生错误: ", serviceName), logx.E(err))
			return errorx.Wrap(err, fmt.Sprintf("[%s]服务关闭失败", serviceName))
		}
		logx.Sys.Info(fmt.Sprintf("[%s]服务完成关闭: ", serviceName))
	}
	return nil
}

func Restart() {
	logx.Sys.Info("服务开始重启")
	for _, srv := range services {
		srv.Stop()
	}
	time.Sleep(10 * time.Second)
	for _, srv := range services {
		srv.Start()
	}
	logx.Sys.Info("服务重启完成")
}

func Status() (srvStatus []map[string]interface{}) {
	srvStatus = make([]map[string]interface{}, 0, 20)
	for _, srv := range services {
		srvStatus = append(srvStatus, map[string]interface{}{
			"service": srv.Name(),
			"node1":   srv.Status(),
		})
	}
	return
}
