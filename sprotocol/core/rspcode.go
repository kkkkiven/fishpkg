package core

var mRsp map[int32]string

const (
	RC_OK                = 0x0000 // 成功
	RC_SYS_ERR           = 0x0001 // 系统错误
	RC_TIMEOUT           = 0x0002 // 超时
	RC_HANDLER_NOT_FOUND = 0x0003 // 接口未找到
	RC_HANDLER_PANIC     = 0x0004 // 接口奔溃
)

func init() {
	mRsp = make(map[int32]string, 0)
	mRsp[RC_OK] = "ok"
	mRsp[RC_SYS_ERR] = "system error"
	mRsp[RC_TIMEOUT] = "timeout"
	mRsp[RC_HANDLER_NOT_FOUND] = "handler not found"
	mRsp[RC_HANDLER_PANIC] = "handler panic"
}

// 获取错误码消息值
func M(code int32) string {
	if v, ok := mRsp[code]; ok {
		return v
	}

	return "unknow"
}
