// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package util

import (
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"sync"
)

var (
	_regPatternChinaMobile = regexp.MustCompile(`^1[3-9][0-9]{9}$`)
	_regPatternNumber      = regexp.MustCompile("^\\d+$")

	_patterns = &sync.Map{}
)

// RegexpOK 检查字符串是否能被正则匹配
func RegexpOK(pattern, str string) bool {
	v, ok := _patterns.Load(pattern)
	var p *regexp.Regexp
	if ok {
		p, ok = v.(*regexp.Regexp)
		if ok {
			return p.MatchString(str)
		}
	}
	
	var err error
	p, err = regexp.Compile(pattern)
	if err != nil {
		return false
	}
	_patterns.Store(pattern, p)
	return p.MatchString(str)
}

// IsChinaMobile 检查是否为中国大陆手机号
func IsChinaMobile(str string) bool {
	return _regPatternChinaMobile.MatchString(str)
}

// IsURL 检查字符串是否是一个 url 地址
func IsURL(str string) bool {
	_, err := url.Parse(str)
	return err == nil
}

// IsEmail 检查字符串是否是一个 Email 地址
func IsEmail(str string) bool {
	_, err := mail.ParseAddress(str)
	return err == nil
}

// IsNumStr 检查字符串是否为纯数字组成
func IsNumStr(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

// IsLongNumStr 检查字符串是否为纯数字组成
func IsLongNumStr(s string) bool {
	return _regPatternNumber.MatchString(s)
}
