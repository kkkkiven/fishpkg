// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package db

import (
	"bytes"
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetBufferPool() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func PutBufferPool(b *bytes.Buffer) {
	b.Reset()
	bufferPool.Put(b)
}
