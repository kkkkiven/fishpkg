// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package util

import (
	"sync/atomic"
)

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount 查表计算一个 2 进制 byte 中 1 的个数
func PopCount(b byte) byte {
	return pc[b]
}

// Counter 原子计数器
type Counter struct {
	v int64
}

// Add 计数加
func (c *Counter) Add(i int64) {
	atomic.AddInt64(&c.v, i)
}

// Set 取计数
func (c *Counter) Set(i int64) {
	atomic.StoreInt64(&c.v, i)
}

// Get 取计数
func (c *Counter) Get() int64 {
	return c.v
}

// Reset 重置计数器
func (c *Counter) Reset() {
	c.Add(c.v * -1)
}
