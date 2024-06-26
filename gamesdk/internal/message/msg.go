package message

import (
	"fmt"
	"sync"

	"git.yuetanggame.com/zfish/fishpkg/gamesdk/api/proto/game"
	"github.com/golang/protobuf/proto"
)

// ST_USER_API UserApi服务类型
const ST_USER_API uint16 = 8

const (
	MSGHeartbeat     uint16 = 0x1000 // 心跳
	MSGHeartbeatResp uint16 = 0x1001
	MSGAuthorize     uint16 = 0x1002 // 连接验证
	MSGAuthorizeResp uint16 = 0x1003

	// GWMSGXXX 网关消息定义 0x02 - 0x0FFF

	// 保留消息段 0x1002 - 0x3FFF
)

// GWMSGXXX 网关消息定义 0x01 - 0x0FFF
const (
	GWMSGKickOff             uint16 = 0x01 // 踢玩家下线
	GWMSGOpProp              uint16 = 0x02 // 操作道具
	GWMSGOpGameInfo          uint16 = 0x03 // 修改游戏属性
	GWMSGBroadcastGame       uint16 = 0x04 // 游戏广播
	GWMSGGetProp             uint16 = 0x05 // 获取玩家实时的道具信息
	GWMSGOpSeniorProp        uint16 = 0x06 // 操作高级道具
	GWMSGGetSeniorProp       uint16 = 0x07 // 获取高级道具
	GWMSGGetGameInfo         uint16 = 0x08 // 获取游戏属性
	GWMSGOpPropAndSeniorProp uint16 = 0x09 // 原子操作普通道具和高级道具
	GWMSGGetUserAllData      uint16 = 0x0A // 获取玩家所有数据
	GWMSGGiveUp              uint16 = 0x0B // 玩家放弃当前比赛
	GWMSGSetProp             uint16 = 0x0C // 设置道具
	GWMSGSetSeniorProp       uint16 = 0x0D // 设置高级道具
	GWMSGSetGameInfo         uint16 = 0x0E // 设置游戏属性
)

var messages = map[uint16]*sync.Pool{
	MSGHeartbeat: {
		New: func() interface{} {
			return &game.ReqHeartbeat{}
		},
	},
	MSGHeartbeatResp: {
		New: func() interface{} {
			return &game.RespHeartbeat{}
		},
	},
	MSGAuthorize: {
		New: func() interface{} {
			return &game.ReqAuthorize{}
		},
	},
	MSGAuthorizeResp: {
		New: func() interface{} {
			return &game.RespAuthorize{}
		},
	},
}

func Get(msgID uint16) interface{} {
	if pool, ok := messages[msgID]; ok {
		return pool.Get()
	}
	panic(fmt.Errorf("can not get msg in message pool ,undefined msg:%04x", msgID))
	return nil
}

func Put(msgID uint16, msg proto.Message) {
	msg.Reset()
	if pool, ok := messages[msgID]; ok {
		pool.Put(msg)
	}
}
