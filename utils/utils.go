package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Rand 随机数字 0 <= n < max
func Rand(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

// RandInt 随机一个数字 min <= n < max
func RandInt(min int, max int) int {
	if max == min {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	if max < min {
		min, max = max, min
	}
	return min + rand.Intn(max-min)
}

// Rand64 随机数字 0 <= n < max
func Rand64(max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max)
}

// RandInt64 随机一个数字 min <= n < max
func RandInt64(min int64, max int64) int64 {
	if max == min {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	if max < min {
		min, max = max, min
	}
	return min + rand.Int63n(max-min)
}

// Max 取两个数较大的一个
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MaxInt64 取两个数较大的一个
func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// Min 取两个数较小的一个
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MinInt64 取两个数较小的一个
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// RandNumStr 生成多位随机数字符串
func RandNumStr(l int) string {
	ret := ""
	for i := 0; i < l; i++ {
		ret += strconv.Itoa(Rand(10))
	}
	return ret
}

// Reverse 字符串反转
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// RandArray 随机一个数组值
func RandArray(arr []string) string {
	return arr[Rand(len(arr))]
}

// 随机一个字符串
func RandString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Range 生成一个数组 不包括n
func Range(m, n int) (b []int) {
	if m >= n {
		return b
	}

	for i := m; i < n; i++ {
		b = append(b, i)
	}

	return b
}

// Keys 返回map中的所有字符串key
func KeysByString(m map[string]interface{}) []string {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}

	return keys
}

// Values 返回map中的所有value
func Values(m map[string]interface{}) []interface{} {
	var values []interface{}
	for _, v := range m {
		values = append(values, v)
	}

	return values
}

// InArray 是否在列表中
func InArray(v string, s []string) int {
	for i, val := range s {
		if val == v {
			return i
		}
	}
	return -1
}

// InSlice 是否在Slice中
/*func InSlice(v interface{}, s interface{}) int {
	if ss, ok := s.([]interface{}); ok {
		for i, val := range ss {
			if val == v {
				return i
			}
		}
	}
	return -1
}*/
func InSlice(v interface{}, s interface{}) int {
	switch reflect.TypeOf(s).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(s)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(v, s.Index(i).Interface()) == true {
				return i
			}
		}
	}
	return -1
}

// Trim 清除左右两边空格
func Trim(str string) string {
	return strings.Trim(str, " \r\n\t")
}

// CheckFatal 检查致命错误
func CheckFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// CheckPanic 检查恐慌
func CheckPanic(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// CheckFormEmpty 检查Form表单的某个值是否为空
// 为空返回真
func CheckFormEmpty(form url.Values, key string) bool {
	if _, ok := form[key]; ok {
		if Trim(form[key][0]) != "" {
			return false
		}
	}
	return true
}

// 检查Form表单的哪个字段值为空
// 发现空值时返回该字段的名称, 否则返回空白字符
func CheckFormEmptyByKeys(form url.Values, keys string) string {
	arr := strings.Split(keys, ",")
	for _, key := range arr {
		k := Trim(key)
		if CheckFormEmpty(form, k) {
			return k
		}
	}
	return ""
}

// PrintStruct 用来打印一个结构的字段与值对应表
func PrintStruct(structPtr interface{}) {
	s := reflect.ValueOf(structPtr).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s (%s) = %v\n",
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

// Len 获取带有中文等非ASCII字符的字符串长度
func Len(str string) int {
	rs := []rune(str)
	return len(rs)
}

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

	//从后数的某个位置向后截取
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
	} else { //从指定位置向前截取
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

// SplitIDStr 把以半角逗号分隔的ID字符串分隔提取到切片中
func SplitIDStr(str string) []int {
	ret := make([]int, 0)
	ids := strings.Split(str, ",")
	for _, idStr := range ids {
		idStr = Trim(idStr)
		if IsNumStr(idStr) {
			ret = append(ret, Atoi(idStr))
		}
	}
	return ret
}

// IntArrToString 将[]int转换成以指定分隔符分隔的字符串
func IntArrToString(a []int, separator ...string) string {
	sep := ","
	if len(separator) == 1 {
		sep = separator[0]
	}
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", sep, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), sep), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), sep), "[]")
}

// StrPad 向字符串中添加指定字符串到指定长度
func StrPad(str interface{}, length int, pad interface{}, padLeft bool) string {
	returnStr := ""
	padStr := ""
	if "int" == GetTypeName(str) {
		returnStr = Itoa(str.(int))
	}

	if "int" == GetTypeName(pad) {
		padStr = Itoa(pad.(int))
	}

	padLen := length - Len(returnStr)
	if padLen > 0 {
		padString := ""
		for i := 0; i < padLen; i++ {
			padString += padStr
		}
		if padLeft {
			returnStr = padString + returnStr
		} else {
			returnStr += padString
		}
	}
	return returnStr
}

// 返回sync.Map类型的长度
func LenSyncMap(m *sync.Map) int {
	var length int
	m.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

// Ternary 模拟三元操作符
func Ternary(b bool, trueVal, falseVal interface{}) interface{} {
	if b {
		return trueVal
	}
	return falseVal
}

// GetMacAddr 获取本地mac地址
func GetLocalMacAddr() ([]byte, error) {

	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil || len(interfaces) == 0 {
		return nil, err
	}
	return interfaces[0].HardwareAddr, nil

}
