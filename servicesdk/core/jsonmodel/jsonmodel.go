package jsonmodel

// 公用响应消息
type RespCommon struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}
