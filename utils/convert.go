package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GetTypeName 获取参数类型
func GetTypeName(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// Btoi 布尔值转整形
func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Atoi 转换成整型
func Atoi(s string, d ...int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}

	return i
}

// AtoUi 转换成无符号整型
func AtoUi(s string) uint {
	return uint(Atoi64(s))
}

// Atoi64 转换成整型int64
func Atoi64(s string, d ...int64) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}

	return i
}

// AtoUi64 转换成整型float64
func AtoUi64(s string, d ...uint64) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}

	return i
}

// Atof 转换成float32整型
func Atof(s string, d ...float32) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}

	return float32(f)
}

// Atof64 转换成整型float64
func Atof64(s string, d ...float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}

	return f
}

// UitoA 32位无符号整形转字符串
func UitoA(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Ui16toA 16位无符号整形转字符串
func Ui16toA(i uint16) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Ui32toA 32位无符号整形转字符串
func Ui32toA(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Ui64toA 64位无符号整形转字符串
func Ui64toA(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// Itoa 整型转字符串
func Itoa(i int) string {
	return strconv.Itoa(i)
}

// I16toA 16位整型转字符串
func I16toA(i int16) string {
	return strconv.FormatInt(int64(i), 10)
}

// I32toA 32位整型转字符串
func I32toA(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// I64toA 64位整形转字符串
func I64toA(i int64) string {
	return strconv.FormatInt(i, 10)
}

// F32toA 32位浮点数转字符串
func F32toA(f float32) string {
	return F64toA(float64(f))
}

// F64toA 64位浮点数转字符串
func F64toA(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// DateFormat 日期格式化
func DateFormat(format string, t time.Time) string {
	patterns := []string{
		// 年
		"Y", "2006", // 4 位数字完整表示的年份
		"y", "06", // 2 位数字表示的年份

		// 月
		"m", "01", // 数字表示的月份，有前导零
		"n", "1", // 数字表示的月份，没有前导零
		"M", "Jan", // 三个字母缩写表示的月份
		"F", "January", // 月份，完整的文本格式，例如 January 或者 March

		// 日
		"d", "02", // 月份中的第几天，有前导零的 2 位数字
		"j", "2", // 月份中的第几天，没有前导零

		"D", "Mon", // 星期几，文本表示，3 个字母
		"l", "Monday", // 星期几，完整的文本格式;L的小写字母

		// 时间
		"g", "3", // 小时，12 小时格式，没有前导零
		"G", "15", // 小时，24 小时格式，没有前导零
		"h", "03", // 小时，12 小时格式，有前导零
		"H", "15", // 小时，24 小时格式，有前导零

		"a", "pm", // 小写的上午和下午值
		"A", "PM", // 小写的上午和下午值

		"i", "04", // 有前导零的分钟数
		"s", "05", // 秒数，有前导零
	}
	replacer := strings.NewReplacer(patterns...)
	format = replacer.Replace(format)
	return t.Format(format)
}

// StrToLocalTime 字符串转本地时间
func StrToLocalTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	zoneName, offset := time.Now().Zone()

	zoneValue := offset / 3600 * 100
	if zoneValue > 0 {
		value += fmt.Sprintf(" +%04d", zoneValue)
	} else {
		value += fmt.Sprintf(" -%04d", zoneValue)
	}

	if zoneName != "" {
		value += " " + zoneName
	}
	return StrToTime(value)
}

// StrToTime 字符串转时间
func StrToTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	layouts := []string{
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05 -0700 MST",
		"2006/01/02 15:04:05 -0700",
		"2006/01/02 15:04:05",
		"2006-01-02 -0700 MST",
		"2006-01-02 -0700",
		"2006-01-02",
		"2006/01/02 -0700 MST",
		"2006/01/02 -0700",
		"2006/01/02",
		"2006-01-02 15:04:05 -0700 -0700",
		"2006/01/02 15:04:05 -0700 -0700",
		"2006-01-02 -0700 -0700",
		"2006/01/02 -0700 -0700",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, err
}

// StructToMap 结构转map
func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Name
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}

// 带xml标签的Struct结构转map，用xml标注做key
func XMLStructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Tag.Get("xml")
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}

// Ip2long IP转长整型
func Ip2long(ipstr string) (ip uint32) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return
	}
	ips := reg.FindStringSubmatch(ipstr)
	if ips == nil {
		return
	}

	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	return
}

// Long2ip 长整型转IP
func Long2ip(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip>>24, ip<<8>>24, ip<<16>>24, ip<<24>>24)
}

// SignStrBuild
func MapToRawStr(data map[string]string) string {

	str := ""
	if len(data) > 0 {
		keys := make([]string, 0)
		for k, _ := range data {
			if k != "sign" && k != "encode" && k != "v" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		//
		for _, k := range keys {
			if data[k] != "" {
				if str == "" {
					str = fmt.Sprintf("%s=%s", k, url.QueryEscape(strings.Replace(data[k], " ", "+", -1)))
				} else {
					str = fmt.Sprintf("%s&%s=%s", str, k, url.QueryEscape(strings.Replace(data[k], " ", "+", -1)))
				}
			}
		}
	}

	return str
}

// GbkToUtf8 GBK转UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbk UTF-8转GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
