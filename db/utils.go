// Copyright (c) 2020. Homeland Interactive Technology Ltd. All rights reserved.

package db

import (
	"database/sql"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/kkkkiven/fishpkg/util"
)

// NewNullString 返回一个带有Null值的数据库字符串
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// NewNullInt32 返回一个带有Null值的数据库整形
func NewNullInt32(s int32, isNull bool) sql.NullInt32 {
	return sql.NullInt32{
		Int32: s,
		Valid: !isNull,
	}
}

// NewNullInt64 返回一个带有Null值的数据库整形
func NewNullInt64(s int64, isNull bool) sql.NullInt64 {
	return sql.NullInt64{
		Int64: s,
		Valid: !isNull,
	}
}

// Quote 对参数转码
func Quote(s string) string {
	return strings.Replace(strings.Replace(s, "'", "", -1), `\`, `\\`, -1)
}

// FullSQL 返回绑定完参数的完整的SQL语句
func FullSQL(str string, args ...interface{}) (string, error) {
	if !strings.Contains(str, "?") {
		return str, nil
	}
	sons := strings.Split(str, "?")

	var ret string
	var argIndex int
	var maxArgIndex = len(args)

	for _, son := range sons {
		ret += son

		if argIndex < maxArgIndex {
			switch v := args[argIndex].(type) {
			case int:
				ret += strconv.Itoa(v)
			case int8:
				ret += strconv.Itoa(int(v))
			case int16:
				ret += strconv.Itoa(int(v))
			case int32:
				ret += util.I64toa(int64(v))
			case int64:
				ret += util.I64toa(v)
			case uint:
				ret += util.UItoa(v)
			case uint8:
				ret += util.UItoa(uint(v))
			case uint16:
				ret += util.UItoa(uint(v))
			case uint32:
				ret += util.UI32toa(v)
			case uint64:
				ret += util.UI64toa(v)
			case float32:
				ret += util.F32toa(v)
			case float64:
				ret += util.F64toa(v)
			case *big.Int:
				ret += v.String()
			case bool:
				if v {
					ret += "true"
				} else {
					ret += "false"
				}
			case string:
				ret += "'" + Quote(v) + "'"
			case nil:
				ret += "NULL"
			default:
				return "", fmt.Errorf(
					"invalid sql argument type: %v => %v (sql: %s)",
					reflect.TypeOf(v).String(), v, str)
			}

			argIndex++
		}
	}

	return ret, nil
}
