/*
 * @Author: reel
 * @Date: 2023-05-16 20:44:40
 * @LastEditors: reel
 * @LastEditTime: 2023-09-12 06:04:50
 * @Description: 配置定时任务
 */
package cron

import (
	"fmt"
	"sync"
	"time"
)

var _ Cron = (*cron)(nil)

type task struct {
	job           func()
	name          string // 作业名称
	operationTime int    // 作业执行时间
	intervalTime  int    // 间隔时间
	failCount     int    // 失败次数, 连续 3 次失败, 停止此作业
	isRunning     bool   // 是否启用
}

type cron struct {
	run  bool
	name string
	job  map[string]*task
	lock *sync.RWMutex
}

type Cron interface {
	Start() error
	Stop() error
	Name() string
	Status() int8
	JobStop(name string)
	JobStart(name string)
	JobState() []interface{}
	AddJob(f func(), name string, interval int)
}

// TODO: 增加日志, 完善相关功能
func (c *cron) Start() error {
	go func() {
		defer func() {
			c.run = false
			err := recover()
			if err != nil {
				fmt.Println("定时执行中发生错误")
			}
		}()
		c.run = true
		for c.run {
			time.Sleep(1 * time.Second)
			c.lock.RLock()
			for _, v := range c.job {
				c.exec(v)
			}
			c.lock.RUnlock()

		}
	}()

	return nil
}

func (c *cron) Stop() (err error) {
	c.run = false
	return
}

func (c *cron) Name() (name string) {
	return c.name
}

func (c *cron) Status() int8 {
	if c.run {
		return 1
	}
	return 0
}

// 日志在程序中完成
func (c *cron) AddJob(f func(), name string, interval int) {
	f1 := func() {
		defer func() {
			err := recover()
			job := c.job[name]
			if err != nil {
				job.failCount += 1
				if job.failCount >= 3 {
					job.isRunning = false
				}
				fmt.Println(fmt.Sprintf("定时作业:[%s]发生错误:%v", name, err))
				return
			}
			job.failCount = 0
		}()
		f()
	}
	// 最小间隔 1 秒
	if interval == 0 {
		interval = 1
	}

	t := &task{
		job:          f1,
		name:         name,
		isRunning:    true,
		intervalTime: interval,
	}
	c.lock.Lock()
	c.job[name] = t
	c.lock.Unlock()
}

// 停用作业
func (c *cron) JobStop(name string) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	job := c.job[name]
	if job != nil {
		job.isRunning = false
	}
}

// 启用作业
func (c *cron) JobStart(name string) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	job := c.job[name]
	if job != nil {
		job.isRunning = true
	}
}

// 执行任务
func (c *cron) exec(t *task) {
	if !t.isRunning {
		return
	}
	if t.operationTime > t.intervalTime {
		t.operationTime = 0
		go t.job()
	}
	t.operationTime += 1
}

// 默认开启秒级别的定时任务
// 不记录日志,
func New() Cron {
	c := &cron{
		name: "Corn",
		lock: &sync.RWMutex{},
		job:  make(map[string]*task, 10),
	}
	return c
}

func (c *cron) JobState() []interface{} {
	list := make([]interface{}, 0, 10)
	c.lock.RLock()
	for _, v := range c.job {
		list = append(list, v.dump())
	}
	c.lock.RUnlock()
	return list
}

func (t *task) dump() (data map[string]interface{}) {
	return map[string]interface{}{
		"name":          t.name,
		"interval_time": t.intervalTime,
		"fail_count":    t.failCount,
		"is_running":    t.isRunning,
	}

}
