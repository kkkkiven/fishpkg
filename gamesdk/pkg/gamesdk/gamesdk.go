package gamesdk

import (
	"context"
	"fmt"
	"net"
	"time"

	"git.yuetanggame.com/zfish/fishpkg/gamesdk/internal/message"

	ntf "git.yuetanggame.com/zfish/fishpkg/gamesdk/api/proto/notify"
	_ "git.yuetanggame.com/zfish/fishpkg/gamesdk/internal/handler"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/internal/server"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/config"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/errors"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/game"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/types"
	"git.yuetanggame.com/zfish/fishpkg/logs"
	sdk "git.yuetanggame.com/zfish/fishpkg/servicesdk/core"
	"git.yuetanggame.com/zfish/fishpkg/sprotocol/core"
	pb "github.com/golang/protobuf/proto"
)

func GetServer() *server.Server {
	return server.GetServer()
}

// Init 初始化库
func Init(conf config.EntityServer, sdkConf sdk.CConfig) error {
	var err error

	// 服务类型固定是游戏服
	if sdkConf.Type == 0 {
		logs.Waring("[!!!] ]he game type ($sdk.CConfig.Type) is not set, the default type (5) will be used")
		sdkConf.Type = types.GAMESERVER
	}
	// 初始化core
	if err = sdk.Init(&sdkConf); err != nil {
		logs.Error("init serviceSdk-core err:", err)
		return err
	}

	config.Init(conf)

	return err
}

// 启动库，开始业务处理
func Run(gm game.Game) error {
	var err error

	// 启动ServiceSDK
	if err = sdk.Start(); err != nil {
		logs.Error("serviceSdk-core start failed:", err)
		return err
	}

	// 有配置WsPort，则打开WebSocket监听
	if config.GetWsPort() != 0 {
		// 初始化webSocket监听
		if err = server.StartWebSocket(uint16(config.GetWsPort())); err != nil {
			logs.Error("websocket listen failed:", err)
			return err
		}
	}

	// 普通Tcp监听
	var listen net.Listener
	if listen, err = net.Listen("tcp", fmt.Sprintf(":%d", config.GetPort())); err != nil {
		logs.Error("server listen err:", err)
		return err
	}

	logs.Infof("***** [Game: %s] serving on port %d and ws port %d *****\n",
		config.GetGameName(), config.GetPort(), config.GetWsPort())

	// 设置游戏
	game.SetGame(gm)

	// 启动服务
	return server.Start("", listen,
		server.MaxConns(config.GetMaxConns()),
		server.TimeOut(time.Duration(config.GetTimeOut())))
}

// Stop ...
func Stop() {
	server.Stop()
}

// AddHandler 添加消息处理
func AddHandler(msgID uint16, handle server.Handler) {
	server.AddHandler(msgID, handle)
}

// Post 给指定玩家发送数据,不等待回应
func Post(playerID int64, msgID uint16, buf []byte) error {
	if player, err := game.GetPlayer(playerID); err == nil {
		msg := core.NewNormalMessage()
		msg.SetFunctionID(msgID)
		msg.SetBody(buf)
		if player.Conn != nil {
			_, err = player.Conn.Send(context.Background(), msg)
			return err
		}
		return errors.LostConnection
	}
	return errors.PlayerNotFound
}

// 发送消息，并等待回应
// 针对同步一问一答的通信模式
func Send(playerID int64, msgID uint16, buf []byte) (resp []byte, err error) {
	if player, err := game.GetPlayer(playerID); err == nil {
		msg := core.NewRequestMessage()
		msg.SetFunctionID(msgID)
		msg.SetBody(buf)
		if player.Conn != nil {
			if ret, err := player.Conn.Send(context.Background(), msg); err != nil {
				return nil, err
			} else {
				return ret.GetBody(), nil
			}
		}
		return nil, errors.LostConnection
	}
	return nil, errors.PlayerNotFound
}

// 回应客户端消息
// 针对同步一问一答的通信模式
func Reply(playerID int64, requestID uint32, msgID uint16, buf []byte) error {
	if player, err := game.GetPlayer(playerID); err == nil {
		msg := core.NewResponseMessage()
		msg.SetRequestID(requestID)
		msg.SetFunctionID(msgID)
		msg.SetBody(buf)
		if player.Conn != nil {
			_, err = player.Conn.Send(context.Background(), msg)
			return err
		}
		return errors.LostConnection
	}
	return errors.PlayerNotFound
}

// RemovePlayer 删除玩家
func RemovePlayer(playerID int64) {
	game.RemovePlayer(playerID)
}

// GetUserAttr 获取用户基础属性
// uid 用户id
// fields 用户属性集合
func GetUserAttr(tc context.Context, uid int64, attrKeys []string, ext ...map[string][]byte) (map[string]string, error) {
	return sdk.Client().GetUserAttr(tc, uid, attrKeys, ext...)
}

// SetUserAttr 设置用户属性
// uid 用户id
// appid 应用id
// channelid 渠道id
// attr 用户属性集合
func SetUserAttr(tc context.Context, uid int64, appid, channelid int32, attr map[string]string, ext ...map[string][]byte) error {
	return sdk.Client().SetUserAttr(tc, uid, appid, channelid, attr, ext...)
}

// GetUserProp 获取用户基础道具
func GetUserProp(tc context.Context, uid int64, propKeys []int32, ext ...map[string][]byte) (map[int32]int64, error) {
	return sdk.Client().GetUserProp(tc, uid, propKeys, ext...)
}

// SetUserProp 设置用户道具信息
func SetUserProp(tc context.Context, uid int64, appid, channelid int32, prop map[int32]int64, ext ...map[string][]byte) error {
	return sdk.Client().SetUserProp(tc, uid, appid, channelid, prop, ext...)
}

// GetUserInfo 获取用户信息（包括属性和道具）
func GetUserInfo(tc context.Context, uid int64, attrKeys []string, propKeys []int32, advancePropKeys []int32,
	gamename string, gameAttrKeys []string, ext ...map[string][]byte) (map[string]string, map[int32]int64, map[int32]string, string, error) {
	return sdk.Client().GetUserInfo(tc, uid, attrKeys, propKeys, advancePropKeys, gamename, gameAttrKeys, ext...)
}

// SetUserInfo 设置用户信息
func SetUserInfo(tc context.Context, uid int64, appid, channelid int32, attr map[string]string, prop map[int32]int64, ext ...map[string][]byte) error {
	return sdk.Client().SetUserInfo(tc, uid, appid, channelid, attr, prop, ext...)
}

// GetUserGameInfo 获取用户游戏属性
func GetUserGameInfo(tc context.Context, uid int64, gamename string, gameAttrKeys []string, ext ...map[string][]byte) (string, error) {
	return sdk.Client().GetUserGameInfo(tc, uid, gamename, gameAttrKeys, ext...)
}

// SetUserGameInfo 设置用户游戏属性
func SetUserGameInfo(tc context.Context, uid int64, appid, channelid int32, gamename string, data string, ext ...map[string][]byte) error {
	return sdk.Client().SetUserGameInfo(tc, uid, appid, channelid, gamename, data, ext...)
}

// GetUserAProp 获取用户高级道具
func GetUserAProp(tc context.Context, uid int64, aPropKeys []int32, ext ...map[string][]byte) (map[int32]string, error) {
	return sdk.Client().GetUserAProp(tc, uid, aPropKeys, ext...)
}

// SetUserAProp 操作用户高级
// opType 操作类型，如充值、看广告等
// option 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
func SetUserAProp(tc context.Context, uid int64, appid, channelid int32, aProp map[int32]string, ext ...map[string][]byte) error {
	return sdk.Client().SetUserAProp(tc, uid, appid, channelid, aProp, ext...)
}

// EnterGame 进入游戏
func EnterGame(tc context.Context, uid int64, gameid string, gamename, roomname string,
	attrKeys []string, propKeys []int32, aPropKeys []int32,
	ext ...map[string][]byte) (map[string]string, map[int32]int64, map[int32]string, string, string, error) {
	return sdk.Client().EnterGame(tc, uid, gameid, gamename, roomname, attrKeys, propKeys, aPropKeys, ext...)
}

// LeaveGame 离开游戏
func LeaveGame(tc context.Context, uid int64, ext ...map[string][]byte) error {
	return sdk.Client().LeaveGame(tc, uid, ext...)
}

// 敏感词检查
func HasBadWord(tc context.Context, str string, ext ...map[string][]byte) ([]string, error) {
	return sdk.Client().HasBadWord(tc, str, ext...)
}

// 敏感词替换
// return 替换后的新串
func ReplaceBadWord(tc context.Context, str string, ext ...map[string][]byte) (string, error) {
	return sdk.Client().ReplaceBadWord(tc, str, ext...)
}

// PushMatchScore 上传比赛积分
// matchType 比赛类型 1    娜迦千倍（周清）
//             2  娜迦至尊(周清)
//             3  鱼券赛(小时清)
//             4  水晶赛(小时清)
//             5  弹头赛(小时清)
//             6  炸弹乐园(周清)
// expire 过期时间
func PushMatchScore(tc context.Context, playID int64, matchType int, score int, expire int) error {
	return sdk.Client().PushMatchScore(tc, playID, matchType, score, expire)
}

// Broadcast 全游戏平台广播
// 发送全游戏平台广播, 广播内容完全透传
// 要接收广播内容，需实现 Broadcaster接口
// type Broadcaster interface {
//	OnBroadcast(msgBody []byte)
// }
func Broadcast(msgId uint32, body []byte) error {
	return sdk.SendBroadcastFish(context.Background(), 5, message.GWMSGBroadcastGame, msgId, body)
}

// Notify 发送通知到Kafka0
// playerID - 0所有用户 非0代表指定用户
// appID - 0全平台消息，1指定平台消息
// tye - 消息类型 -自定义
func Notify(playerID int64, appID, typ int32, msg []byte, args ...string) error {
	notify := &ntf.NotifyMsg{}
	notify.UserId = playerID
	notify.AppId = appID
	notify.Type = typ
	notify.NotifyMsg = msg
	body, err := pb.Marshal(notify)

	if err != nil {
		logs.Debug("body: %v ", body)
		return err
	}

	topic := "t_notify"
	if len(args) > 0 {
		topic = args[0]
	}
	return sdk.GetService().KafkaProducer().SendMessage(topic, body)
}

// GetCount 获取当前连接和玩家数
// "players" - 当家玩家数
// "socket" - 普通Socket连接数
// "web_sockets" - WebSocket 连接数
func GetCount() map[string]int32 {
	counter := make(map[string]int32, 3)
	counter["players"] = game.GetPlayerCount()
	count, wsCount := server.GetConnCount()
	counter["sockets"] = count
	counter["web_sockets"] = wsCount
	return counter
}
