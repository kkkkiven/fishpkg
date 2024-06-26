// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package util

import (
	"bytes"
	"errors"
	"strconv"
	"sync"
)

var (
	errJSONFormat = errors.New("json format error")

	poolFlag = &sync.Pool{
		New: func() interface{} {
			return &flag{}
		},
	}

	poolDictKV = &sync.Pool{
		New: func() interface{} {
			return &dictKV{}
		},
	}

	poolDict = &sync.Pool{
		New: func() interface{} {
			return NewStringDict()
		},
	}
)

const (
	flagQuoteBegin = 1 << iota
	flagColonBegin
	flagCommaBegin
	flagEscapeBegin
	flagInValue
)

func getFlag() *flag {
	return poolFlag.Get().(*flag)
}

func putFlag(f *flag) {
	f.quote = 0
	f.flag = 0
	poolFlag.Put(f)
}

type flag struct {
	quote byte
	flag  int8
}

func (jf *flag) isQuoteBegin() bool {
	return jf.flag&flagQuoteBegin == flagQuoteBegin
}

func (jf *flag) quoteBegin() {
	jf.flag |= flagQuoteBegin
}

func (jf *flag) quoteEnd() {
	jf.flag &= ^flagQuoteBegin
}

func (jf *flag) isColonBegin() bool {
	return jf.flag&flagColonBegin == flagColonBegin
}

func (jf *flag) colonBegin() {
	jf.flag |= flagColonBegin
}

func (jf *flag) isCommaBegin() bool {
	return jf.flag&flagCommaBegin == flagCommaBegin
}

func (jf *flag) commaBegin() {
	jf.flag |= flagCommaBegin
}

func (jf *flag) isEscapeBegin() bool {
	return jf.flag&flagEscapeBegin == flagEscapeBegin
}

func (jf *flag) escapeBegin() {
	jf.flag |= flagEscapeBegin
}

func (jf *flag) escapeEnd() {
	jf.flag &= ^flagEscapeBegin
}

func (jf *flag) isInValue() bool {
	return jf.flag&flagInValue == flagInValue
}

func (jf *flag) inValue() {
	jf.flag |= flagInValue
}

func (jf *flag) outValue() {
	jf.flag &= ^flagInValue
}

func getKVWithValue(k, v string) *dictKV {
	kv := getKV()
	kv.Key = append(kv.Key[:0], k...)
	kv.Val = append(kv.Val[:0], v...)
	return kv
}

func getKV() *dictKV {
	return poolDictKV.Get().(*dictKV)
}

func putKV(kv *dictKV) {
	kv.Key = kv.Key[:0]
	kv.Val = kv.Val[:0]
	poolDictKV.Put(kv)
}

type dictKV struct {
	Key []byte
	Val []byte
}

// GetStringDict 从对象池中获取一个 *StringDict
func GetStringDict() *StringDict {
	return poolDict.Get().(*StringDict)
}

// NewStringDict 创建字符串字典
func NewStringDict() *StringDict {
	return &StringDict{
		KVs: make([]*dictKV, 0, 8),
		buf: new(bytes.Buffer),
	}
}

// StringDict 一个可复用的 map[string]stirng 实现
// 注意! 该字典同样是无序, 且非线程安全
type StringDict struct {
	KVs []*dictKV
	buf *bytes.Buffer
}

// Index 查找元素所在索引
func (dict *StringDict) Index(k string) int {
	for i, kv := range dict.KVs {
		if bytes.Equal(kv.Key, StringToBytes(k)) {
			return i
		}
	}
	return -1
}

// Add 添加元素
func (dict *StringDict) Add(k, v string) *StringDict {
	dict.KVs = append(dict.KVs, getKVWithValue(k, v))
	return dict
}

// Set 设置元素, 跟 Add 不同的是, 如果 key 存在则替换其值, 而不是新增
func (dict *StringDict) Set(k, v string) *StringDict {
	i := dict.Index(k)
	if i == -1 {
		dict.Add(k, v)
	} else {
		dict.KVs[i].Val = append(dict.KVs[i].Val[:0], v...)
	}
	return dict
}

// Has 获取值, 如果不存在 bool 返回值将是 false
func (dict *StringDict) Has(k string) (string, bool) {
	i := dict.Index(k)
	if i == -1 {
		return "", false
	}
	return string(dict.KVs[i].Val), true
}

// Get 获取值
func (dict *StringDict) Get(k string) string {
	i := dict.Index(k)
	if i == -1 {
		return ""
	}
	return string(dict.KVs[i].Val)
}

// Del 删除元素
func (dict *StringDict) Del(k string) *StringDict {
	i := dict.Index(k)
	if i != -1 {
		dict.remove(i)
	}
	return dict
}

// Len 获取当前数据数量
func (dict *StringDict) Len() int {
	return len(dict.KVs)
}

// GetAndDel 删除指定元素并返回其值
func (dict *StringDict) GetAndDel(k string) string {
	i := dict.Index(k)
	if i == -1 {
		return ""
	}
	v := dict.KVs[i].Val
	dict.remove(i)
	return string(v)
}

// remove 删除元素
func (dict *StringDict) remove(i int) {
	if i >= len(dict.KVs) || i < 0 {
		return
	}
	putKV(dict.KVs[i])
	dict.KVs[i] = dict.KVs[len(dict.KVs)-1]
	dict.KVs = dict.KVs[:len(dict.KVs)-1]
}

// Range 遍历数据, 如果处理函数返回 false 则会中断遍历
func (dict *StringDict) Range(f func(k, v string) bool) {
	for _, kv := range dict.KVs {
		if !f(string(kv.Key), string(kv.Val)) {
			break
		}
	}
}

// IsExists 判断键是否存在
func (dict *StringDict) IsExists(k string) bool {
	return dict.Index(k) != -1
}

// Reset 重置
func (dict *StringDict) Reset() {
	for i := range dict.KVs {
		putKV(dict.KVs[i])
	}
	dict.KVs = dict.KVs[:0]
	dict.buf.Reset()
}

// Release 收回字典
func (dict *StringDict) Release() {
	dict.Reset()
	poolDict.Put(dict)
}

// MarshalJSON 自定义 JSON 编码
func (dict *StringDict) MarshalJSON() (b []byte, err error) {
	if len(dict.KVs) == 0 {
		return []byte("{}"), nil
	}
	dict.buf.WriteString("{")
	for i, kv := range dict.KVs {
		if i > 0 {
			dict.buf.WriteString(",")
		}
		dict.buf.WriteString(strconv.Quote(BytesToString(kv.Key)))
		dict.buf.WriteString(":")
		dict.buf.WriteString(strconv.Quote(BytesToString(kv.Val)))
	}
	dict.buf.WriteString("}")
	return dict.buf.Bytes(), nil
}

// UnmarshalJSON 自定义 JSON 解码
func (dict *StringDict) UnmarshalJSON(data []byte) error {
	var pos int
	kv := poolDictKV.Get().(*dictKV)
	f := getFlag()
	defer putFlag(f)
	for i, c := range data {
		if f.isEscapeBegin() {
			if !f.isQuoteBegin() {
				return errJSONFormat
			}
			if f.isInValue() {
				kv.Val = append(kv.Val, data[i])
			} else {
				kv.Key = append(kv.Key, data[i])
			}
			pos = i + 1
			f.escapeEnd()
			continue
		}
		switch c {
		case '\'', '"':
			if !f.isQuoteBegin() {
				f.quoteBegin()
				f.quote = c
				pos = i + 1
			} else if f.quote == c {
				if f.isInValue() {
					kv.Val = append(kv.Val, data[pos:i]...)
					dict.KVs = append(dict.KVs, kv)
					kv = &dictKV{}
					f.outValue()
				} else {
					kv.Key = append(kv.Key, data[pos:i]...)
				}
				f.quoteEnd()
			}
		case ':':
			if f.isQuoteBegin() {
				continue
			}
			if !f.isColonBegin() {
				if f.isInValue() {
					return errJSONFormat
				}
				f.inValue()
			}
		case '\\':
			f.escapeBegin()
			if i > pos {
				if f.isInValue() {
					kv.Val = append(kv.Val, data[pos:i]...)
				} else {
					kv.Key = append(kv.Key, data[pos:i]...)
				}
			}
			pos = i + 1
		}
	}
	return nil
}

func (dict *StringDict) String() string {
	if dict.Len() == 0 {
		return "{}"
	}
	b, _ := dict.MarshalJSON()
	return BytesToString(b)
}
