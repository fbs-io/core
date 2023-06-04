/*
 * @Author: reel
 * @Date: 2023-05-11 22:28:37
 * @LastEditors: reel
 * @LastEditTime: 2023-05-16 21:42:36
 * @Description: 请填写简介
 */
package filex

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// 判断文件路径是否存在
func IsExists(path string) (os.FileInfo, error) {
	f, err := os.Stat(path)
	return f, err
}

// 判断文件路径是否存在, 如果不存在则创建
func CreatDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return err
		}
	}
	return nil
}

// 字符串/文本按行读取, 根据 mac, Windows, Linux 等不同平台进行分割
func SplitStrToLine(str string) (strList []string) {
	//windows 换行判断,
	strList = strings.Split(str, "\r\n")
	if len(strList) == 1 {
		// mac 换行判断
		strList = strings.Split(str, "\r")
	}
	if len(strList) == 1 {
		//linux 换行判断
		strList = strings.Split(str, "\n")
	}
	linLen := len(strList)
	if strList[linLen-1] == "" {
		strList = strList[:linLen-1]
	}
	return
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", val, suffix)
}

func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}
