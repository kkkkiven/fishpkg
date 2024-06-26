package utils

import (
	"reflect"
	"regexp"
	"strings"
	"net/http"
	"strconv"
	"time"
	"math"
)

// IsIDCard 检查是否是身份证
func IsIDCard(str string, args ...bool) bool {
	if m, _ := regexp.MatchString("^(?:[\\dxX]{15}|[\\dxX]{18})$", str); !m {
		return false
	}

	areaMap := map[int]string{
		11 : "北京",
		12 : "天津",
		13 : "河北",
		14 : "山西",
		15 : "内蒙古",
		21 : "辽宁",
		22 : "吉林",
		23 : "黑龙江",
		31 : "上海",
		32 : "江苏",
		33 : "浙江",
		34 : "安徽",
		35 : "福建",
		36 : "江西",
		37 : "山东",
		41 : "河南",
		42 : "湖北",
		43 : "湖南",
		44 : "广东",
		45 : "广西",
		46 : "海南",
		50 : "重庆",
		51 : "四川",
		52 : "贵州",
		53 : "云南",
		54 : "西藏",
		61 : "陕西",
		62 : "甘肃",
		63 : "青海",
		64 : "宁夏",
		65 : "新疆",
		71 : "台湾",
		81 : "香港",
		82 : "澳门",
		91 : "国外",
	}

	allowLen15 := true
	if len(args) == 1 {
		allowLen15 = args[0]
	}

	province := Atoi(Substr(str, 0, 2))
	if _, ok := areaMap[province]; !ok {
		return false
	}

	//如果是15位身份证
	if len(str) == 15 && allowLen15 {
		// 如果身份证顺序码是996 997 998 999，这些是为百岁以上老人的特殊编码
		code := Substr(str, 12, 3)
		if code == "996" || code == "997" || code == "998" || code == "999" {
			str = Substr(str, 0, 6) + "18" + Substr(str, 6, 9)
		} else {
			str = Substr(str, 0, 6) + "19" + Substr(str, 6, 9)
		}
		str += idcardVerifyNumber(str)
	}
	if !checkIdcardFormat(str) {
		return false
	}
	return idcardVerifyNumber(Substr(str, 0, 17)) == strings.ToUpper(Substr(str, 17, 1))
}

// checkIdcardFormat 检查身份证号码的格式
// 针对近期出现身份证号码使用210000000000000000注册的用户
// 验证地区码(3-6位)以及生日(6-14)位
func checkIdcardFormat(str string) bool {
	if len(str) == 18 {
		birthYear := Atoi(Substr(str, 6, 4))
		birthMonth := Atoi(Substr(str, 10, 2))
		birthDay := Atoi(Substr(str, 10, 2))
		if birthYear < 1900 || (birthMonth == 0 || birthMonth > 12) || (birthDay == 0 || birthDay > 31) {
			return false
		}
		return true
	}
	return false
}

// idcardVerifyNumber 计算身份证校验码，根据国家标准GB 11643-1999
func idcardVerifyNumber(str string) string {
	if len(str) != 17 {
		return ""
	}

	//加权因子
	factor := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	//校验码对应值
	verifyNumberList := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	checkSum := 0
	for i, l := 0, len(str); i < l; i++ {
		checkSum += Atoi(Substr(str, i, 1)) * factor[i]
	}
	checkSum = checkSum % 11
	return verifyNumberList[checkSum]
}

// GetInfoByIDCard 通过身份证号获取公民信息
// 注意: 调用此接口前请自行通过 IsIDCard 验证身份证的合法性
// 返回值 sex 0为女, 1为男
func GetInfoByIDCard(idcard string) (region string, birthday time.Time, sex int) {
	b := []byte(idcard)
	region = string(b[:6])
	switch len(idcard) {
	case 15:
		code := string(b[12:15])
		if code == "996" || code == "997" || code == "998" || code == "999" {
			birthday, _ = time.ParseInLocation("20060102", "18" + string(b[6:12]), time.Local)
		} else {
			birthday, _ = time.ParseInLocation("20060102", "19" + string(b[6:12]), time.Local)
		}
		sex = Atoi(string(b[14:])) % 2
	case 18:
		birthday, _ = time.ParseInLocation("20060102", string(b[6:14]), time.Local)
		sex = Atoi(string(b[16:17])) % 2
	}
	return
}

// 是否为空
func IsEmpty(val interface{}) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			b = true
		}
	}()
	v := reflect.ValueOf(val)

	switch v.Kind() {
	case reflect.Bool:
		b = (val.(bool) == false)
	case reflect.String:
		b = (val.(string) == "")
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		b = (v.Len() == 0)
	default:
		b = (v.Interface() == reflect.ValueOf(0).Interface() || v.Interface() == reflect.ValueOf(0.0).Interface())
	}

	return b
}

// 判断是否正确的ip地址
func IsIp(ip string) bool {
	ips := strings.Split(ip, ".")
	if len(ips) != 4 {
		return false
	}
	for _, v := range ips {
		i := Atoi(v, -1)
		if i < 0 || i > 255 {
			return false
		}
	}

	return true
}

//检测是否为手机设备
func IsMobileDevice(r *http.Request) bool {
	agent := r.UserAgent()
	if strings.Contains(agent,"iPhone") || strings.Contains(agent,"iPad") || strings.Contains(agent,"iOS") || strings.Contains(agent,"Android") {
		return true;
	}
	return false;
}

// 检查字符串是否为纯数字组成
func IsNumStr(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

// 检查字符串是否为纯数字组成
func IsLongNumStr(s string) bool {
	if m, _ := regexp.MatchString("^\\d+$", s); m {
		return true
	}
	return false
}

// IsMobileNum 检查是否是手机号
func IsMobileNum(str string) bool {
	if m, _ := regexp.MatchString("^1[3-9][0-9]{9}$", str); m {
		return true
	}
	return false
}


// IsVersionStr 检查字符串是否是一个版本格式
func IsVersionStr(str string) bool {
	return IsVersion(str, "{2}")
}

// IsVersion 检查字符串是否是一个版本格式
func IsVersion(str string, length ...string) bool {
	vL := "{1,3}"
	if len(length) == 1 {
		vL = length[0]
	}
	if m, _ := regexp.MatchString("^\\d+(?:\\.\\d+)" + vL + "$", str); m {
		return true
	}
	return false
}

// IsDeviceCode 检查字符串是否是一个机器码
func IsDeviceCode(str string) bool {
	if m, _ := regexp.MatchString("^[a-zA-Z0-9]{40}$", str); m {
		return true
	}
	return false
}

// RegexpOK 检查字符串被正则匹配
func RegexpOK(pattern, str string) bool {
	if m, _ := regexp.MatchString(pattern, str); m {
		return true
	}
	return false
}

// DiffDay 比较两个日志间隔的天数
// 第1个日期小于第2个日期时返回正数
// 第1个日期大于第2个日期时返回负数
func DiffDay(t1, t2 string, fm ...string) (int, error) {
	format := "2006-01-02"
	if len(fm) == 1 {
		format = fm[0]
	}
	tm1, err := time.ParseInLocation(format, t1, time.Local)
	if err != nil {
		return 0, err
	}
	tm2, err := time.ParseInLocation(format, t2, time.Local)
	if err != nil {
		return 0, err
	}
	hour := tm2.Sub(tm1).Hours()
	if -24 >= hour || hour >= 24 {
		return int(math.Ceil(hour / 24)), nil
	}
	return 0, nil
}