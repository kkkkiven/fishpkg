package utils

import (
	"strconv"
	"strings"
	"net/http"
	"fmt"
)

// buildLuaResponse构造Lua响应结果
func BuildLuaResponse(m interface{}) string {
	var ret string
	if mm, ok := m.(map[interface{}]interface{}); ok {
		for k, i := range mm {
			vStr := ""
			switch v := i.(type) {
			case int:
				vStr = strconv.Itoa(v)
			case int8:
				vStr = strconv.Itoa(int(v))
			case int16:
				vStr = strconv.Itoa(int(v))
			case int32:
				vStr = I64toA(int64(v))
			case int64:
				vStr = I64toA(v)
			case uint:
				vStr = UitoA(v)
			case uint8:
				vStr = UitoA(uint(v))
			case uint16:
				vStr = UitoA(uint(v))
			case uint32:
				vStr = Ui32toA(v)
			case uint64:
				vStr = Ui64toA(v)
			case bool:
				if v {
					vStr = "true"
				}
			case []interface{}:
				tmp := make(map[interface{}]interface{})
				for kk, vv := range v {
					tmp[kk] = vv
				}
				vStr = BuildLuaResponse(tmp)
			case map[interface{}]interface{}:
				vStr = BuildLuaResponse(v)
			case map[string]string:
				for ks, vs := range v {
					vStr += `,` + ks + `=` + strings.Replace(vs, `"`, `\"`, -1)
				}
				vStr = "{" + Substr(vStr, 1) + "}"
			case string:
				vStr = `"` + strings.Replace(v, `"`, `\"`, -1) + `"`
			case map[string]interface{}:
				for ks, vs := range v {
					vStr += `,` + ks + `=` + `"`+strings.Replace(fmt.Sprint(vs), `"`, `\"`, -1)+`"`
				}
				vStr = "{" + Substr(vStr, 1) + "}"
			default:
				if vv, ok := v.(map[interface{}]interface{}); ok {
					vStr = BuildLuaResponse(vv)
				}
			}
			if vStr != "" {
				switch v := k.(type) {
				case int:
					ret += `,[` + strconv.Itoa(v) + `]=` + vStr
				case string:
					ret += `,` + v + `=` + vStr
				}
			}
		}
		return "{" + Substr(ret, 1) + "}"
	}
	return "{}"
}

// outputLuaMsg 输出lua格式提示信息
func OutputLua(w http.ResponseWriter, m map[interface{}]interface{}) {
	w.Write([]byte("return " + BuildLuaResponse(m)))
}

// outputLuaMsg 输出lua格式提示信息
func OutputLuaMsg(w http.ResponseWriter, msg string, code ...int) bool {
	var c int = 120
	if len(code) >= 1 {
		c = code[0]
	}
	OutputLua(w, NewLuaResult(msg, c))
	return c == 0
}

// outputLuaOk 输出lua格式成功信息
func OutputLuaOk(w http.ResponseWriter, msg string) {
	OutputLuaMsg(w, msg, 0)
}

// 获取lua返回结果map
func NewLuaResult(params ...interface{}) map[interface{}]interface{} {
	ret := make(map[interface{}]interface{})
	if len(params) > 0 {
		ret["msg"] = params[0]

		if len(params) > 1 {
			ret["status"] = params[1]
		} else {
			ret["status"] = 0
		}
	}
	return ret
}