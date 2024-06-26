package game

import (
	"fmt"

	"git.yuetanggame.com/zfish/fishpkg/sprotocol/core"
)

// Player 定义了玩家基本数据,和底层连接和对应关系
type Player struct {
	ID      int64 // 玩家ID
	DeskID  int32 // 玩家所在桌号
	SeatID  int32 // 玩家所在座次
	Conn    *core.Socket
	context interface{} // 自定义数据
}

// Reset 重置Player，以便对象复用
func (p *Player) Reset() {
	p.ID = 0
	p.Conn = nil
	p.context = nil
}

func (p *Player) GetContext() interface{} {
	return p.context
}

func (p *Player) SetContext(context interface{}) {
	p.context = context
}

func (p Player) String() string {
	if p.ID == 0 {
		return "unknown"
	}
	return fmt.Sprintf("id:%d", p.ID)
}
