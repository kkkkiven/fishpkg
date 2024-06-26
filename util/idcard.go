// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package util

import (
	"strings"
)

var (
	// 加权因子
	idCardFactor = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

	// 校验码对应值
	idCardVerifyNumberList = []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	// 省份定义
	idCardAreaMap = map[string]string{
		"11": "北京",
		"12": "天津",
		"13": "河北",
		"14": "山西",
		"15": "内蒙古",
		"21": "辽宁",
		"22": "吉林",
		"23": "黑龙江",
		"31": "上海",
		"32": "江苏",
		"33": "浙江",
		"34": "安徽",
		"35": "福建",
		"36": "江西",
		"37": "山东",
		"41": "河南",
		"42": "湖北",
		"43": "湖南",
		"44": "广东",
		"45": "广西",
		"46": "海南",
		"50": "重庆",
		"51": "四川",
		"52": "贵州",
		"53": "云南",
		"54": "西藏",
		"61": "陕西",
		"62": "甘肃",
		"63": "青海",
		"64": "宁夏",
		"65": "新疆",
		"71": "台湾",
		"81": "香港",
		"82": "澳门",
		"91": "国外",
	}
)

// IsIDCard 检查是否是身份证
// param: allowLen15 允许 15 位旧身份证号
func IsIDCard(idCard string, allowLen15 ...bool) bool {
	if len(idCard) != 15 && len(idCard) != 18 {
		return false
	}

	if _, ok := idCardAreaMap[idCard[:2]]; !ok {
		return false
	}

	// 如果是15位身份证
	if len(idCard) == 15 && len(allowLen15) == 1 && allowLen15[0] {
		idCard = IDCardUpgrade(idCard)
	}
	if !checkIDCardFormat(idCard) {
		return false
	}
	return strings.EqualFold(idCardVerifyNumber(idCard[:17]), idCard[17:18])
}

// IDCardUpgrade 15 位身份证升级为 18 位
func IDCardUpgrade(idcard string) string {
	if len(idcard) != 15 {
		return idcard
	}
	// 如果身份证顺序码是996 997 998 999，这些是为百岁以上老人的特殊编码
	switch idcard[12:15] {
	case "996", "997", "998", "999":
		idcard = idcard[:6] + "18" + idcard[6:15]
	default:
		idcard = idcard[:6] + "19" + idcard[6:15]
	}
	return idcard + idCardVerifyNumber(idcard)
}

// idCardVerifyNumber 计算身份证校验码，根据国家标准GB 11643-1999
func idCardVerifyNumber(idcard string) string {
	if len(idcard) != 17 {
		return ""
	}

	checkSum := 0
	for i, l := 0, len(idcard); i < l; i++ {
		checkSum += Atoi(idcard[i:i+1]) * idCardFactor[i]
	}
	checkSum %= 11
	return idCardVerifyNumberList[checkSum]
}

// checkIDCardFormat 检查身份证号码的格式
// 针对近期出现身份证号码使用210000000000000000注册的用户
// 验证地区码(3-6位)以及生日(6-14)位
func checkIDCardFormat(idcard string) bool {
	if len(idcard) == 18 {
		birthYear := Atoi(idcard[6:10])
		birthMonth := Atoi(idcard[10:12])
		birthDay := Atoi(idcard[12:14])
		if birthYear < 1900 || (birthMonth == 0 || birthMonth > 12) || (birthDay == 0 || birthDay > 31) {
			return false
		}
		return true
	}
	return false
}
