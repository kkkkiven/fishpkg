package core

import (
	"bytes"
	"compress/zlib"
	"sync"
)

var (
	zlibP = NewZlibProcessor()
)

type zlibProcessor struct {
	pool sync.Pool
}

func NewZlibProcessor() *zlibProcessor {
	return &zlibProcessor{
		pool: sync.Pool{
			New: func() interface{} {
				return new(zlib.Writer)
			},
		},
	}
}

func (zp *zlibProcessor) Compress(src []byte) []byte {
	buf := &bytes.Buffer{}
	zw := zp.pool.Get().(*zlib.Writer)
	defer func() {
		zp.pool.Put(zw)
	}()
	zw.Reset(buf)
	_, _ = zw.Write(src)
	_ = zw.Close()
	return buf.Bytes()
}
