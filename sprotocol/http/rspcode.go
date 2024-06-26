package http

var mRsp map[uint8]string

const (
	RC_OK                = 0x00 // 成功
	RC_SYS_ERR           = 0x01 // 系统错误
	RC_HANDLER_NOT_FOUND = 0x02 // 接口未找到
	RC_HANDLER_PANIC     = 0x03 // 接口panic
	RC_NETWORK_ERR       = 0x04 // 网络错误
	RC_TIMEOUT           = 0x05 // 超时
	RC_SVR_NOT_FOUND     = 0x06 // 服务未找到
	RC_HYSTRIX_LIMIT     = 0x07 // 熔断限制
)

func init() {
	mRsp = make(map[uint8]string, 0)
	mRsp[RC_OK] = "ok"
	mRsp[RC_SYS_ERR] = "system error"
	mRsp[RC_HANDLER_NOT_FOUND] = "handler not found"
	mRsp[RC_HANDLER_PANIC] = "handler panic"
	mRsp[RC_NETWORK_ERR] = "network err"
	mRsp[RC_TIMEOUT] = "timeout"
	mRsp[RC_SVR_NOT_FOUND] = "service not found"
	mRsp[RC_HYSTRIX_LIMIT] = "hystrix limited"
}

// 获取错误码消息值
func M(code uint8) string {
	if v, ok := mRsp[code]; ok {
		return v
	}
	return "unknow"
}
