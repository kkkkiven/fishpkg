package utils

import "unsafe"

// String2Bytes 高性能string转Byte
func String2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// String2Bytes 高性能Byte转String
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
