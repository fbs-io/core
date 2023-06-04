/*
 * @Author: reel
 * @Date: 2023-05-16 21:19:26
 * @LastEditors: reel
 * @LastEditTime: 2023-06-04 22:35:53
 * @Description: base62加密解密封装
 */
package encrypt

import (
    "math"
    "strconv"
    "strings"
)

var (
    EDOC = map[string]int{
        "0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
        "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16, "h": 17, "i": 18,
        "j": 19, "k": 20, "l": 21, "m": 22, "n": 23, "o": 24, "p": 25, "q": 26, "r": 27,
        "s": 28, "t": 29, "u": 30, "v": 31, "w": 32, "x": 33, "y": 34, "z": 35, "A": 36,
        "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42, "H": 43, "I": 44, "J": 45,
        "K": 46, "L": 47, "M": 48, "N": 49, "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54,
        "T": 55, "U": 56, "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61}
)

const (
    CODE62     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    CODE_LENTH = 62
)

func Decode(str string) int {
    str = strings.TrimSpace(str)
    var result int = 0
    for index, _ := range str {
        result += EDOC[str[len(str)-index-1:len(str)-index]] * int(math.Pow(CODE_LENTH, float64(index)))
    }
    return result
}

func Encode(num int) string {
    var result string
    if num == 0 {
        return "0"
    }
    for num > 0 {
        remain := num % CODE_LENTH
        result = string(CODE62[remain]) + result
        num = num / CODE_LENTH
    }
    return result
}

func DecodeWBMid(str string) string {
    base62MidList := splitWithNumRight(str, 4)
    var res string
    for x, m := range base62MidList {
        result := Decode(m)
        a := strconv.Itoa(result)

        for i := 0; i <= (7 - len(a)); i++ {
            if x < len(base62MidList)-1 {
                if len(a) < 7 {
                    a = "0" + a
                }
            }
        }
        res = a + res
    }

    return res
}

// 由 mid 生成 base62 编码
func EncodeWBMid(mid string) string {

    result := ""
    sMidList := splitWithNumRight(mid, 7)
    for _, ns := range sMidList {
        num, _ := strconv.Atoi(ns)

        result = Encode(num) + result
    }
    return string(result)
}

// 按照字符串长度切割字符
// 从右向左, 按照 step 长度切割
func splitWithNumRight(str string, step int) (numlist []string) {
    numlist = make([]string, 0, 10)
    if str == "" {
        return
    }
    l := len(str)
    count := l / step
    remain := l % step
    for i := 0; i < count; i++ {
        numlist = append(numlist, str[l-((i+1)*step):l-(i*step)])
    }
    if remain > 0 {
        numlist = append(numlist, str[:remain])
    }
    return

}
