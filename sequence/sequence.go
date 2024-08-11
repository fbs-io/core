/*
 * @Author: reel
 * @Date: 2023-09-19 05:06:15
 * @LastEditors: reel
 * @LastEditTime: 2024-07-07 16:00:29
 * @Description: 业务编码生成器
 */
package sequence

import (
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/consts"
	"github.com/fbs-io/core/store/cache"
)

// TODO: 定义规则
// type sequenceRule struct {
// }

type sequenceMode struct {
	name        string      // 生成器名称
	split       string      //分隔符
	prefix      string      //前缀
	dateFormat  string      //日期格式
	sequenceLen int         //流水号长度
	cache       cache.Store //缓存,用于记录业务号
	lock        chan struct{}
	// serial      map[string]chan int32 //通过管道自动获取业务流水号, 业务量大的情况下可以考虑
}

type Sequence interface {
	sequenceP()
	Code(fs ...WarpCodeFunc) string
}

func New(c core.Core, name string, funcs ...OptFunc) Sequence {
	opt := &opts{
		split:       "",
		prefix:      "",
		dateFormat:  "20060102",
		sequenceLen: 4,
	}

	for _, f := range funcs {
		f(opt)
	}

	s := &sequenceMode{
		name:        name,
		split:       opt.split,
		prefix:      opt.prefix,
		dateFormat:  opt.dateFormat,
		sequenceLen: opt.sequenceLen,
		cache:       c.Cache(),
		lock:        make(chan struct{}, 1),
	}
	s.setLock()

	return s
}

func (s *sequenceMode) sequenceP() {}

func (s *sequenceMode) setLock() {
	s.lock <- struct{}{}
}

func (s *sequenceMode) getLock() {
	<-s.lock
}

// 生成code
func (s *sequenceMode) Code(fs ...WarpCodeFunc) string {
	s.getLock()
	defer s.setLock()
	w := &warpCode{}
	for _, f := range fs {
		f(w)
	}
	codeList := make([]string, 0, 10)

	// 前缀
	if s.prefix != "" {
		codeList = append(codeList, s.prefix)
	}
	// 自定义前缀
	if w.varPrefix != "" {
		codeList = append(codeList, w.varPrefix)
	}
	// 自定义code
	if w.customCode != "" {
		codeList = append(codeList, w.customCode)
	}
	// 日期格式
	if s.dateFormat != "" {
		codeList = append(codeList, time.Now().Format(s.dateFormat))
	}
	expired := 0 // 一天的秒数
	switch s.dateFormat {
	case consts.DATE_FORMAT_Y, consts.DATE_FORMAT_SHORT_Y:
		expired = 86400 * 366
	case consts.DATE_FORMAT_YM, consts.DATE_FORMAT_SHORT_YM:
		expired = 86400 * 31
	case consts.DATE_FORMAT_YMD, consts.DATE_FORMAT_SHORT_YMD:
		expired = 86400
	case consts.DATE_FORMAT_YMDH, consts.DATE_FORMAT_SHORT_YMDH:
		expired = 3600
	case consts.DATE_FORMAT_YMDHM, consts.DATE_FORMAT_SHORT_YMDHM:
		expired = 60
	case consts.DATE_FORMAT_YMDHMS, consts.DATE_FORMAT_SHORT_YMDHMS:
		expired = 1
	}

	// 获取序列号及设置序列号
	var sequenceValue int
	var index = strings.Join(codeList, s.split)
	value := s.cache.Get(index)
	if value == "" {
		sequenceValue = 1
	} else {
		sequenceValue, _ = strconv.Atoi(value)
		sequenceValue += 1
	}
	codeList = append(codeList, fmt.Sprintf("%0*d", s.sequenceLen, sequenceValue))
	if expired != 0 {
		s.cache.Set(index, fmt.Sprintf("%d", sequenceValue), cache.SetTTL(time.Duration(expired)))
	} else {
		s.cache.Set(index, fmt.Sprintf("%d", sequenceValue))
	}
	if w.suffix != "" {
		codeList = append(codeList, w.suffix)
	}
	return strings.Join(codeList, s.split)
}
