package core

import (
	p "git.yuetanggame.com/zfish/fishpkg/sprotocol/core"
)

type Filter struct {
}

// 接收到裸消息回调, 当不希望再次调用OnRecv时返回true
func (this *Filter) OnRecv(so *p.Socket, data []byte) bool {
	return false
}

// 响应心跳请求消息
func (this *Filter) OnMessage(so *p.Socket, msg *p.Message) bool {
	return false
}
