// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package util

import (
	"bytes"
	"math/rand"
	"sync"
	"time"
)

const defaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

var defaultStringRand = NewStringRand(defaultAlphabet, 8)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandString 生成指定长度的字符串
func RandString(length int) string {
	return defaultStringRand.GenerateWithLength(length)
}

// NewStringRand 新建一个随机字符串生成器
// nolint
func NewStringRand(alphabet string, length int) *StringRand {
	return &StringRand{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
		bytePool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		alphabet:      []byte(alphabet),
		alphabetCount: len(alphabet),
		genLength:     length,
	}
}

type StringRand struct {
	r             *rand.Rand
	bytePool      *sync.Pool
	alphabet      []byte
	alphabetCount int
	genLength     int
}

// SetGenerateLength 设置字符串生成长度
func (sr *StringRand) SetGenerateLength(length int) {
	sr.genLength = length
}

func (sr *StringRand) getBuffer() *bytes.Buffer {
	return sr.bytePool.Get().(*bytes.Buffer)
}
func (sr *StringRand) putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	sr.bytePool.Put(buf)
}

// Generate 生成字符串
func (sr *StringRand) Generate() string {
	return sr.GenerateWithLength(sr.genLength)
}

// GenerateWithLength 生成指定长度的字符串
func (sr *StringRand) GenerateWithLength(length int) string {
	ret := sr.getBuffer()
	defer sr.putBuffer(ret)
	for i := 0; i < length; i++ {
		ret.WriteByte(sr.alphabet[sr.r.Intn(sr.alphabetCount)])
	}
	return ret.String()
}

// Rand 随机数字 0 <= n < max
// nolint
func Rand(max int) int {
	return rand.Intn(max)
}

// RandInt 随机一个数字 min <= n < max
func RandInt(min, max int) int {
	if max == min {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return min + Rand(max-min)
}

// Rand64 随机数字 0 <= n < max
// nolint
func Rand64(max int64) int64 {
	return rand.Int63n(max)
}

// RandInt64 随机一个数字 min <= n < max
func RandInt64(min, max int64) int64 {
	if max == min {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return min + Rand64(max-min)
}
