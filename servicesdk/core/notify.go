package core

import (
	"time"

	"git.yuetanggame.com/zfish/fishpkg/logs"
	p "git.yuetanggame.com/zfish/fishpkg/sprotocol/core"
)

type Notify struct {
	ping int
}

// 关闭tcp连接回调
func (this *Notify) OnClose(so *p.Socket) {
	logs.Waringf("- %s - Reconnection", so.GetConn().RemoteAddr().String())

	if key, ok := so.GetContext().(string); ok {
		gwList.Add(key)
		time.Sleep(1 * time.Second)
	}
}

// 读取数据超时回调
func (this *Notify) OnTimeout(so *p.Socket) {
	go func() {
		msg := p.NewRequestMessage()
		msg.SetToSvrType(ST_GW_CORE)
		msg.SetFunctionID(F_ID_PING)
		if _, err := so.Send(nil, msg); err != nil {
			logs.Errorf("- %s - Send PING message failed", so.GetConn().RemoteAddr().String())

			this.ping++
			if this.ping > 2 {
				so.Close()
			}

			return
		}

		this.ping = 0
		logs.Tracef("- %s - Send PING message success", so.GetConn().RemoteAddr().String())
	}()
}
