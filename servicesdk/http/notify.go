package http

import (
	"encoding/binary"
	"time"

	"git.yuetanggame.com/zfish/fishpkg/logs"
	p "git.yuetanggame.com/zfish/fishpkg/sprotocol/http"
)

type Notify struct {
	ping int
}

// 关闭tcp连接回调
func (this *Notify) OnClose(so *p.Socket) {
	logs.Errorf("- %v - Close and reconnection", so.GetConn().RemoteAddr().String())

	if key, ok := so.GetContext().(string); ok {
		gwList.Add(key)
		time.Sleep(1 * time.Second)
	}

}

// 读取数据超时回调
func (this *Notify) OnTimeout(so *p.Socket) {
	go func() {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(time.Now().UnixNano()))
		msg := p.NewPingMessage()
		msg.SetBody(buf)
		rsp, err := so.Send(msg)
		if rsp != nil {
			this.ping = 0
			logs.Tracef("- %s - Send PING message", so.GetConn().RemoteAddr().String())
			return
		}

		this.ping++
		logs.Errorf("- %v - Send PING message failed, err:%s", so.GetConn().RemoteAddr().String(), err.Error())

		if this.ping > 2 {
			so.Close()
		}
	}()
}
