package http

import (
	"github.com/kkkkiven/fishpkg/logs"
	p "github.com/kkkkiven/fishpkg/sprotocol/http"
)

type Filter struct {
}

// 接收到裸消息回调, 当不希望再次调用OnRecv时返回true
func (this *Filter) OnRecv(so *p.Socket, data []byte) bool {
	return false
}

// 响应心跳请求消息
func (this *Filter) OnMessage(so *p.Socket, msg *p.Message) bool {
	if msg.GetMessageType() == p.MT_PING {
		rsp := p.NewPongMessage()
		rsp.SetRequestID(msg.GetRequestID())
		rsp.SetBody(msg.GetBody())
		so.Send(rsp)
		logs.Tracef("- %v - Receive PING message", so.GetConn().RemoteAddr().String())
		return true
	}

	return false
}
