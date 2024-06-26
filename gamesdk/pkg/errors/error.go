package errors

import (
	"errors"
	"fmt"
)

// ErrorCode 错误码定义格式如下
// xx          xx        xxx
// 二位服务号   二位模块号 三位错误码
type ErrorCode int32

var (
	SUCCESS               = errors.New("ok")
	ErrAuthorized         = errors.New("unauthorized request") // 验证失败
	Unknown               = errors.New("unknown error")
	System                = errors.New("system error")
	Server                = errors.New("server disable")
	ConnGateWay           = errors.New("communicate with core gateway error")
	Media                 = errors.New("unsupported media types")
	ParamInvalid          = errors.New("invalid parameters, please refer to the document") // 参数无效，可能是类型错误或者超出范围
	Busy                  = errors.New("system busy")
	TimeOut               = errors.New("time out")
	MissParam             = errors.New("missing required parameters, please refer to the document")
	BadMsg                = errors.New("bad message")
	RemoteSvr             = errors.New("some errors occurred in the remote service") // 跟远程服务通信成功，但远程服务执行逻辑失败
	PlayerInOtherRoom     = errors.New("the player has already in other game")       // 玩家正在其他房间
	PhoneNo               = errors.New("invalid cell phone number")
	Email                 = errors.New("invalid email")
	AssignDesk            = errors.New("some errors occurred when assignDesk")
	UnknownBroadcastBound = errors.New("unknown broadcast bound")                                          // 未知的广播范围
	PlayerOffline         = errors.New("the player is currently offline")                                  // 玩家当前离线
	NotSupportMsg         = errors.New("message not support in the game ,please check the game configure") // 不支持的消息
	MatchRegRepeat        = errors.New("the player has already registered in the game")
	PlayerJoin            = errors.New("the player join failed, please trying again later") // 玩家JoinGame失败
	PlayerJoinFirst       = errors.New("please join game first")                            // 发游戏消息前要先JoinGame
	PlayerNotFound        = errors.New("player not found, please check the player id")      // 指定桌子不存在
	DeskFull              = errors.New("have no more desks to assign") // 没有空闲桌子
	PlayerNotInDesk       = errors.New("the player not in this assign")              // 玩家不在当前桌子
	PlayerNotInGame       = errors.New("the player is not in game")                // 玩家不在游戏中
	LostConnection        = errors.New("the player is currently lost connection")  // 玩家断线
	ResourceLack          = errors.New("lack of resources")                        // 指定资源不足
	PropsLack             = errors.New("prop not enough")                          // 道具不足
	MatchStarted          = errors.New("the match has started")
	NotInTime             = errors.New("not in time")
	TimesLimit            = errors.New("maximum number of times")
	NotRegister           = errors.New("player o register match") // 未报名，比赛

)

var errs = map[error]int32{
	SUCCESS:               0,
	ErrAuthorized:         500401,
	Unknown:               501000,
	System:                501001,
	Server:                501002,
	ConnGateWay:           501003,
	Media:                 501004,
	ParamInvalid:          501005,
	Busy:                  501006,
	TimeOut:               501007,
	MissParam:             501008,
	BadMsg:                501009,
	RemoteSvr:             501010,
	PlayerInOtherRoom:     502000,
	PhoneNo:               502001,
	Email:                 502002,
	AssignDesk:            502003,
	UnknownBroadcastBound: 502004,
	PlayerOffline:         502005,
	NotSupportMsg:         502006,
	MatchRegRepeat:        502007,
	PlayerJoin:            502008,
	PlayerJoinFirst:       502009,
	PlayerNotFound:        502010,
	DeskFull:              502011,
	PlayerNotInDesk:       502014,
	LostConnection:        502015,
	PlayerNotInGame:       502016,
	ResourceLack:          503001,
	PropsLack:             503002,
	MatchStarted:          504001,
	NotInTime:             504002,
	TimesLimit:            504003,
	NotRegister:           504004,
}

// ErrCode 获取错误码
func ErrCode(err error) int32 {
	if err == nil {
		return 0
	}

	if code, ok := errs[err]; ok {
		return code
	}

	return errs[System]
}

func Error(err error) string {
	if err == nil {
		return SUCCESS.Error()
	}

	if code, ok := errs[err]; ok {
		return fmt.Sprintf("(%d)%v", code, err)
	}

	return fmt.Sprintf("(%d)%v", errs[System], err)
}
