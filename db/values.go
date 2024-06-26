// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package db

import (
	"sync"
)

const defaultValueCapacity = 10

var (
	poolValues = sync.Pool{
		New: func() interface{} {
			return NewValues()
		},
	}
)

// NewValues 创建 Values 实体
func NewValues() *Values {
	return &Values{
		Keys: make([]string, 0, defaultValueCapacity),
		Vals: make([]interface{}, 0, defaultValueCapacity),
	}
}

// GetValues 从对象池中获取一个 *Values 对象
// 注意: SQLBuilder 调用 Exec, Query, QueryOne, QueryAllRow, QueryRow 这些方法后会自动释放传入的 *Values 对象,
// 		如无必要请使用 NewValues() 替代此方法会更加安全
func GetValues() *Values {
	return poolValues.Get().(*Values)
}

// PutValues 将 *Values 对象回收放回对象池
func PutValues(vs *Values) {
	vs.Reset()
	poolValues.Put(vs)
}

type Values struct {
	Keys []string
	Vals []interface{}
}

// Reset 重置对象
func (vs *Values) Reset() {
	vs.Keys = vs.Keys[:0]
	vs.Vals = vs.Vals[:0]
}

// Add 添加元素
func (vs *Values) Add(k string, v interface{}) *Values {
	vs.Keys = append(vs.Keys, k)
	vs.Vals = append(vs.Vals, v)
	return vs
}

// Adds 添加多个元素
func (vs *Values) Adds(k string, val ...interface{}) *Values {
	if len(val)&1 == 0 {
		return vs
	}
	vs.Keys = append(vs.Keys, k)
	for i, v := range val {
		if i&1 == 0 {
			vs.Vals = append(vs.Vals, v)
		} else {
			vs.Keys = append(vs.Keys, v.(string))
		}
	}
	return vs
}

// AddMap 以字典方式添加值, 不推荐使用
func (vs *Values) AddMap(m map[string]interface{}) *Values {
	for k, v := range m {
		vs.Add(k, v)
	}
	return vs
}

// AddSBValues 以 SBValues 方式添加值, 不推荐使用
func (vs *Values) AddSBValues(m SBValues) *Values {
	for k, v := range m {
		vs.Add(k, v)
	}
	return vs
}

// Set 设置元素, 跟 Add 不同的是, 如果 key 存在则替换其值, 而不是新增
func (vs *Values) Set(k string, v interface{}) *Values {
	i := vs.findKey(k)
	if i == -1 {
		vs.Add(k, v)
	} else {
		vs.Vals[i] = v
	}
	return vs
}

func (vs *Values) findKey(s string) int {
	for i, k := range vs.Keys {
		if k == s {
			return i
		}
	}
	return -1
}

// Del 删除元素
func (vs *Values) Del(k string) *Values {
	_ = vs.GetAndDel(k)
	return vs
}

// Len 获取当前数据数量
func (vs *Values) Len() int {
	return len(vs.Keys)
}

// GetAndDel 删除指定元素并返回其值
func (vs *Values) GetAndDel(k string) interface{} {
	i := vs.findKey(k)
	if i == -1 {
		return nil
	}
	v := vs.Vals[i]
	vs.remove(i)
	return v
}

// remove 删除元素
func (vs *Values) remove(i int) {
	if i >= len(vs.Keys) || i < 0 {
		return
	}
	vs.Keys[i] = vs.Keys[len(vs.Keys)-1]
	vs.Keys = vs.Keys[:len(vs.Keys)-1]
	vs.Vals[i] = vs.Vals[len(vs.Vals)-1]
	vs.Vals = vs.Vals[:len(vs.Vals)-1]
}

// Range 遍历数据, 如果处理函数返回 false 则会中断遍历
func (vs *Values) Range(f func(k string, v interface{}) bool) {
	for i, k := range vs.Keys {
		if !f(k, vs.Vals[i]) {
			break
		}
	}
}

// IsExists 判断键是否存在
func (vs *Values) IsExists(k string) bool {
	return vs.findKey(k) != -1
}

// Get 根据 key 获取值
func (vs *Values) Get(k string) interface{} {
	i := vs.findKey(k)
	if i == -1 {
		return nil
	}
	return vs.Vals[i]
}

func (vs *Values) GetString(k string) string {
	v := vs.Get(k)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (vs *Values) GetInt64(k string) int64 {
	v := vs.Get(k)
	if v != nil {
		if i, ok := v.(int64); ok {
			return i
		}
	}
	return 0
}

func (vs *Values) GetUint64(k string) uint64 {
	v := vs.Get(k)
	if v != nil {
		if i, ok := v.(uint64); ok {
			return i
		}
	}
	return 0
}

func (vs *Values) GetInt(k string) int {
	v := vs.Get(k)
	if v != nil {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}

func (vs *Values) GetUint(k string) uint {
	v := vs.Get(k)
	if v != nil {
		if i, ok := v.(uint); ok {
			return i
		}
	}
	return 0
}

func (vs *Values) GetInt8(k string) int8 {
	v := vs.Get(k)
	if v != nil {
		if i, ok := v.(int8); ok {
			return i
		}
	}
	return 0
}

func (vs *Values) GetUint8(k string) uint8 {
	v := vs.Get(k)
	if v != nil {
		if i, ok := v.(uint8); ok {
			return i
		}
	}
	return 0
}
