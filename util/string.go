// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package util

import "strings"

// Substr 截取字符串
// 例: abc你好1234
// Substr(str, 0) == abc你好1234
// Substr(str, 2) == c你好1234
// Substr(str, -2) == 34
// Substr(str, 2, 3) == c你好
// Substr(str, 0, -2) == 34
// Substr(str, 2, -1) == b
// Substr(str, -3, 2) == 23
// Substr(str, -3, -2) == 好1
func Substr(str string, start int, length ...int) string {
	rs := []rune(str)
	lth := len(rs)
	end := 0

	if start > lth {
		return ""
	}

	if len(length) == 1 {
		end = length[0]
	}

	// 从后数的某个位置向后截取
	if start < 0 {
		if -start >= lth {
			start = 0
		} else {
			start = lth + start
		}
	}

	if end == 0 {
		end = lth
	} else if end > 0 {
		end += start
		if end > lth {
			end = lth
		}
	} else { // 从指定位置向前截取
		if start == 0 {
			start = lth
		}
		start, end = start+end, start
	}
	if start < 0 {
		start = 0
	}

	return string(rs[start:end])
}

// Trim 清除左右两边空格
func Trim(str string) string {
	return strings.Trim(str, " \r\n\t")
}

// StrInSlice 检查字符串是否存在与某切片中
func StrInSlice(s []string, item string) bool {
	for _, s := range s {
		if s == item {
			return true
		}
	}
	return false
}
