package game

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kkkkiven/fishpkg/gamesdk/pkg/errors"
	"github.com/kkkkiven/fishpkg/logs"
)

type Broadcaster interface {
	OnBroadcast(msgID uint32, msgBody []byte)
}

type GiveUper interface {
	// 玩家放弃当前比赛
	OnGiveUp(playerID int64, msg /*消息*/ string)
}

type UserInfoSetter interface {
	// 外部设置游戏属性
	OnSetGameInfo(playerID int64, gameInfo map[string]interface{}) error

	// 外部设置普通道具
	OnSetProp(playerID int64, prop map[int32]int64) error

	// 外部设置高级道具
	OnSetSeniorProp(playerID int64, prop map[int32]string) error
}

type Game interface {
	// 消息处理
	OnMessage(player *Player, frameID uint32, msgID uint16, msgBody []byte)

	// 玩家连接
	// reConn 标记玩家是否是重连
	// 返回值 nil 表示允许玩家进入游戏， 否则将阻止玩家进入，并断开连接
	OnConnect(player *Player, reConn /*重连标记*/ bool) error

	// 连接断开
	OnPlayerLost(player *Player)

	// 玩家被踢通知
	OnKickOff(playerID /*玩家id*/ int64, msg /*消息*/ string)

	// 返回指定的道具信息
	GetProp(playerID int64, propIDs /*要获取的道具id列表, 列表为空表示返回全部*/ []int32) (prop map[int32]int64, err error)

	// 返回指定的高级道具信息
	GetSeniorProp(playerID int64, propIDs /*要获取的道具id列表, 列表为空表示返回全部*/ []int32) (prop map[int32]string, err error)

	// 返回游戏属性
	GetGameInfo(playerID int64, fields []string) (prop map[string]interface{}, err error)

	// 外部操作普通道具
	// 参数 data，结构如下
	//	data["optype"] - 操作类型，如充值、看广告等
	//	data["option"] - 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
	//	data["prop"] - 道具列表 (增量)
	//	data["ext"] - 扩展
	// 返回值：
	// before -操作前道具数量
	// after - 操作后道具数量
	OnOperateProp(playerID int64, option int32, opType string, data map[int32]int64, ext ...map[string][]byte) (before map[int32]int64, after map[int32]int64, err error)

	// 外部操作高级道具
	OnOperateSeniorProp(playerID int64, option int32, opType string, data map[int32]string, ext ...map[string][]byte) error

	// 外部操作普通道具和高级道具
	OnOperatePropAndSeniorProp(playerID int64, option int32, opType string,
		prop map[int32]int64, seniorProp map[int32]string, ext ...map[string][]byte) (
		afterProp map[int32]int64, /*操作后普通道具数量*/
		afterSeniorProp map[int32]string, /*操作后高级道具数量*/
		err error)

	// 外部操作游戏属性
	OnOperateGameInfo(playerID int64, gameName string, option int32, opType string, data map[string]int64, ext ...map[string][]byte) error
}

var gm Game

// 玩家列表
var players sync.Map
var playerCount int32

func init() {
	players = sync.Map{}
}

// GetPlayer 获取玩家
func GetPlayer(playerID int64) (Player, error) {
	if v, ok := players.Load(playerID); ok {
		player := v.(*Player)
		return *player, nil
	}
	return Player{}, errors.PlayerNotFound
}

// 获取当前玩家数
func GetPlayerCount() int32 {
	return atomic.LoadInt32(&playerCount)
}

// AddPlayer 增加player,原来不存在则人数+1
func AddPlayer(player *Player) (err error) {
	if _, ok := players.Load(player.ID); !ok {
		// 安排桌子的座次
		if assign != nil {
			if player.DeskID, player.SeatID, err = assign.Assign(); err != nil {
				return
			}
		}
		atomic.AddInt32(&playerCount, 1)

	}
	players.Store(player.ID, player)
	return
}

// UpdatePlayer 更新Player
func UpdatePlayer(playerID int64, player Player) error {
	if v, ok := players.Load(playerID); ok {
		p, _ := v.(*Player)
		p.Conn = player.Conn
		p.context = player.context

		return nil
	}

	return errors.PlayerNotFound
}

// 删除玩家，并Close连接；不会触发 PlayerLost事件
func RemovePlayer(playerID int64) {
	if p, err := GetPlayer(playerID); err == nil {
		if p.Conn != nil {
			p.Conn.SetContext(nil)
			p.Conn.Close()
		}

		players.Delete(playerID)

		// 收回桌子的座次
		if assign != nil {
			assign.Recycle(p.DeskID, p.SeatID)
		}

		atomic.AddInt32(&playerCount, -1)
	}
}

// 减玩家数量
func SubPlayerNum(playerID int64) {
	players.Delete(playerID)
	atomic.AddInt32(&playerCount, -1)
}

// 设置自定义数据
func SetContext(playerID int64, context interface{}) error {
	if v, ok := players.Load(playerID); ok {
		player := v.(*Player)
		player.SetContext(context)
		return nil
	}
	return errors.PlayerNotFound
}
func SetGame(gm_ Game) {
	gm = gm_
}

// 玩家进入
func OnConnect(player *Player, reConn bool) error {
	return gm.OnConnect(player, reConn)
}

// 消息处理
func OnMessage(player *Player, requestID /*请求包ID*/ uint32, msgID uint16, msgBody []byte) {
	gm.OnMessage(player, requestID, msgID, msgBody)
}

// 连接断开
func OnPlayerLost(player *Player) {
	gm.OnPlayerLost(player)
}

// 玩家被踢通知
func OnKickOff(uid int64, msg string) {
	gm.OnKickOff(uid, msg)
}

func OnGiveUp(uid int64, msg string) {
	if g, ok := gm.(GiveUper); ok {
		g.OnGiveUp(uid, msg)
		return
	}
	logs.Error("'OnGiveUp' function not be found in game instance")
}

func OnBroadcast(msgID uint32, body []byte) {
	if g, ok := gm.(Broadcaster); ok {
		g.OnBroadcast(msgID, body)
	}
}

// 外部操作普通道具
// 参数 data，结构如下
//
//	data["optype"] - 操作类型，如充值、看广告等
//	data["option"] - 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
//	data["prop"] - 道具列表 (增量)
//	data["ext"] - 扩展
func OnOperateProp(playerID int64, option int32, opType string, data map[int32]int64, ext ...map[string][]byte) (before map[int32]int64, after map[int32]int64, err error) {
	return gm.OnOperateProp(playerID, option, opType, data, ext...)
}

// 外部操作高级道具
func OnOperateSeniorProp(playerID int64, option int32, opType string, data map[int32]string, ext ...map[string][]byte) error {
	return gm.OnOperateSeniorProp(playerID, option, opType, data, ext...)
}

// 外部操作普通道具和高级道具
func OnOperatePropAndSeniorProp(playerID int64, option int32, opType string,
	prop map[int32]int64, seniorProp map[int32]string, ext ...map[string][]byte) (
	afterProp map[int32]int64, afterSeniorProp map[int32]string, err error) {
	return gm.OnOperatePropAndSeniorProp(playerID, option, opType, prop, seniorProp, ext...)
}

// 外部操作游戏属性
func OnOperateGameInfo(playerID int64, gameName string, option int32, opType string, data map[string]int64, ext ...map[string][]byte) error {
	return gm.OnOperateGameInfo(playerID, gameName, option, opType, data, ext...)
}

// 返回游戏属性
func GetGameInfo(playerID int64, fields []string) (info map[string]interface{}, err error) {
	return gm.GetGameInfo(playerID, fields)
}

// 返回指定道具信息
func GetProp(playerID int64, ids []int32) (prop map[int32]int64, err error) {
	return gm.GetProp(playerID, ids)
}

// 返回指定高级道具
func GetSeniorProp(playerID int64, ids []int32) (prop map[int32]string, err error) {
	return gm.GetSeniorProp(playerID, ids)
}

func OnSetProp(playerID int64, prop map[int32]int64) error {
	if setter, ok := gm.(UserInfoSetter); ok {
		return setter.OnSetProp(playerID, prop)
	}
	return fmt.Errorf("%T does not implement UserInfoSetter (missing OnSetProp method)", gm)
}

func OnSetSeniorProp(playerID int64, sProp map[int32]string) error {
	if setter, ok := gm.(UserInfoSetter); ok {
		return setter.OnSetSeniorProp(playerID, sProp)
	}
	return fmt.Errorf("%T does not implement UserInfoSetter (missing OnSetSeniorProp method)", gm)
}

func OnSetGameInfo(playerID int64, gameInfo map[string]interface{}) error {
	if setter, ok := gm.(UserInfoSetter); ok {
		return setter.OnSetGameInfo(playerID, gameInfo)
	}
	return fmt.Errorf("%T does not implement UserInfoSetter (missing OnSetGameInfo method)", gm)
}
