// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"sync"
)

var (
	md5Pool = sync.Pool{
		New: func() interface{} {
			return md5.New()
		},
	}

	sha256Pool = sync.Pool{
		New: func() interface{} {
			return sha256.New()
		},
	}
)

// StringMD5 sum md5 of string
func StringMD5(text string) string {
	h := md5Pool.Get().(hash.Hash)
	_, err := h.Write(StringToBytes(text))
	if err != nil {
		return ""
	}
	ret := hex.EncodeToString(h.Sum(nil))
	h.Reset()
	md5Pool.Put(h)
	return ret
}

// StringSha256 sum sha256 of string
func StringSha256(text string) string {
	h := sha256Pool.Get().(hash.Hash)
	h.Reset()
	_, err := h.Write(StringToBytes(text))
	if err != nil {
		return ""
	}
	ret := hex.EncodeToString(h.Sum(nil))
	h.Reset()
	sha256Pool.Put(h)
	return ret

}
