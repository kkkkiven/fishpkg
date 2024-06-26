// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package util

import (
	"encoding/binary"
	"reflect"
	"strconv"
	"unsafe"
)

// StringToBytes 字符串转字节切片
// 需要注意的是该方法极不安全，使用过程中应足够谨慎，防止各类访问越界的问题
// nolint
func StringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

// BytesToString 字节切片转字符串
// nolint
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Btoi boolean to int
func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Btoi8 boolean to int8
func Btoi8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

// Atoi8 string to int8
func Atoi8(s string, d ...int8) int8 {
	i, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	return int8(i)
}

// Atoi16 string to int8
func Atoi16(s string, d ...int16) int16 {
	i, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	return int16(i)
}

// Atoi string to int
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

// Atoi32 string to int32
func Atoi32(s string, d ...int32) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}

	return int32(i)
}

// Atoi64 string to int64
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

// Atoui8 string to uint8
func Atoui8(s string, d ...uint8) uint8 {
	i, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	return uint8(i)
}

// Atoui16 string to uint16
func Atoui16(s string, d ...uint16) uint16 {
	i, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	return uint16(i)
}

// Atoui string to uint
func Atoui(s string, d ...uint) uint {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	return uint(i)
}

// Atoui32 string to uint32
func Atoui32(s string, d ...uint32) uint32 {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	return uint32(i)
}

// Atoui64 string to uint64
func Atoui64(s string, d ...uint64) uint64 {
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

// Atof string to float32
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

// Atof64 string to float64
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

// I8toa int8 转字符串
func I8toa(i int8) string {
	return strconv.FormatInt(int64(i), 10)
}

// I16toa int16 转字符串
func I16toa(i int16) string {
	return strconv.FormatInt(int64(i), 10)
}

// Itoa int 转字符串
func Itoa(i int) string {
	return strconv.Itoa(i)
}

// I32toa int32 转字符串
func I32toa(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// I64toa int64 转字符串
func I64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}

// UI8toa uint8 转字符串
func UI8toa(i uint8) string {
	return strconv.FormatUint(uint64(i), 10)
}

// UI16toa uint16 转字符串
func UI16toa(i uint16) string {
	return strconv.FormatUint(uint64(i), 10)
}

// UItoa uint 转字符串
func UItoa(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

// UI32toa uint32 转字符串
func UI32toa(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

// UI64toa uint64 转字符串
func UI64toa(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// F32toa float32 转字符串
func F32toa(f float32) string {
	return F64toa(float64(f))
}

// F64toa float64 转字符串
func F64toa(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// Int32Tobytes int32 转 []byte
func Int32Tobytes(i int32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	return bytes
}

// BytesToInt32 []byte 转 int64
func BytesToInt32(bytes []byte) int32 {
	return int32(binary.LittleEndian.Uint32(bytes))
}

// IntTobytes int 转 []byte, 这将把 int 按 4 个字节处理
func IntTobytes(i int) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	return bytes
}

// BytesToInt []byte 转 int, 这将把 int 按 4 个字节处理
func BytesToInt(bytes []byte) int {
	return int(int32(binary.LittleEndian.Uint32(bytes)))
}
