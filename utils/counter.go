package utils

import (
	"sync/atomic"
)

// Counter 原子计数器
type Counter struct {
	v int64
}

// Add 计数加
func (c *Counter) Add(i int64) {
	atomic.AddInt64(&c.v, i)
}

// Get 取计数
func (c *Counter) Get() int64 {
	return c.v
}

// Reset 重置计数器
func (c *Counter) Reset() {
	c.Add(c.v * -1)
}