// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package db

// NewSBValues 获取一个 SBValues 对象
func NewSBValues() SBValues {
	return SBValues{}
}

// SBValues 值对象
type SBValues map[string]interface{}

// Add 向值对象中加入值
func (v SBValues) Add(key string, val interface{}) {
	v[key] = val
}

// Del 删除值对象中的某个值
func (v SBValues) Del(key string) {
	delete(v, key)
}

// IsExist 判断指定键是否存在
func (v SBValues) IsExist(key string) bool {
	if _, exist := v[key]; exist {
		return true
	}
	return false
}

// Get 获取键的整形值
func (v SBValues) Get(key string) interface{} {
	if val, exist := v[key]; exist {
		return val
	}
	return nil
}

// GetString 获取键的字符串值
func (v SBValues) GetString(key string) string {
	if val, exist := v[key]; exist {
		if trueVal, ok := val.(string); ok {
			return trueVal
		}
	}
	return ""
}

// GetInt 获取键的整形值
func (v SBValues) GetInt(key string) int {
	if val, exist := v[key]; exist {
		if trueVal, ok := val.(int); ok {
			return trueVal
		}
	}
	return 0
}

// GetUint 获取键的无符号整形值
func (v SBValues) GetUint(key string) uint {
	if val, exist := v[key]; exist {
		if trueVal, ok := val.(uint); ok {
			return trueVal
		}
	}
	return 0
}

// GetInt64 获取键的64位整形值
func (v SBValues) GetInt64(key string) int64 {
	if val, exist := v[key]; exist {
		if trueVal, ok := val.(int64); ok {
			return trueVal
		}
	}
	return 0
}
