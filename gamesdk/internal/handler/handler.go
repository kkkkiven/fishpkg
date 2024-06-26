package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"

	gameproto "git.yuetanggame.com/zfish/fishpkg/gamesdk/api/proto/game"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/internal/message"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/internal/server"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/internal/utils"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/config"
	errs "git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/errors"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/game"
	"git.yuetanggame.com/zfish/fishpkg/logs"
	sdk "git.yuetanggame.com/zfish/fishpkg/servicesdk/core"
	usr "git.yuetanggame.com/zfish/fishpkg/servicesdk/core/pb/userapi"
	pc "git.yuetanggame.com/zfish/fishpkg/sprotocol/core"
	pb "github.com/golang/protobuf/proto"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var responselPool = &sync.Pool{
	New: func() interface{} {
		return pc.NewResponseMessage()
	},
}

var normalPool = &sync.Pool{
	New: func() interface{} {
		return pc.NewNormalMessage()
	},
}

func init() {
	// 客户端消息
	_ = pc.AddHandler(message.MSGHeartbeat, OnHeartbeatHandle)
	_ = pc.AddHandler(message.MSGAuthorize, OnAuthorizeHandle)

	// 网关消息
	sdk.AddHandler(message.GWMSGKickOff, OnKickOffHandle)
	sdk.AddHandler(message.GWMSGGiveUp, OnGiveUpHandle)
	sdk.AddHandler(message.GWMSGBroadcastGame, OnBroadcast)
	sdk.AddHandler(message.GWMSGOpProp, OnOperatePropHandle)
	sdk.AddHandler(message.GWMSGOpGameInfo, OnOperateGameInfoHandle)
	sdk.AddHandler(message.GWMSGOpSeniorProp, OnOperateSeniorPropHandle)
	sdk.AddHandler(message.GWMSGOpPropAndSeniorProp, OnOperatePropAndSPropHandle)
	sdk.AddHandler(message.GWMSGGetProp, OnGetUserProp)
	sdk.AddHandler(message.GWMSGGetGameInfo, OnGetUserGameInfo)
	sdk.AddHandler(message.GWMSGGetSeniorProp, OnGetUserSeniorProp)
	sdk.AddHandler(message.GWMSGGetUserAllData, OnGetUserAllData)
	sdk.AddHandler(message.GWMSGSetProp, OnSetUserProp)
	sdk.AddHandler(message.GWMSGSetSeniorProp, OnSetUserSeniorProp)
	sdk.AddHandler(message.GWMSGSetGameInfo, OnSetUserGameInfo)
}

// OnHeartbeatHandle 心跳
func OnHeartbeatHandle(ctx context.Context, conn *pc.Socket, msg *pc.Message) {
	var resp *pc.Message
	var pool *sync.Pool

	if msg.GetMessageType() == pc.MT_REQUEST {
		pool = responselPool
	} else {
		pool = normalPool
	}

	resp = pool.Get().(*pc.Message)
	resp.SetRequestID(msg.GetRequestID())
	resp.SetFunctionID(message.MSGHeartbeatResp)

	var err error
	var respData = message.Get(message.MSGHeartbeatResp).(*gameproto.RespHeartbeat)

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Errcode = errs.ErrCode(err)
		respBody, _ := pb.Marshal(respData)
		resp.SetBody(respBody)
		if _, err := conn.Send(ctx, resp); err != nil {
			logs.Error(err)
		}

		resp.Reset()
		pool.Put(resp)
	}()

	logs.Debugf(" [heartbeat] - %s", conn.GetRemoteIPStr())
	respData.Timestamp = time.Now().Unix()
}

// 连接授权
func OnAuthorizeHandle(ctx context.Context, conn *pc.Socket, msg *pc.Message) {
	var resp *pc.Message
	var pool *sync.Pool

	if msg.GetMessageType() == pc.MT_REQUEST {
		pool = responselPool
	} else {
		pool = normalPool
	}

	resp = pool.Get().(*pc.Message)
	resp.SetRequestID(msg.GetRequestID())
	resp.SetFunctionID(message.MSGAuthorizeResp)

	var err error
	var respData = message.Get(message.MSGAuthorizeResp).(*gameproto.RespAuthorize)

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Errcode = errs.ErrCode(err)
		respBody, _ := pb.Marshal(respData)
		resp.SetBody(respBody)
		logs.Infof("%+v", resp)
		if _, err := conn.Send(ctx, resp); err != nil {
			logs.Error(err)
		}

		resp.Reset()
		pool.Put(resp)

		if err != nil && err != errs.SUCCESS {
			conn.SetFilter(server.GetServer())
			conn.Close()
		}
	}()

	req := message.Get(message.MSGAuthorize).(*gameproto.ReqAuthorize)
	defer message.Put(message.MSGAuthorize, req)

	err = pb.Unmarshal(msg.GetBody(), req)

	if err != nil {
		logs.Errorf("close client:%s", conn.GetRemoteIPStr(), err)
		err = errs.BadMsg
		return
	}

	uid := req.GetUid()
	if uid == 0 {
		logs.Errorf("close client: uid required", conn.GetRemoteIPStr())
		err = errs.ErrAuthorized
		return
	}

	// 要求签名验证
	if config.Authorize() {
		cliSign := req.GetSign()
		if cliSign == "" {
			logs.Errorf("close client: sign required", conn.GetRemoteIPStr())
			err = errs.ErrAuthorized
			return
		}

		stamp := req.GetTimestamp()
		// 请求时间必须在30分钟内，防止重放攻击
		ts := time.Unix(stamp, 0)
		if time.Since(ts) > 30*time.Minute || time.Since(ts) < -30*time.Minute {
			logs.Errorf("close client: request expires, please check both system time", conn.GetRemoteIPStr())
			err = errs.ErrAuthorized
			return
		}

		nonce := req.GetNonce()

		var buf bytes.Buffer
		buf.WriteString(nonce)
		buf.WriteString(fmt.Sprintf("%d", stamp))
		buf.WriteString(fmt.Sprintf("%d", uid))

		var info map[string]string
		info, err = sdk.Client().GetUserAttr(context.Background(), uid, []string{"token"}, nil)
		if err != nil {
			logs.Errorf("close client:%s", conn.GetRemoteIPStr(), err.Error())
			err = errs.ErrAuthorized
			return
		}

		buf.WriteString(info["token"])

		sign := utils.Sign(buf.String())

		if sign != cliSign {
			logs.Errorf("close client: sign error", conn.GetRemoteIPStr())
			err = errs.ErrAuthorized
			return
		}
	}

	reConn := false // 重连标记
	player := &game.Player{ID: uid, Conn: conn}

	// 签名ok
	// 若玩家已在游戏中，但连接不同，则先断开旧连接
	if p, err := game.GetPlayer(uid); err == nil {
		player.SetContext(p.GetContext())
		reConn = true

		if p.Conn != conn {
			originConn := p.Conn
			p.Conn = conn
			_ = game.UpdatePlayer(uid, p)

			if originConn != nil {
				msg := pc.NewNormalMessage() // 发送顶号处理
				msg.SetFunctionID(0x8005)
				buf, _ := pb.Marshal(&gameproto.NotifyOtherLogin{})
				msg.SetBody(buf)
				_, err = originConn.Send(context.Background(), msg)
				// 把Context设成nil，防止触发player的 Lost
				originConn.SetContext(nil)
				originConn.SetFilter(server.GetServer())
				originConn.Close()
				logs.Waringf("player(%v) is already in the game, original connection has be closed", p)
			}
		}
	} else if err = game.AddPlayer(player); err != nil {
		return
	}

	conn.SetContext(player)
	// 回调游戏
	if err = game.OnConnect(player, reConn); err != nil {
		game.RemovePlayer(uid)
		return
	}
	logs.Infof("player(%v) authorize success", player)
	return
}

// OnKickOffHandle 踢玩家
func OnKickOffHandle(ctx *sdk.SDKContext) {
	var err error
	var respData = &usr.RspKickOff{}

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqKickOff{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		logs.Error(err)
		err = errs.BadMsg
		return
	}

	uid := reqBody.GetId()
	logs.Infof("player(uid:%d) kicked off request: %s", uid, reqBody.Msg)

	// 通知游戏处理玩家被踢，游戏此时应该回应玩家数据，并调用LeaveGame
	game.OnKickOff(reqBody.GetId(), reqBody.GetMsg())
}

// OnGiveUpHandle 玩家放弃当前比赛
func OnGiveUpHandle(ctx *sdk.SDKContext) {
	var err error
	var respData = &usr.RspGiveUp{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(id:%d) give up failed:%v", uid, err)
		} else {
			logs.Infof("player(%v) give up success", uid)
		}

		respBody, _ := pb.Marshal(respData)

		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqGiveUp{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()

	game.OnGiveUp(uid, reqBody.Msg)
}

// OnBroadcast 接收到游戏全局广播消息
func OnBroadcast(ctx *sdk.SDKContext) {
	game.OnBroadcast(ctx.GetMsg().GetRequestID(), ctx.GetBody())
}

// OnOperatePropHandle 购买道具通知
func OnOperatePropHandle(ctx *sdk.SDKContext) {
	var err error
	var respData = &usr.RspNotifyOpProp{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) operate prop failed:%v", uid, errs.Error(err))
		} else {
			logs.Infof("player(uid:%d) operate prop success", uid)
		}

		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	// 解析接收消息体
	reqBody := &usr.ReqNotifyOpProp{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()
	prop := reqBody.GetProp()
	var before, after map[int32]int64

	// 通知到game， game 要返回道具操作前后的值
	before, after, err = game.OnOperateProp(uid, reqBody.GetOption(), reqBody.GetOptype(), prop, reqBody.GetExt())
	if err != nil {
		return
	}

	// 设置道具操作前后值到返回结果
	record := make(map[int32]*usr.RecordProp, len(prop))

	for k := range prop {
		if newVal, ok := after[k]; ok {
			record[k] = &usr.RecordProp{NewProp: newVal, OldProp: before[k]}
		}
	}
	respData.Prop = record
}

// OnOperateGameInfoHandle 外部操作游戏属性，如锻造
func OnOperateGameInfoHandle(ctx *sdk.SDKContext) {
	var err error
	var respData = &usr.RspNotifyOpGame{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) operate game info failed:%v", uid, err)
		} else {
			logs.Infof("player(uid:%d) operate game info success", uid)
		}
		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	// 解析消息体
	reqBody := &usr.ReqNotifyOpGame{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()

	// 通知到game
	err = game.OnOperateGameInfo(uid, reqBody.GetName(), reqBody.GetOption(), reqBody.GetOptype(), reqBody.GetInfo(), reqBody.GetExt())
}

// OnOperateSeniorPropHandle 外部操作高级道具
func OnOperateSeniorPropHandle(ctx *sdk.SDKContext) {
	var err = errs.SUCCESS
	var respData = &usr.RspNotifyOpAdProp{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) operate senior prop failed:%v", uid, err)
		} else {
			logs.Infof("player(uid:%d) operate senior prop success", uid)
		}
		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqNotifyOpAdProp{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()
	// 通知到game
	err = game.OnOperateSeniorProp(uid, reqBody.GetOption(), reqBody.GetOptype(), reqBody.GetAdProp(), reqBody.GetExt())
}

// OnOperatePropAndSPropHandle 原子操作普通道具和高级道具
func OnOperatePropAndSPropHandle(ctx *sdk.SDKContext) {
	var err = errs.SUCCESS
	var respData = &usr.RspOpPropAndAdProp{}
	var playerID int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) operate prop and senior prop failed:%v", playerID, err)
		} else {
			logs.Infof("player(uid:%d) operate prop and senior prop success", playerID)
		}
		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqOpPropAndAdProp{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		err = errs.BadMsg
		return
	}

	playerID = reqBody.GetId()

	respData.Prop, respData.AdProp, err = game.OnOperatePropAndSeniorProp(playerID,
		reqBody.GetOption(),
		reqBody.GetOptype(),
		reqBody.GetProp(),
		reqBody.GetAdProp(),
		reqBody.GetExt())
}

// OnGetUserProp 实时获取玩家道具
func OnGetUserProp(ctx *sdk.SDKContext) {
	var err = errs.SUCCESS
	var respData = &usr.RspProp{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) get user prop failed:%v", uid, err)
		}

		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqProp{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		logs.Error(err)
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()
	propIDS := reqBody.GetProp()

	var props map[int32]int64
	props, err = game.GetProp(uid, propIDS)
	if err == nil {
		respData.Prop = props

		logs.Infof("player(uid:%d) get user prop:%+v success", uid, respData.GetProp())
		return
	}
}

// OnGetUserSeniorProp 获取玩家高级道具
func OnGetUserSeniorProp(ctx *sdk.SDKContext) {
	var err = errs.SUCCESS
	var respData = &usr.RspAdProp{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) get user senior prop failed:%v", uid, err)
		}
		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqAdProp{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		logs.Error(err)
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()

	respData.AdProp, err = game.GetSeniorProp(uid, reqBody.GetAdProp())

	if err == nil {
		logs.Infof("player(uid:%d) get user senior prop:%+v success", uid, respData.GetAdProp())
	}
}

// OnGetUserGameInfo 实时获取玩家游戏属性
func OnGetUserGameInfo(ctx *sdk.SDKContext) {
	var err error
	var respData = &usr.RspGameInfo{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) get user game info failed:%v", uid, err)
		}
		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqGameInfo{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		logs.Error(err)
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()
	gameName := reqBody.GetName()

	if gameName == config.GetGameName() {
		var mInfo map[string]interface{}
		mInfo, err = game.GetGameInfo(uid, reqBody.GetFields())

		if err == nil {
			var buf []byte
			buf, err = json.Marshal(&mInfo)
			if err == nil {
				respData.Info = utils.Bytes2String(buf)
			} else {
				logs.Error(err)
			}
			logs.Infof("player(uid:%d) get game info:%+v success", uid, mInfo)
		}
		return
	}

	err = errs.ParamInvalid
	logs.Errorf("player(uid:%d) get game info failed: game name must be %s ,got %s", uid, config.GetGameName(), gameName)
}

// OnGetUserAllData 获取玩家数据
func OnGetUserAllData(ctx *sdk.SDKContext) {
	var err error
	var respData = &usr.RspAll{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) get user all data failed:%v", uid, err)
		}
		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqAll{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		logs.Error(err)
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()
	propID := reqBody.GetProp()
	sPropID := reqBody.GetAdvancedProp()
	gameName := reqBody.GetGameName()
	gameField := reqBody.GetGameFields()

	if len(propID) > 0 {
		// 获取普通道具
		respData.Prop, err = game.GetProp(uid, propID)
		if err != nil {
			logs.Errorf("player(uid:%d) game.GetProp failed:%v", uid, err)
			return
		}
	}

	// 获取高级道具
	if len(sPropID) > 0 {
		respData.AdvancedProp, err = game.GetSeniorProp(uid, sPropID)
		if err != nil {
			logs.Errorf("player(uid:%d) game.GetSeniorProp failed:%v", uid, err)
			return
		}
	}

	if len(gameField) > 0 {
		if gameName == config.GetGameName() {
			mGameInfo, mGameInfoErr := game.GetGameInfo(uid, gameField)
			if mGameInfoErr != nil {
				err = mGameInfoErr
				logs.Errorf("player(uid:%d) game.GetGameInfo failed:%v", uid, err)
				return
			}
			if len(mGameInfo) > 0 {
				buf, marErr := json.Marshal(&mGameInfo)
				if marErr == nil {
					respData.GameInfo = utils.Bytes2String(buf)
				} else {
					err = marErr
					logs.Error(err)
				}
			}
			return
		}
		logs.Waringf("player(uid:%d) get game info failed: game name must be %s ,got %s", uid, config.GetGameName(), gameName)
	}
}

// OnSetUserProp 设置玩家数据
func OnSetUserProp(ctx *sdk.SDKContext) {
	var err error
	var respData = &usr.RspSetProp{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) set user prop failed:%v", uid, err)
		}
		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqSetProp{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		logs.Error(err)
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()

	err = game.OnSetProp(uid, reqBody.GetProp())
}

// OnSetUserSeniorProp 设置玩家高级道具
func OnSetUserSeniorProp(ctx *sdk.SDKContext) {
	var err error
	var respData = &usr.RspSetAdProp{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) set user senior prop failed:%v", uid, err)
		}
		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqSetAdProp{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		logs.Error(err)
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()

	err = game.OnSetSeniorProp(uid, reqBody.GetAdProp())
}

// OnSetUserGameInfo 设置玩家游戏属性
func OnSetUserGameInfo(ctx *sdk.SDKContext) {
	var err error
	var respData = &usr.RspSetGameInfo{}
	var uid int64

	defer func() {
		respData.Msg = errs.Error(err)
		respData.Code = errs.ErrCode(err)

		if err != nil && err != errs.SUCCESS {
			logs.Errorf("player(uid:%d) set user game info failed:%v", uid, err)
		}
		respBody, _ := pb.Marshal(respData)
		if err := ctx.SendResponse(respBody); err != nil {
			logs.Error(err)
		}
	}()

	reqBody := &usr.ReqSetGameInfo{}
	err = pb.Unmarshal(ctx.GetBody(), reqBody)
	if err != nil {
		logs.Error(err)
		err = errs.BadMsg
		return
	}

	uid = reqBody.GetId()
	gameName := reqBody.GetName()
	info := reqBody.GetInfo()

	if gameName != config.GetGameName() {
		logs.Errorf("set user gameInfo failed:GameName must be %s, but got %s", config.GetGameName(), gameName)
		err = errors.New("game name error")
		return
	}

	gameInfo := make(map[string]interface{})
	if err = json.Unmarshal([]byte(info), &gameInfo); err != nil {
		return
	}

	err = game.OnSetGameInfo(uid, gameInfo)
}
