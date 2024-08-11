/*
 * @Author: reel
 * @Date: 2024-07-07 15:29:15
 * @LastEditors: reel
 * @LastEditTime: 2024-07-07 15:29:17
 * @Description: 请填写简介
 */
/*
 * @Author: reel
 * @Date: 2023-09-19 05:19:41
 * @LastEditors: reel
 * @LastEditTime: 2023-09-27 21:55:46
 * @Description: 生成编码参数
 */

package sequence

import "github.com/fbs-io/core/pkg/consts"

var (
	mapping = map[string]string{
		consts.DATE_FORMAT_Y:            "2006",
		consts.DATE_FORMAT_YM:           "200601",
		consts.DATE_FORMAT_YMD:          "20060102",
		consts.DATE_FORMAT_YMDH:         "2006010215",
		consts.DATE_FORMAT_YMDHM:        "200601021504",
		consts.DATE_FORMAT_YMDHMS:       "20060102150405",
		consts.DATE_FORMAT_SHORT_Y:      "06",
		consts.DATE_FORMAT_SHORT_YM:     "0601",
		consts.DATE_FORMAT_SHORT_YMD:    "060102",
		consts.DATE_FORMAT_SHORT_YMDH:   "06010215",
		consts.DATE_FORMAT_SHORT_YMDHM:  "0601021504",
		consts.DATE_FORMAT_SHORT_YMDHMS: "060102150405",
	}
)

type opts struct {
	split       string // 分隔符
	prefix      string // 编码前缀
	dateFormat  string // 日期格式
	sequenceLen int    // 流水号长度

}

type OptFunc func(*opts)

// 设置分隔符
func SetSplit(split string) func(*opts) {
	return func(opt *opts) {
		opt.split = split
	}
}

// 设置编码前缀
func SetPrefix(prefix string) func(*opts) {
	return func(opt *opts) {
		opt.prefix = prefix
	}
}

// 设置日期格式
//
// 支持 YY YYYY YYMM YYYYMM YYMMDD YYYYMMDD YYMMDDhh YYYYMMDDhh YYMMDDhhmm YYYYMMDDhhmm YYMMDDhhmmss YYYYMMDDhhmmss
//
// 日期格式不正确, 将统一处理为不适用日期作为可变序列前缀
func SetDateFormat(dateFormat string) func(*opts) {
	return func(opt *opts) {
		opt.dateFormat = mapping[dateFormat]
	}
}

// 设置流水号长度
func SetSequenceLen(sequenceLen int) func(*opts) {
	return func(opt *opts) {
		opt.sequenceLen = sequenceLen
	}
}

type warpCode struct {
	varPrefix  string
	customCode string
	suffix     string
}

type WarpCodeFunc func(w *warpCode)

// 设置可变前缀
func SetVarPrefix(varPrefix string) WarpCodeFunc {
	return func(w *warpCode) {
		w.varPrefix = varPrefix
	}
}

// 设置自定义code
func SetCustomCode(customCode string) WarpCodeFunc {
	return func(w *warpCode) {
		w.customCode = customCode
	}
}

// 设置后缀
//
// 后缀不作为查询
func Setsuffix(suffix string) WarpCodeFunc {
	return func(w *warpCode) {
		w.suffix = suffix
	}
}
