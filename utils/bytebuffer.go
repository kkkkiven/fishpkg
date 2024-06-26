package utils

import (
	"bytes"
	"encoding/binary"
	//	"errors"
	//	"fmt"
	//	"reflect"
	//	"unsafe"
)

// Put8bit 加入8字节到[]byte中
func Put8bit(buf []byte, n byte) []byte {
	return append(buf, n)
}

// Put16bit 加入16字节
func Put16bit(buf []byte, n uint16) []byte {
	var by [2]byte

	by[0] = byte((n >> 8) & 0xff)
	by[1] = byte(n & 0xff)

	return append(buf, by[:]...)
}

// Put32bit 加入32字节
func Put32bit(buf []byte, n uint32) []byte {
	var by [4]byte

	by[0] = byte((n >> 24) & 0xff)
	by[1] = byte((n >> 16) & 0xff)
	by[2] = byte((n >> 8) & 0xff)
	by[3] = byte(n & 0xff)

	return append(buf, by[:]...)
}

// Put64bit 加入64字节
func Put64bit(buf []byte, n uint64) []byte {
	var by [8]byte

	by[0] = byte((n >> 56) & 0xff)
	by[1] = byte((n >> 48) & 0xff)
	by[2] = byte((n >> 40) & 0xff)
	by[3] = byte((n >> 32) & 0xff)
	by[4] = byte((n >> 24) & 0xff)
	by[5] = byte((n >> 16) & 0xff)
	by[6] = byte((n >> 8) & 0xff)
	by[7] = byte(n & 0xff)

	return append(buf, by[:]...)
}

// Get8bit 获取8bit
func Get8bit(buf []byte, start int) byte {
	return buf[start]
}

// Get16bit 获取16bit
func Get16bit(buf []byte, start int) uint16 {
	var ret uint16

	ret = uint16(buf[start]) << 8
	ret |= uint16(buf[start+1])

	return ret
}

// Get32bit 获取32big
func Get32bit(buf []byte, start int) uint32 {
	var ret uint32

	ret = uint32(buf[start]) << 24
	ret |= uint32(buf[start+1]) << 16
	ret |= uint32(buf[start+2]) << 8
	ret |= uint32(buf[start+3])

	return ret
}

// Get64bit 获取64bit
func Get64bit(buf []byte, start int) uint64 {
	var ret uint64

	ret = uint64(buf[start]) << 56
	ret |= uint64(buf[start+1]) << 48
	ret |= uint64(buf[start+2]) << 40
	ret |= uint64(buf[start+3]) << 32
	ret |= uint64(buf[start+4]) << 24
	ret |= uint64(buf[start+5]) << 16
	ret |= uint64(buf[start+6]) << 8
	ret |= uint64(buf[start+7])

	return ret
}

// BytesToInt 字节数组转int
func BytesToInt(b []byte) int {
	buf := bytes.NewBuffer(b)
	var x int
	binary.Read(buf, binary.BigEndian, &x)

	return int(x)
}

// IntToBytes int转字节数组
func IntToBytes(n int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)

	return buf.Bytes()
}

// BytesToInt16 字节数组转int16
func BytesToInt16(b []byte) int16 {
	buf := bytes.NewBuffer(b)
	var x int16
	binary.Read(buf, binary.BigEndian, &x)

	return int16(x)
}

// Int16ToBytes int16转字节数组
func Int16ToBytes(n int16) []byte {
	x := int16(n)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, x)

	return buf.Bytes()
}

// BytesToInt32 字节数组转int32
func BytesToInt32(b []byte) int32 {
	buf := bytes.NewBuffer(b)
	var x int32
	binary.Read(buf, binary.BigEndian, &x)

	return int32(x)
}

//	Int32ToBytes int32转字节数组
func Int32ToBytes(n int32) []byte {
	x := int32(n)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, x)

	return buf.Bytes()
}

// 字节数组转int64
func BytesToInt64(b []byte) int64 {
	buf := bytes.NewBuffer(b)
	var x int64
	binary.Read(buf, binary.BigEndian, &x)

	return int64(x)
}

// Int64ToBytes int64转字节数组
func Int64ToBytes(n int64) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)

	return buf.Bytes()
}

// BytesToUInt64 字节数组转uint64
func BytesToUInt64(b []byte) uint64 {
	buf := bytes.NewBuffer(b)
	var x uint64
	binary.Read(buf, binary.BigEndian, &x)

	return uint64(x)
}

// UInt64ToBytes uint64转字节数组
func UInt64ToBytes(n uint64) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)

	return buf.Bytes()
}

// UInt32ToBytes uint32转字节数组
func UInt32ToBytes(n uint32) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)

	return buf.Bytes()
}

// BytesToBool 字节数组转bool
func BytesToBool(b []byte) bool {
	buf := bytes.NewBuffer(b)
	var x bool
	binary.Read(buf, binary.BigEndian, &x)
	return x
}

// BoolToBytes bool转字节数组
func BoolToBytes(x bool) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, x)

	return buf.Bytes()
}
