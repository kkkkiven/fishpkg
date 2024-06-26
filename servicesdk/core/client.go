package core

import (
	"context"
	"errors"
	"fmt"

	"git.yuetanggame.com/zfish/fishpkg/logs"
	"git.yuetanggame.com/zfish/fishpkg/servicesdk/core/jsonmodel"
	pbSMS "git.yuetanggame.com/zfish/fishpkg/servicesdk/core/pb/sms"
	pbUA "git.yuetanggame.com/zfish/fishpkg/servicesdk/core/pb/userapi"
	"git.yuetanggame.com/zfish/fishpkg/servicesdk/pkg/utils"
	pCore "git.yuetanggame.com/zfish/fishpkg/sprotocol/core"

	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var errRemoteServer = errors.New("some errors occurred in the remote service")

var json jsoniter.API

func init() {
	extra.RegisterFuzzyDecoders()
	json = jsoniter.ConfigCompatibleWithStandardLibrary
}

type ISdk interface {
	/**
	** USERAPI 接口
	 */
	// GetUserAttr 获取用户基础属性
	// tc 链路上下文（初始链路节点传context.TODO()）
	// uid 用户id
	// attrKeys 属性key集合
	// ext 扩展字段，暂时无用
	GetUserAttr(tc context.Context, uid int64, attrKeys []string, ext ...map[string][]byte) (map[string]string, error)

	// GetUsersAttr 批量获取用户基础属性
	// tc 链路上下文（初始链路节点传context.TODO()）
	// uid 用户id
	// uids 需要获取用户用户ids
	// attrKeys 属性key集合
	// ext 扩展字段，暂时无用
	GetUsersAttr(tc context.Context, id int64, uids []int64, attrKeys []string, ext ...map[string][]byte) (map[int64]map[string]string, error)

	// GetUsersAttrAndGameAttr 批量获取用户基础属性&游戏属性
	// tc 链路上下文（初始链路节点传context.TODO()）
	// uid 用户id
	// uids 需要获取用户用户ids
	// attrKeys 属性key集合
	// gameAttrKeys 游戏属性key集合
	// ext 扩展字段，暂时无用
	GetUsersAttrAndGameAttr(tc context.Context, id int64, uids []int64,
		attrKeys, gameAttrKeys []string, gamename string, ext ...map[string][]byte) (map[int64]map[string]string, map[int64]string, error)

	// SetUserAttr 设置用户基础属性
	// uid 用户id
	// appid 应用id
	// channelid 渠道id
	// attr 要设置的属性kv
	SetUserAttr(tc context.Context, uid int64, appid, channelid int32, attr map[string]string, ext ...map[string][]byte) error

	// GetUserProp 获取用户基础道具
	// uid 用户id
	// propKeys 道具key集合
	GetUserProp(tc context.Context, uid int64, propKeys []int32, ext ...map[string][]byte) (map[int32]int64, error)

	// GetBatchUserInfoAndProp 批量获取用户信息和基础道具
	// uid 请求用户id
	// uids 需要获取用户id集合
	// attrKeys 基础信息集合
	// propKeys 道具key集合
	GetBatchUserInfoAndProp(tc context.Context, uid int64, uids []int64, attrKeys []string, propKeys []int32, ext ...map[string][]byte) (map[int64]map[string]string, map[int64]map[int32]int64, error)

	// GetBatchUserInfosAndProp 批量获取用户信息和基础道具和游戏属性
	// uid 请求用户id
	// uids 需要获取用户id集合
	// attrKeys 基础信息集合
	// propKeys 道具key集合
	// gameKeys 游戏信息集合
	GetBatchUserInfosAndProp(tc context.Context, uid int64, uids []int64, attrKeys []string, propKeys []int32, gameAttrKeys []string, gamename string, ext ...map[string][]byte) (
		map[int64]map[string]string, map[int64]map[int32]int64, map[int64]string, error)

	// SetUserProp 设置用户道具信息
	SetUserProp(tc context.Context, uid int64, appid, channelid int32, prop map[int32]int64, ext ...map[string][]byte) error

	// GetUserInfo 获取用户信息（包括属性和道具）
	// uid 用户id
	// attrKeys 要获取的基础属性key集合
	// propKeys 道具key集合
	// advancePropKeys 高级道具key集合
	// gamename 游戏名称
	// gameAttrKeys 游戏属性key集合
	GetUserInfo(tc context.Context, uid int64, attrKeys []string, propKeys []int32, advancePropKeys []int32,
		gamename string, gameAttrKeys []string, ext ...map[string][]byte) (map[string]string,
		map[int32]int64, map[int32]string, string, error)

	// SetUserInfo 设置用户信息(此接口仅支持同时操作基础属性和基础道具，不支持高级道具和游戏属性)
	SetUserInfo(tc context.Context, uid int64, appid, channelid int32, attr map[string]string, prop map[int32]int64, ext ...map[string][]byte) error

	// OpProp 道具操作
	// opType 操作类型，如充值、看广告等
	// option 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
	OpProp(tc context.Context, uid int64, appid, channelid int32, opType string, option int32,
		prop map[int32]int64, ext ...map[string][]byte) (map[int32]int64, error)

	// SetAll 设置用户属性、普通道具、高级道具
	SetAll(tc context.Context, uid int64, appid, channelid int32, attr map[string]string, prop map[int32]int64, adProp map[int32]string, ext ...map[string][]byte) error

	// GetUserGameInfo 获取用户游戏属性
	GetUserGameInfo(tc context.Context, uid int64, gamename string, gameAttrKeys []string, ext ...map[string][]byte) (string, error)

	// SetUserGameInfo 设置用户游戏属性
	SetUserGameInfo(tc context.Context, uid int64, appid, channelid int32, gamename string, data string, ext ...map[string][]byte) error

	// OpUserGameInfo 操作用户游戏属性
	// opType 操作类型，如充值、看广告等
	// option 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
	OpUserGameInfo(tc context.Context, uid int64, appid, channelid int32, opType string, option int32,
		gamename string, data map[string]int64, ext ...map[string][]byte) (map[string]int64, error)

	// GetUserAProp 获取用户高级道具
	GetUserAProp(tc context.Context, uid int64, aPropKeys []int32, ext ...map[string][]byte) (map[int32]string, error)

	// SetUserAProp 操作用户高级
	// opType 操作类型，如充值、看广告等
	// option 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
	SetUserAProp(tc context.Context, uid int64, appid, channelid int32, aProp map[int32]string, ext ...map[string][]byte) error

	// OpUserAProp 操作用户高级
	// opType 操作类型，如充值、看广告等
	// option 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
	OpUserAProp(tc context.Context, uid int64, appid, channelid int32, opType string, option int32,
		aProp map[int32]string, ext ...map[string][]byte) (map[int32]string, error)

	// OpPropAndAdProp 原子操作普通道具与高级道具
	OpPropAndAdProp(tc context.Context, uid int64, appid, channelid int32, opType string, option int32,
		prop map[int32]int64, adProp map[int32]string, ext ...map[string][]byte) (map[int32]int64, map[int32]string, error)

	// EnterGame 进入游戏
	// uid 用户id
	// gameid 游戏（房间）id
	EnterGame(tc context.Context, uid int64, gameid string, gamename, roomname string,
		attrKeys []string, propKeys []int32, aPropKeys []int32,
		ext ...map[string][]byte) (map[string]string, map[int32]int64, map[int32]string, string, string, error)

	// LeaveGame 离开游戏
	LeaveGame(tc context.Context, uid int64, ext ...map[string][]byte) error

	// IsUsernameReg 用户名是否已注册
	IsUsernameReg(tc context.Context, username string, ext ...map[string][]byte) (bool, error)

	// SetPassword 重置密码
	SetPassword(tc context.Context, uid int64, pwd string, ext ...map[string][]byte) error

	// UpdatePassword 更新密码
	UpdatePassword(tc context.Context, uid int64, oldPwd, newPwd string, ext ...map[string][]byte) error

	// HasBadWord 敏感词检测
	// return 敏感词串，空值表示无敏感词
	HasBadWord(tc context.Context, str string, ext ...map[string][]byte) ([]string, error)

	// 敏感词替换
	// return 替换后的新串
	ReplaceBadWord(tc context.Context, str string, ext ...map[string][]byte) (string, error)

	// CheckIdCardExisted 检查用户idcard是否存在
	CheckIdCardExisted(tc context.Context, idcard string, ext ...map[string][]byte) (bool, error)

	// GetUsersByPhoneNumber 通过电话号码获取用户列表
	GetUsersByPhoneNumber(tc context.Context, phone string, ext ...map[string][]byte) ([]*pbUA.RspGetUsersByPhoneNumberUser, error)

	// AddDeputyAccount 添加辅助登录账号
	AddDeputyAccount(tc context.Context, userId int64, username, password string, userfrom ...int32) error

	// DelDeputyAccount 删除辅助登录账号
	DelDeputyAccount(tc context.Context, userId int64, username string, userfrom ...int32) (bool, error)

	// GetDeputyAccounts 获取辅助登录账号
	GetDeputyAccounts(tc context.Context, userId int64, userfrom ...int32) ([]string, error)

	// IsNicknameReg 昵称是否已注册
	IsNicknameReg(tc context.Context, nickname string, ext ...map[string][]byte) (bool, error)

	GetUserAPropEx(tc context.Context, uid int64, adPropsEx []int32, ext ...map[string][]byte) (map[int32][]pbUA.AdvancedPropEx, error)

	GetAllEx(ct context.Context,
		uid int64,
		attrs []string,
		props []int32,
		adProps []int32,
		adPropExs []int32,
		gameName string,
		gameAttrs []string,
		exp ...map[string][]byte) (map[string]string, map[int32]int64, map[int32]string, map[int32][]pbUA.AdvancedPropEx, string, error)

	SetUserAPropEx(tc context.Context, uid int64, appid, channelid int32, adPropsEx map[int32][]pbUA.AdvancedPropEx, ext ...map[string][]byte) error

	SetAllEx(tc context.Context,
		uid int64,
		appid int32,
		channelid int32,
		attrs map[string]string,
		props map[int32]int64,
		adProps map[int32]string,
		adPropExs map[int32][]pbUA.AdvancedPropEx,
		ext ...map[string][]byte) error

	/**
	** 短信服务接口
	 */
	// SendCaptcha 发送验证码
	// uid 用户id
	// purpose 意图
	// clientIP 客户端ip（long）
	// ext 扩展字段（预留）
	// return:interval 每次验证码发送的间隔
	// return:surplus 距下次发送剩余的秒数
	// return:token 唯一标识串
	SendCaptcha(tc context.Context, uid int64, phone, purpose string, clientIp int64, ext ...map[string]string) (int32, int32, string, error)

	// VerifyCaptcha 校验验证码
	// token : 验证码标识
	// code : 验证码
	VerifyCaptcha(tc context.Context, uid int64, phone, token, content, purpose string, ext ...map[string]string) error

	/**
	** 邮件服务接口
	 */

	// SendMail 发送邮件
	// uid 用户id
	// title 标题
	// body 内容
	// sender 发送人(默认传 微乐捕鱼官方运营团队)
	// award 奖品
	// mailType 邮件类型
	SendMail(tc context.Context, uids []int64, title, body, sender string, award map[int64]int64, mailType int) error

	// SendMailV2 发送邮件
	// uid 用户id
	// title 标题
	// body 内容
	// sender 发送人(默认传 微乐捕鱼官方运营团队)
	// award 奖品
	// rxKeyNumber 瑞雪邮件标识
	// mailType 邮件类型
	SendMailV2(tc context.Context, uids []int64, title, body, sender, rxKeyNumber string, award map[int64]int64, mailType int) error

	// SendAllUserMail 发送全服邮件
	// title 标题
	// body 内容
	// sender 发送人(默认传 微乐捕鱼官方运营团队)
	// award 奖品
	// mailType 邮件类型
	SendAllUserMail(tc context.Context, title, body, sender string, award map[int64]int64, mailType int) error

	// SendBackMail 后台发送邮件
	SendBackMail(tc context.Context, content []byte) error

	// MailAwardList 邮件奖励列表
	// uid 用户id
	// mailid 邮件id  0代表全部邮件
	MailAwardList(tc context.Context, uid, mailid int64) ([]int64, map[int32]int64, error)

	// MailAllDetail 邮件奖励列表详情
	// uid 用户id
	// mailid 邮件id  0代表全部邮件
	MailAllDetail(tc context.Context, uid, mailid int64) ([]int64, map[int64]jsonmodel.MailDetailData, error)

	// MailEvaluation 邮件服务评价
	// 用户id
	// 邮件id
	// 是否满意（0-不满意，1-满意）
	MailEvaluation(tc context.Context, uid int64, mailid int64, pleased int) error

	// MailUpdateStatus 更新邮件状态
	// uid 用户id
	// mailids 邮件id
	MailUpdateStatus(tc context.Context, uid int64, mailids []int64) error

	// MailBatchRead 邮件一键已读,对应邮件置已读状态
	// uid 用户id
	// mailids 邮件id
	// ids 一键已读生效的邮件id
	MailBatchRead(tc context.Context, uid int64, mailids []int64) (ids []int64, err error)

	// MailBatchAward 邮件一键领取,对应邮件置已读&已领取
	// uid 用户id
	// mailids 邮件id
	// ids 一键领取生效的邮件id
	// content 邮件类型(2赠送邮件) body 内容
	MailBatchAward(tc context.Context, uid int64, mailids []int64) (ids []int64, award map[int32]int64, content []string, err error)

	// MailBatchAward 邮件一键删除,对应邮件置已删除
	// uid 用户id
	// mailids 邮件id
	// ids 一键删除生效的邮件id
	MailBatchDel(tc context.Context, uid int64, mailids []int64) (ids []int64, err error)

	/**
	** 游戏服务接口
	 */
	// KickOff 踢用户下线
	// uid : 用户id
	// gid : 游戏服务id
	// comment : 备注
	KickOff(tc context.Context, uid int64, gid uint32, comment string, typ ...uint16) error

	// GiveUp 玩家放弃当前比赛
	// uid : 用户id
	// gid : 游戏服务id
	// comment : 备注
	GiveUp(tc context.Context, uid int64, gid uint32, comment string) error

	/**
	** 排名服务接口
	 */
	// PushMatchScore 上传比赛积分
	// matchType 比赛类型
	//    1  娜迦千倍（周清）
	//    2  娜迦至尊(周清)
	//    3  鱼券赛(小时清)
	//    4  水晶赛(小时清)
	//    5  弹头赛(小时清)
	//    6  炸弹乐园(周清)
	PushMatchScore(tc context.Context, playID int64, matchType int, score int, expire int) error

	/**
	** 排行榜服务接口
	 */
	// GetRankList 获取排行榜
	// raceType 比赛类型
	//		1:  娜迦千倍(周清)
	//      2:  娜迦至尊(周清)
	//		3:  鱼券赛(小时清)
	//		4:  水晶赛(小时清)
	//		5:  弹头赛(小时清)
	//		6:  炸弹乐园(周清)
	//		7:  王者榜(周清)
	//		8:  富豪榜(周清)
	//		9:  弹头榜(周清)

	// []jsonmodel.RankData 排行列表信息
	// jsonmodel.MyRankData 玩家排行
	// string 上周排行key
	GetRankList(tc context.Context, raceType, front, isView int, userId int64) ([]jsonmodel.RankData, jsonmodel.MyRankData, string, error)

	/**
	** 排行榜服务接口
	 */
	// GetRankListByKey 获取排行榜列表
	// Key 存在redis key
	// size 获取前size
	// []jsonmodel.RankData 排行列表信息
	GetRankListByKey(tc context.Context, key string, size int) ([]jsonmodel.RankData, error)

	/**
	** 调度服务接口
	 */
	// GetRoomCardIDs批量申请房卡号
	// roomid 房间id
	// serverid 服务id（ip地址转int64）
	// num 申请的数量
	// return 房卡号集合
	GetRoomCardIDs(tc context.Context, roomid, serverid int64, num int) ([]string, error)

	// ReleaseRoomCardIDs 批量释放房卡号
	// roomid 房间id
	// serverid 服务id（ip地址转int64）
	// ids 释放的房卡号
	ReleaseRoomCardIDs(tc context.Context, roomid, serverid int64, ids []string) error

	// PayOrderCheck 支付下单校验
	// content 客户端请求参数json字节流
	PayOrderCheck(tc context.Context, content []byte, svrType uint16, svrid uint32) (*jsonmodel.RespPay, error)

	// PayCallBack 支付回调通知大厅
	// orderid  大厅订单ID
	// transactionid 平台订单ID
	PayCallBack(tc context.Context, orderid, transactionid string, svrType uint16, svrid uint32) (*jsonmodel.RespPay, error)

	// SendKFKMessage 异步发送kafka消息
	// topic 发送的主题
	// msg 发送的消息
	// key 取key[0]进行哈希，确定消息的分区
	SendKFKMessage(topic string, msg []byte, key ...string) error

	// 获得排行榜信息(通用)
	// gameName 	string	游戏名
	// raceType 	string	排行类型
	// condition int64	前多少名	<=0不限制
	GetCommonRankList(c context.Context, gameName string, raceType string, condition int64) (interface{}, error)

	// 设置玩家排名(通用)
	// args *jsonmodel.ReqSetCommonRank 需要用到的参数
	SetCommonUserRank(c context.Context, args *jsonmodel.ReqSetCommonRank) error

	// 清除排行榜信息(通用)
	// gameName 	string	游戏名
	// raceType 	string	排行类型
	ClearCommonRankList(c context.Context, gameName string, raceType string) error

	// 获得指定排行榜信息(通用)
	// gameName 	string		游戏名
	// raceType 	string		排行类型
	// userID 	string		玩家自己ID
	// list 		[][]int32	指定排名
	GetSpecifiedRankList(c context.Context, gameName string,
		raceType string, userID string, list [][]int32) (res *jsonmodel.SpecifiedRankData, err error)
}

type SdkClient struct {
	so *pCore.Socket
}

func Client() ISdk {
	so := gwList.Roll()
	return &SdkClient{so}
}

// GetUserAttr 获取用户基础属性
// uid 用户id
// fields 用户属性集合
func (c *SdkClient) GetUserAttr(tc context.Context, uid int64, attrKeys []string, ext ...map[string][]byte) (map[string]string, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, fmt.Errorf("param[uid] err")
	}
	if attrKeys == nil || len(attrKeys) <= 0 {
		return nil, fmt.Errorf("param[attrKeys] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_INFO)
	mData := &pbUA.ReqFields{}

	mData.Id = uid
	mData.Fields = attrKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)

	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &pbUA.RspFields{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.GetInfo(), nil
}

// GetUsersAttr 批量获取用户基础属性
// uids 用户id列表
// fields 用户属性集合
func (c *SdkClient) GetUsersAttr(tc context.Context, id int64, uids []int64, attrKeys []string, ext ...map[string][]byte) (map[int64]map[string]string, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if len(uids) <= 0 {
		return nil, fmt.Errorf("param[uids] err")
	}
	for _, id := range uids {
		if id <= 0 {
			return nil, fmt.Errorf("param[uids] err")
		}
	}
	if attrKeys == nil || len(attrKeys) <= 0 {
		return nil, fmt.Errorf("param[attrKeys] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_USERS_INFO)
	mData := &pbUA.ReqUsersFields{}

	mData.Id = id
	mData.Ids = uids
	mData.Fields = attrKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)

	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &pbUA.RspUsersFields{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	rspData := make(map[int64]map[string]string)
	for id, data := range mRespData.GetUsersInfo() {
		rspData[id] = data.GetInfo()
	}
	return rspData, nil
}

// GetUsersAttrAndGameAttr 批量获取用户基础属性&游戏属性
// uids 用户id列表
// fields 用户属性集合
func (c *SdkClient) GetUsersAttrAndGameAttr(tc context.Context, id int64, uids []int64,
	attrKeys, gameAttrKeys []string, gamename string, ext ...map[string][]byte) (map[int64]map[string]string, map[int64]string, error) {
	if c.so == nil {
		return nil, nil, errors.New("no available gateway online")
	}
	if len(uids) <= 0 {
		return nil, nil, fmt.Errorf("param[uids] err")
	}
	for _, id := range uids {
		if id <= 0 {
			return nil, nil, fmt.Errorf("param[uids] err")
		}
	}
	if len(attrKeys) <= 0 && len(gameAttrKeys) <= 0 {
		return nil, nil, fmt.Errorf("param[attrKeys,gameAttrKeys] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_USERS_INFO_EX)
	mData := &pbUA.ReqUsersFieldsEx{}

	mData.Id = id
	mData.Ids = uids
	mData.GameName = gamename
	mData.Fields = attrKeys
	mData.GameFields = gameAttrKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)

	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &pbUA.RspUsersFieldsEx{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	rspData := make(map[int64]map[string]string)
	for id, data := range mRespData.GetUsersInfo() {
		rspData[id] = data.GetInfo()
	}
	return rspData, mRespData.GetGameInfo(), nil
}

// SetUserAttr 设置用户属性
// uid 用户id
// appid 应用id
// channelid 渠道id
// attr 用户属性集合
func (c *SdkClient) SetUserAttr(tc context.Context, uid int64, appid, channelid int32, attr map[string]string, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}
	if attr == nil || len(attr) <= 0 {
		return fmt.Errorf("param[attr] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_SET_INFO)
	mData := &pbUA.ReqSetFields{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Info = attr
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &pbUA.RspSetFields{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}

	return nil
}

// GetUserProp 获取用户基础道具
func (c *SdkClient) GetUserProp(tc context.Context, uid int64, propKeys []int32, ext ...map[string][]byte) (map[int32]int64, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, fmt.Errorf("param[uid] err")
	}
	if propKeys == nil || len(propKeys) <= 0 {
		return nil, fmt.Errorf("param[propKeys] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_PROP)
	mData := &pbUA.ReqProp{}
	mData.Id = uid
	mData.Prop = propKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &pbUA.RspProp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.GetProp(), nil

}

// GetBatchUserInfoAndProp 批量获取用户信息和基础道具
func (c *SdkClient) GetBatchUserInfoAndProp(tc context.Context, uid int64, uids []int64,
	attrKeys []string, propKeys []int32, ext ...map[string][]byte) (map[int64]map[string]string, map[int64]map[int32]int64, error) {
	if c.so == nil {
		return nil, nil, errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, nil, fmt.Errorf("param[uid] err")
	}

	if len(uids) <= 0 {
		return nil, nil, fmt.Errorf("param[uids] err")
	}
	for _, id := range uids {
		if id <= 0 {
			return nil, nil, fmt.Errorf("param[uids] err")
		}
	}

	if (propKeys == nil || len(propKeys) <= 0) && (attrKeys == nil || len(attrKeys) <= 0) {
		return nil, nil, fmt.Errorf("param[attrKeys || propKeys] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_BATCH_INFO_AND_PROP)
	mData := &pbUA.ReqBatchInfoAndProp{}
	mData.Id = uid
	mData.Ids = uids
	mData.Fields = attrKeys
	mData.Prop = propKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &pbUA.RspBatchInfoAndProp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	rspPropData := make(map[int64]map[int32]int64)
	for id, data := range mRespData.GetUsersProp() {
		rspPropData[id] = data.GetProp()
	}

	rspInfoData := make(map[int64]map[string]string)
	for id, data := range mRespData.GetUsersInfo() {
		rspInfoData[id] = data.GetInfo()
	}

	return rspInfoData, rspPropData, nil

}

// GetBatchUserInfosAndProp 批量获取用户信息和基础道具和游戏属性信息
func (c *SdkClient) GetBatchUserInfosAndProp(tc context.Context, uid int64, uids []int64,
	attrKeys []string, propKeys []int32, gameAttrKeys []string, gamename string, ext ...map[string][]byte) (map[int64]map[string]string, map[int64]map[int32]int64, map[int64]string, error) {
	if c.so == nil {
		return nil, nil, nil, errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, nil, nil, fmt.Errorf("param[uid] err")
	}

	if len(uids) <= 0 {
		return nil, nil, nil, fmt.Errorf("param[uids] err")
	}
	for _, id := range uids {
		if id <= 0 {
			return nil, nil, nil, fmt.Errorf("param[uids] err")
		}
	}

	if (propKeys == nil || len(propKeys) <= 0) && (attrKeys == nil || len(attrKeys) <= 0) && (gameAttrKeys == nil || len(gameAttrKeys) <= 0) {
		return nil, nil, nil, fmt.Errorf("param[attrKeys || propKeys] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_BATCH_INFOS_AND_PROP)
	mData := &pbUA.ReqBatchInfosAndProp{}
	mData.Id = uid
	mData.Ids = uids
	mData.GameName = gamename
	mData.Fields = attrKeys
	mData.Prop = propKeys
	mData.GameFields = gameAttrKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &pbUA.RspBatchInfosAndProp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	rspPropData := make(map[int64]map[int32]int64)
	for id, data := range mRespData.GetUsersProp() {
		rspPropData[id] = data.GetProp()
	}

	rspInfoData := make(map[int64]map[string]string)
	for id, data := range mRespData.GetUsersInfo() {
		rspInfoData[id] = data.GetInfo()
	}

	return rspInfoData, rspPropData, mRespData.GameInfo, nil

}

// SetUserProp 设置用户道具信息
func (c *SdkClient) SetUserProp(tc context.Context, uid int64, appid, channelid int32, prop map[int32]int64, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}
	if prop == nil || len(prop) <= 0 {
		return fmt.Errorf("param[prop] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_SET_PROP)
	mData := &pbUA.ReqSetProp{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Prop = prop
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)

	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &pbUA.RspSetProp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil

}

// GetUserInfo 获取用户信息（包括属性和道具）
func (c *SdkClient) GetUserInfo(tc context.Context, uid int64, attrKeys []string, propKeys []int32, advancePropKeys []int32,
	gamename string, gameAttrKeys []string, ext ...map[string][]byte) (map[string]string, map[int32]int64, map[int32]string, string, error) {
	if c.so == nil {
		return nil, nil, nil, "", errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, nil, nil, "", fmt.Errorf("param[uid] err")
	}

	if (attrKeys == nil || len(attrKeys) <= 0) &&
		(propKeys == nil || len(propKeys) <= 0) &&
		(advancePropKeys == nil || len(advancePropKeys) <= 0) &&
		(gamename == "" || (gameAttrKeys == nil || len(gameAttrKeys) <= 0)) {
		return nil, nil, nil, "", fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_ALL)
	mData := &pbUA.ReqAll{}
	mData.Id = uid
	mData.Fields = attrKeys
	mData.Prop = propKeys
	mData.AdvancedProp = advancePropKeys
	mData.GameName = gamename
	mData.GameFields = gameAttrKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, nil, "", fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspAll{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, nil, "", fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, nil, "", fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.GetInfo(), mRespData.GetProp(), mRespData.GetAdvancedProp(), mRespData.GetGameInfo(), nil

}

// SetUserInfo 设置用户信息
func (c *SdkClient) SetUserInfo(tc context.Context, uid int64, appid, channelid int32, attr map[string]string, prop map[int32]int64, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}

	if (attr == nil || len(attr) <= 0) &&
		(prop == nil || len(prop) <= 0) {
		return fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_SET_ALL)
	mData := &pbUA.ReqSetAll{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Info = attr
	mData.Prop = prop
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspSetAll{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil

}

// OpProp 道具操作
// opType 操作类型，如充值、看广告等
// option 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
func (c *SdkClient) OpProp(tc context.Context, uid int64, appid, channelid int32, opType string, option int32,
	prop map[int32]int64, ext ...map[string][]byte) (map[int32]int64, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}

	if uid <= 0 {
		return nil, fmt.Errorf("param[uid] err")
	}

	if prop == nil || len(prop) <= 0 {
		return nil, fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_OP_PROP)
	mData := &pbUA.ReqOpProp{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Optype = opType
	mData.Option = option
	mData.Prop = prop
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspOpProp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.GetProp(), nil
}

// SetAll 设置用户属性、普通道具、高级道具
func (c *SdkClient) SetAll(tc context.Context, uid int64, appid, channelid int32, attr map[string]string, prop map[int32]int64, adProp map[int32]string, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}
	if attr == nil && prop == nil && adProp == nil {
		return nil
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_SET_ALL)
	mData := &pbUA.ReqSetAll{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Info = attr
	mData.Prop = prop
	mData.AdProp = adProp
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}

	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspSetAll{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

// GetUserGameInfo 获取用户游戏属性
func (c *SdkClient) GetUserGameInfo(tc context.Context, uid int64, gamename string, gameAttrKeys []string, ext ...map[string][]byte) (string, error) {
	if c.so == nil {
		return "", errors.New("no available gateway online")
	}
	if uid <= 0 {
		return "", fmt.Errorf("param[uid] err")
	}

	if gamename == "" || gameAttrKeys == nil || len(gameAttrKeys) <= 0 {
		return "", fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_GAME_INFO)
	mData := &pbUA.ReqGameInfo{}
	mData.Id = uid
	mData.Name = gamename
	mData.Fields = gameAttrKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return "", fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspGameInfo{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return "", fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return "", fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.GetInfo(), nil

}

// SetUserGameInfo 设置用户游戏属性
func (c *SdkClient) SetUserGameInfo(tc context.Context, uid int64, appid, channelid int32, gamename string, data string, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}

	if gamename == "" || data == "" {
		return fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_SET_GAME_INFO)
	mData := &pbUA.ReqSetGameInfo{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Name = gamename
	mData.Info = data
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspSetGameInfo{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

// OpUserGameInfo 操作用户游戏属性
// opType 操作类型，如充值、看广告等
// option 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
func (c *SdkClient) OpUserGameInfo(tc context.Context, uid int64, appid, channelid int32, opType string, option int32,
	gamename string, data map[string]int64, ext ...map[string][]byte) (map[string]int64, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, fmt.Errorf("param[uid] err")
	}

	if gamename == "" || data == nil {
		return nil, fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_OP_GAME_INFO)
	mData := &pbUA.ReqOpGameInfo{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Optype = opType
	mData.Option = option
	mData.Name = gamename
	mData.Info = data
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspOpGameInfo{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.GetInfo(), nil
}

// GetUserAProp 获取用户高级道具
func (c *SdkClient) GetUserAProp(tc context.Context, uid int64, aPropKeys []int32, ext ...map[string][]byte) (map[int32]string, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, fmt.Errorf("param[uid] err")
	}

	if aPropKeys == nil || len(aPropKeys) == 0 {
		return nil, fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_ADVANCE_PROP)
	mData := &pbUA.ReqAdProp{}
	mData.Id = uid
	mData.AdProp = aPropKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspAdProp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.GetAdProp(), nil
}

// SetUserAProp 操作用户高级
// opType 操作类型，如充值、看广告等
// option 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
func (c *SdkClient) SetUserAProp(tc context.Context, uid int64, appid, channelid int32, aProp map[int32]string, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}

	if aProp == nil || len(aProp) == 0 {
		return fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_SET_ADVANCE_PROP)
	mData := &pbUA.ReqSetAdProp{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.AdProp = aProp
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspSetAdProp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

// OpUserAProp 操作用户高级
// opType 操作类型，如充值、看广告等
// option 操作选项(位运算) 1: 要求必须登陆, 2: 允许在游戏中进行扣除操作
func (c *SdkClient) OpUserAProp(tc context.Context, uid int64, appid, channelid int32, opType string, option int32,
	aProp map[int32]string, ext ...map[string][]byte) (map[int32]string, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, fmt.Errorf("param[uid] err")
	}

	if aProp == nil || len(aProp) == 0 {
		return nil, fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_OP_ADVANCE_PROP)
	mData := &pbUA.ReqOpAdProp{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Optype = opType
	mData.Option = option
	mData.AdProp = aProp
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspOpAdProp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.GetAdProp(), nil
}

// OpPropAndAdProp 原子操作普通道具与高级道具
func (c *SdkClient) OpPropAndAdProp(tc context.Context, uid int64, appid, channelid int32, opType string, option int32,
	prop map[int32]int64, adProp map[int32]string, ext ...map[string][]byte) (map[int32]int64, map[int32]string, error) {
	if c.so == nil {
		return nil, nil, errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, nil, fmt.Errorf("param[uid] err")
	}

	if prop == nil || adProp == nil {
		return nil, nil, fmt.Errorf("param err")
	}
	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_OP_PROP_AND_ADVANCED_PROP)
	mData := &pbUA.ReqOpPropAndAdProp{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Optype = opType
	mData.Option = option
	mData.Prop = prop
	mData.AdProp = adProp
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}

	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspOpPropAndAdProp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}

	return mRespData.GetProp(), mRespData.GetAdProp(), nil
}

// EnterGame 进入游戏
func (c *SdkClient) EnterGame(tc context.Context, uid int64, gameid string, gamename, roomname string,
	attrKeys []string, propKeys []int32, aPropKeys []int32,
	ext ...map[string][]byte) (map[string]string, map[int32]int64, map[int32]string, string, string, error) {
	if c.so == nil {
		return nil, nil, nil, "", "", errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, nil, nil, "", "", fmt.Errorf("param[uid] err")
	}

	if gameid == "" {
		return nil, nil, nil, "", "", fmt.Errorf("param[gameid] err")
	}
	if gamename == "" {
		return nil, nil, nil, "", "", fmt.Errorf("param[gamename] err")
	}

	if (attrKeys == nil || len(attrKeys) == 0) &&
		(propKeys == nil || len(propKeys) == 0) &&
		(aPropKeys == nil || len(aPropKeys) == 0) {
		return nil, nil, nil, "", "", fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_ENTER_GAME)
	mData := &pbUA.ReqEnterGame{}
	mData.Id = uid
	mData.RoomId = roomname
	mData.Gid = gameid
	mData.GameName = gamename
	mData.Fields = attrKeys
	mData.Prop = propKeys
	mData.AdvancedProp = aPropKeys
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, nil, "", "", fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspEnterGame{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, nil, "", "", fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, nil, "", mRespData.RoomId, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.GetInfo(), mRespData.GetProp(), mRespData.GetAdvancedProp(), mRespData.GetGameInfo(), "", nil
}

// LeaveGame 离开游戏
func (c *SdkClient) LeaveGame(tc context.Context, uid int64, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_LEAVE_GAME)
	mData := &pbUA.ReqLeaveGame{}
	mData.Id = uid
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspLeaveGame{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

// IsUsernameReg 检查用户名是否已注册
func (c *SdkClient) IsUsernameReg(tc context.Context, username string, ext ...map[string][]byte) (bool, error) {
	if c.so == nil {
		return false, errors.New("no available gateway online")
	}
	if username == "" {
		return false, fmt.Errorf("param[username] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_CHECK_USERNAME_EXISTED)
	mData := &pbUA.ReqUsernameExist{}
	mData.Username = username
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return false, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspUsernameExist{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return false, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return false, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.Existed, nil
}

// SetPassword 重置密码
func (c *SdkClient) SetPassword(tc context.Context, uid int64, pwd string, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}
	if pwd == "" {
		return fmt.Errorf("param[pwd] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_SET_PASS)
	mData := &pbUA.ReqSetPass{}
	mData.Id = uid
	mData.Pass = pwd
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspSetPass{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

func (c *SdkClient) UpdatePassword(tc context.Context, uid int64, oldPwd, newPwd string, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}
	if oldPwd == "" {
		return fmt.Errorf("param[oldPwd] err")
	}
	if newPwd == "" {
		return fmt.Errorf("param[newPwd] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_UPDATE_PASS)
	mData := &pbUA.ReqUpdatePass{}
	mData.Id = uid
	mData.OldPass = oldPwd
	mData.NewPass = newPwd
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspUpdatePass{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

// HasBadWord 敏感词检测
// return 敏感词串，空值表示无敏感词
func (c *SdkClient) HasBadWord(tc context.Context, str string, ext ...map[string][]byte) ([]string, error) {
	if str == "" {
		return nil, nil
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_HAS_BAD_WROD)
	mData := &pbUA.ReqHasBadWord{}
	mData.Content = str
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspHasBadWord{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.Words, nil
}

// 敏感词替换
// return 替换后的新串
func (c *SdkClient) ReplaceBadWord(tc context.Context, str string, ext ...map[string][]byte) (string, error) {

	if str == "" {
		return "", nil
	}
	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_REPLACE_BAD_WORD)
	mData := &pbUA.ReqReplaceBadWord{}
	mData.Content = str
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return "", fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspReplaceBadWord{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return "", fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return "", fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.Content, nil
}

// CheckIdCardExisted 检查用户idcard是否存在
func (c *SdkClient) CheckIdCardExisted(tc context.Context, idcard string, ext ...map[string][]byte) (bool, error) {
	if c.so == nil {
		return false, errors.New("no available gateway online")
	}
	if idcard == "" {
		return false, fmt.Errorf("param[idcard] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_CHECK_IDCARD_EXISTED)
	mData := &pbUA.ReqCheckIdCardExisted{}
	mData.IdCard = idcard
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return false, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspCheckIdCardExisted{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return false, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return false, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.Existed, nil
}

// GetUsersByPhoneNumber 通过电话号码获取用户列表
func (c *SdkClient) GetUsersByPhoneNumber(tc context.Context, phone string, ext ...map[string][]byte) ([]*pbUA.RspGetUsersByPhoneNumberUser, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if phone == "" {
		return nil, fmt.Errorf("param[phone] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_USERS_BY_PHONE_NUMBER)
	mData := &pbUA.ReqGetUsersByPhoneNumber{}
	mData.Phone = phone
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspGetUsersByPhoneNumber{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.Users, nil
}

// AddDeputyAccount 添加辅助登录账号
func (c *SdkClient) AddDeputyAccount(tc context.Context, userId int64, username, password string, userfrom ...int32) error {
	if c.so == nil {
		return fmt.Errorf("no available gateway online")
	}

	if userId == 0 || username == "" || password == "" {
		return fmt.Errorf("param err")
	}

	var from int32 = 1
	if len(userfrom) > 0 {
		from = userfrom[0]
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_ADD_DEPUTY_ACCOUT)
	mData := &pbUA.ReqAddDeputyAccount{}
	mData.Id = userId
	mData.Username = username
	mData.Userfrom = from
	mData.Password = password

	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspAddDeputyAccount{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}

	return nil
}

// DelDeputyAccount 删除指定用户的某个辅助登录账号
func (c *SdkClient) DelDeputyAccount(tc context.Context, userId int64, username string, userfrom ...int32) (bool, error) {
	if c.so == nil {
		return false, fmt.Errorf("no available gateway online")
	}

	if userId == 0 || username == "" {
		return false, fmt.Errorf("param err")
	}

	var from int32 = 1
	if len(userfrom) > 0 {
		from = userfrom[0]
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_DEL_DEPUTY_ACCOUT)
	mData := &pbUA.ReqDelDeputyAccount{}
	mData.Id = userId
	mData.Username = username
	mData.Userfrom = from

	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return false, fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspDelDeputyAccount{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return false, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return false, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}

	return mRespData.Success, nil
}

// GetDeputyAccounts 获取用户的辅助登录账号
func (c *SdkClient) GetDeputyAccounts(tc context.Context, userId int64, userfrom ...int32) ([]string, error) {
	if c.so == nil {
		return nil, fmt.Errorf("no available gateway online")
	}

	if userId == 0 {
		return nil, fmt.Errorf("param err")
	}

	var from int32 = 1
	if len(userfrom) > 0 {
		from = userfrom[0]
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_DEPUTY_ACCOUT)
	mData := &pbUA.ReqGetDeputyAccount{}
	mData.Id = userId
	mData.Userfrom = from

	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspGetDeputyAccount{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}

	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}

	return mRespData.Username, nil
}

// IsNicknameReg 检查昵称是否已注册
func (c *SdkClient) IsNicknameReg(tc context.Context, nickname string, ext ...map[string][]byte) (bool, error) {
	if c.so == nil {
		return false, errors.New("no available gateway online")
	}
	if nickname == "" {
		return false, fmt.Errorf("param[username] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_CHECK_NICKNAME_EXISTED)
	mData := &pbUA.ReqNicknameExist{}
	mData.Nickname = nickname
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return false, fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspNicknameExist{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return false, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return false, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.Existed, nil
}

// GetUserAPropEx 获取用户扩展高级道具
func (c *SdkClient) GetUserAPropEx(tc context.Context, uid int64, adPropsEx []int32, ext ...map[string][]byte) (map[int32][]pbUA.AdvancedPropEx, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, fmt.Errorf("param[uid] err")
	}

	if len(adPropsEx) == 0 {
		return nil, fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_ADVANCED_PROP_EX)
	mData := &pbUA.ReqAdPropsEx{}
	mData.Id = uid
	mData.AdPropsEx = adPropsEx
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspAdPropsEx{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}

	apxMap := make(map[int32][]pbUA.AdvancedPropEx, len(mRespData.AdPropsEx))
	for k, v := range mRespData.AdPropsEx {
		var apxArr []pbUA.AdvancedPropEx
		if err := json.Unmarshal([]byte(v), &apxArr); err != nil {
			continue
		}
		apxMap[k] = apxArr
	}

	return apxMap, nil
}

// GetAllEx 获取用户属性、普通道具、高级道具、扩展高级道具、游戏属性
func (c *SdkClient) GetAllEx(tc context.Context,
	uid int64,
	attrs []string,
	props []int32,
	adProps []int32,
	adPropExs []int32,
	gameName string,
	gameAttrs []string,
	ext ...map[string][]byte) (map[string]string, map[int32]int64, map[int32]string, map[int32][]pbUA.AdvancedPropEx, string, error) {
	if c.so == nil {
		return nil, nil, nil, nil, "", errors.New("no available gateway online")
	}
	if uid <= 0 {
		return nil, nil, nil, nil, "", fmt.Errorf("param[uid] err")
	}

	if len(attrs) == 0 && len(props) == 0 && len(adProps) == 0 && len(adPropExs) == 0 && gameName == "" && len(gameAttrs) == 0 {
		return nil, nil, nil, nil, "", fmt.Errorf("param err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_GET_ALL_EX)
	mData := &pbUA.ReqAllEx{}
	mData.Id = uid
	mData.Fields = attrs
	mData.Props = props
	mData.AdProps = adProps
	mData.AdPropsEx = adPropExs
	mData.GameName = gameName
	mData.GameFields = gameAttrs
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, nil, nil, "", fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspAllEx{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, nil, nil, "", fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, nil, nil, "", fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}

	apxMap := make(map[int32][]pbUA.AdvancedPropEx, len(mRespData.AdProps))
	for k, v := range mRespData.AdPropsEx {
		var apxArr []pbUA.AdvancedPropEx
		if err := json.Unmarshal([]byte(v), &apxArr); err != nil {
			continue
		}
		apxMap[k] = apxArr
	}
	return mRespData.Info, mRespData.Props, mRespData.AdProps, apxMap, mRespData.GameInfo, nil
}

// SetUserAPropEx 设置用户扩展高级道具
func (c *SdkClient) SetUserAPropEx(tc context.Context, uid int64, appid, channelid int32, adPropsEx map[int32][]pbUA.AdvancedPropEx, ext ...map[string][]byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}
	if len(adPropsEx) == 0 {
		return fmt.Errorf("param err")
	}

	apxMap := make(map[int32]string, 0)
	for k, v := range adPropsEx {
		buf, _ := json.Marshal(v)
		apxMap[k] = string(buf)
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_SET_ADVANCED_PROP_EX)
	mData := &pbUA.ReqSetAdPropsEx{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.AdPropsEx = apxMap
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspSetAdPropsEx{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil

}

// SetAllEx 设置用户属性、普通道具、高级道具、扩展高级道具
func (c *SdkClient) SetAllEx(tc context.Context,
	uid int64,
	appid int32,
	channelid int32,
	attrs map[string]string,
	props map[int32]int64,
	adProps map[int32]string,
	adPropsEx map[int32][]pbUA.AdvancedPropEx,
	ext ...map[string][]byte) error {

	if c.so == nil {
		return errors.New("no available gateway online")
	}

	if uid <= 0 {
		return fmt.Errorf("param[uid] err")
	}

	if len(attrs) == 0 && len(props) == 0 && len(adProps) == 0 && len(adPropsEx) == 0 {
		return nil
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_USER_API)
	m.SetFunctionID(F_ID_SET_ALL_EX)
	mData := &pbUA.ReqSetAllEx{}
	mData.Id = uid
	mData.Appid = appid
	mData.Channelid = channelid
	mData.Info = attrs
	mData.Props = props
	mData.AdProps = adProps

	if len(adPropsEx) != 0 {
		apxMap := make(map[int32]string, 0)
		for k, v := range adPropsEx {
			buf, _ := json.Marshal(v)
			apxMap[k] = string(buf)

			mData.AdPropsEx = apxMap
		}
	}

	if len(ext) > 0 {
		mData.Ext = ext[0]
	}

	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &pbUA.RspSetAllEx{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

// ----------------------------------------------------------------------------------
// -- 短信接口 start
// ----------------------------------------------------------------------------------
// SendCaptcha 发送验证码
// return
// interval 每次验证码发送的间隔
// surplus 距下次发送剩余的秒数
// token 唯一标识串
func (c *SdkClient) SendCaptcha(tc context.Context, uid int64, phone, purpose string, clientIp int64, ext ...map[string]string) (int32, int32, string, error) {
	if c.so == nil {
		return 0, 0, "", fmt.Errorf("no available gateway online")
	}
	// if uid <= 0 {
	// 	return 0, 0, "", fmt.Errorf("param[uid] err")
	// }
	if phone == "" {
		return 0, 0, "", fmt.Errorf("param[phone] err")
	}
	if purpose == "" {
		return 0, 0, "", fmt.Errorf("param[purpose] err")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_SMS)
	m.SetFunctionID(F_ID_SMS_SEND)
	mData := &pbSMS.ReqSend{}
	mData.Userid = uid
	mData.Phone = phone
	mData.Purpose = purpose
	mData.Ip = clientIp
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}
	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return 0, 0, "", fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbSMS.RespSend{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return 0, 0, "", fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return 0, 0, "", fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return mRespData.Interval, mRespData.Surplus, mRespData.Token, nil
}

// VerifyCaptcha 校验验证码
// token : 验证码标识
// code : 验证码
// return : 短信服务响应消息
func (c *SdkClient) VerifyCaptcha(tc context.Context, uid int64, phone, token, content, purpose string, ext ...map[string]string) error {
	if c.so == nil {
		return fmt.Errorf("no available gateway online")
	}
	// if uid <= 0 {
	// 	return fmt.Errorf("param[uid] err")
	// }
	if phone == "" {
		return fmt.Errorf("param[phone] err")
	}
	if token == "" {
		return fmt.Errorf("param[token] err")
	}
	if content == "" {
		return fmt.Errorf("param[content] err")
	}
	if purpose == "" {
		return fmt.Errorf("param[purpose] err")
	}

	// 发送请求到短信服务
	m := pCore.NewMessage(pCore.MT_REQUEST)
	m.SetToSvrType(ST_SMS)
	m.SetFunctionID(F_ID_SMS_VERIFY)
	mData := &pbSMS.ReqVerify{}
	mData.Userid = uid
	mData.Phone = phone
	mData.Token = token
	mData.Content = content
	mData.Purpose = purpose
	if len(ext) > 0 {
		mData.Ext = ext[0]
	}

	mDataBuf, _ := proto.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		// logs.Errorf("fail to verifySMS,req : %+v,err:%v", req, err)
		return fmt.Errorf("fail to verifySMS, : err:%v", err)
	}
	// 接受短信服务响应
	mRespData := &pbSMS.RespVerify{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("fail to decode resp from sendSMS err:%v", err.Error())
	}

	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

// ----------------------------------------------------------------------------------
// -- 短信接口 end
// ----------------------------------------------------------------------------------

// ----------------------------------------------------------------------------------
// -- 邮件接口接口 start
// ----------------------------------------------------------------------------------
// SendMail 发送邮件
func (c *SdkClient) SendMail(tc context.Context, uids []int64, title, body, sender string, award map[int64]int64, mailType int) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if len(uids) <= 0 {
		return errors.New("user is nil")
	}
	if title == "" || body == "" {
		return errors.New("body or title is nil")
	}
	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_SEND)
	mData := &jsonmodel.MailMsg{}
	mData.Userid = uids
	mData.Title = title
	mData.Body = body
	mData.Sender = sender
	for k, v := range award {
		mData.Attach.Award = append(mData.Attach.Award, jsonmodel.Award{DataId: int32(k), DataValue: v})
	}
	mData.Type = mailType
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RespMail{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return errors.New(mRespData.Msg)
	}
	return nil
}

// SendBackMail 发送后台邮件
func (c *SdkClient) SendBackMail(tc context.Context, content []byte) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if len(content) <= 0 {
		return errors.New("content is nil")
	}
	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_BACK_SEND)
	m.SetBody(content)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RespMail{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return errors.New(mRespData.Msg)
	}
	return nil
}

// MailAwardList 邮件奖励
func (c *SdkClient) MailAwardList(tc context.Context, uid, mailid int64) ([]int64, map[int32]int64, error) {
	if c.so == nil {
		return nil, nil, errors.New("no available gateway online")
	}
	if uid == 0 {
		return nil, nil, errors.New("user is nil")
	}
	if mailid < 0 {
		return nil, nil, errors.New("mailid is nil")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_AWARD)
	mData := &jsonmodel.MailAward{}
	mData.Userid = uid
	mData.Mailid = mailid
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RspMailAward{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, errors.New(mRespData.Msg)
	}
	return mRespData.MailIds, mRespData.Data, nil
}

// MailAllDetail 邮件详情
func (c *SdkClient) MailAllDetail(tc context.Context, uid, mailid int64) ([]int64, map[int64]jsonmodel.MailDetailData, error) {
	if c.so == nil {
		return nil, nil, errors.New("no available gateway online")
	}
	if uid == 0 {
		return nil, nil, errors.New("user is nil")
	}
	if mailid < 0 {
		return nil, nil, errors.New("mailid is nil")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_ALL_DETAIL)
	mData := &jsonmodel.MailAward{}
	mData.Userid = uid
	mData.Mailid = mailid
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RspMailAllDetail{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, errors.New(mRespData.Msg)
	}
	return mRespData.MailIds, mRespData.Data, nil
}

// 邮件评价
func (c *SdkClient) MailEvaluation(tc context.Context, uid int64, mailid int64, pleased int) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid == 0 {
		return errors.New("user is nil")
	}
	if mailid <= 0 {
		return errors.New("mailid is nil")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_EVALUATION)
	mData := &jsonmodel.MailEvaluation{}
	mData.Userid = uid
	mData.MailsId = mailid
	mData.Pleased = pleased
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RespCommon{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return errors.New(mRespData.Msg)
	}
	return nil
}

// MailUpdateStatus 更新邮件状态
func (c *SdkClient) MailUpdateStatus(tc context.Context, uid int64, mailids []int64) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if uid == 0 {
		return errors.New("user is nil")
	}
	if len(mailids) == 0 {
		return errors.New("mailids is nil")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_STATUS)
	mData := &jsonmodel.MailStatus{}
	mData.Userid = uid
	mData.MailsIds = mailids
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RspMailStatus{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return errors.New(mRespData.Msg)
	}
	return nil
}

// MailBatchRead 一键已读
func (c *SdkClient) MailBatchRead(tc context.Context, uid int64, mailids []int64) ([]int64, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if uid == 0 {
		return nil, errors.New("user is nil")
	}
	if len(mailids) == 0 {
		return nil, errors.New("mailids is nil")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_BATCH_READ)
	mData := &jsonmodel.MailBatch{}
	mData.Userid = uid
	mData.MailsIds = mailids
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RspMailBatchRead{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, errors.New(mRespData.Msg)
	}
	return mRespData.MailIds, nil
}

// MailBatchAward 一键领取
func (c *SdkClient) MailBatchAward(tc context.Context, uid int64, mailids []int64) ([]int64, map[int32]int64, []string, error) {
	if c.so == nil {
		return nil, nil, nil, errors.New("no available gateway online")
	}
	if uid == 0 {
		return nil, nil, nil, errors.New("user is nil")
	}
	if len(mailids) == 0 {
		return nil, nil, nil, errors.New("mailids is nil")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_BATCH_AWARD)
	mData := &jsonmodel.MailBatch{}
	mData.Userid = uid
	mData.MailsIds = mailids
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RspMailBatchAward{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, nil, nil, errors.New(mRespData.Msg)
	}
	return mRespData.MailIds, mRespData.Data, mRespData.Content, nil
}

// ----------------------------------------------------------------------------------
// -- 邮件接口接口 start
// ----------------------------------------------------------------------------------
// SendMail 发送邮件
func (c *SdkClient) SendMailV2(tc context.Context, uids []int64, title, body, sender, rxKeyNumber string, award map[int64]int64, mailType int) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if len(uids) <= 0 {
		return errors.New("user is nil")
	}
	if title == "" || body == "" {
		return errors.New("body or title is nil")
	}
	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_SEND)
	mData := &jsonmodel.MailMsg{}
	mData.Userid = uids
	mData.Title = title
	mData.Body = body
	mData.Sender = sender
	mData.RXKeyNumber = rxKeyNumber
	for k, v := range award {
		mData.Attach.Award = append(mData.Attach.Award, jsonmodel.Award{DataId: int32(k), DataValue: v})
	}
	mData.Type = mailType
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RespMail{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return errors.New(mRespData.Msg)
	}
	return nil
}

func (c *SdkClient) SendAllUserMail(tc context.Context, title, body, sender string, award map[int64]int64, mailType int) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	if title == "" || body == "" {
		return errors.New("body or title is nil")
	}
	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_SEND_ALL)
	mData := &jsonmodel.MailMsg{}
	mData.Title = title
	mData.Body = body
	mData.Sender = sender
	for k, v := range award {
		mData.Attach.Award = append(mData.Attach.Award, jsonmodel.Award{DataId: int32(k), DataValue: v})
	}
	mData.Type = mailType
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &jsonmodel.RespMail{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return errors.New(mRespData.Msg)
	}
	return nil
}

// MailBatchDel 一键删除
func (c *SdkClient) MailBatchDel(tc context.Context, uid int64, mailids []int64) (ids []int64, err error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if uid == 0 {
		return nil, errors.New("user is nil")
	}
	if len(mailids) == 0 {
		return nil, errors.New("mailids is nil")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_MAIL)
	m.SetFunctionID(F_ID_MAIL_BATCH_DEL)
	mData := &jsonmodel.MailBatch{}
	mData.Userid = uid
	mData.MailsIds = mailids
	mDataBuf, _ := json.Marshal(mData)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RspMailBatchDel{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, errors.New(mRespData.Msg)
	}
	return mRespData.MailIds, nil
}

// ----------------------------------------------------------------------------------
// -- 邮件接口 end
// ----------------------------------------------------------------------------------

// ----------------------------------------------------------------------------------
// -- 游戏接口 start
// ----------------------------------------------------------------------------------
// KickOff 踢用户下线
// uid : 用户id
// gid : 游戏服务id
// comment : 备注
func (c *SdkClient) KickOff(tc context.Context, uid int64, gid uint32, comment string, typ ...uint16) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}

	svrType := uint16(ST_GAME)
	if len(typ) != 0 {
		svrType = typ[0]
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(svrType)
	m.SetToSvrID(gid)
	m.SetFunctionID(F_ID_GAME_KICKOFF)
	mData := &pbUA.ReqKickOff{}
	mData.Id = uid
	mData.Msg = comment

	buf, _ := proto.Marshal(mData)
	m.SetBody(buf)

	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspKickOff{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

// ----------------------------------------------------------------------------------
// -- 游戏接口 start
// ----------------------------------------------------------------------------------
// GiveUp 玩家放弃当前比赛
// uid : 用户id
// gid : 游戏服务id
// comment : 备注
func (c *SdkClient) GiveUp(tc context.Context, uid int64, gid uint32, comment string) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_GAME)
	m.SetToSvrID(gid)
	m.SetFunctionID(F_ID_GAME_GIVEUP)
	mData := &pbUA.ReqGiveUp{}
	mData.Id = uid
	mData.Msg = comment

	buf, _ := proto.Marshal(mData)
	m.SetBody(buf)

	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &pbUA.RspGiveUp{}
	err = proto.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}
	return nil
}

// PushMatchScore 上传比赛积分
// matchType 比赛类型 1    娜迦千倍（周清）
//
//	2  娜迦至尊(周清)
//	3  鱼券赛(小时清)
//	4  水晶赛(小时清)
//	5  弹头赛(小时清)
//	6  炸弹乐园(周清)
func (c *SdkClient) PushMatchScore(tc context.Context, playID int64, matchType int, score int, expire int) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	req := pCore.NewNormalMessage()
	req.SetToSvrType(ST_RANK)
	req.SetFunctionID(F_ID_RANK_PUSH_SCORE)

	data := make(map[string]interface{})
	data["user_id"] = playID
	data["race_type"] = matchType
	data["score"] = score
	data["expire"] = expire
	body, _ := utils.EncodeJson(&data)

	req.SetBody(body)
	_, err := c.so.Send(tc, req)
	return err
}

// GetRankList 获取排行榜
// raceType 比赛类型
//
//			1:  娜迦千倍(周清)
//	     2:  娜迦至尊(周清)
//			3:  鱼券赛(小时清)
//			4:  水晶赛(小时清)
//			5:  弹头赛(小时清)
//			6:  炸弹乐园(周清)
//			7:  王者榜(周清)
//			8:  富豪榜(周清)
//			9:  弹头榜(周清)
func (c *SdkClient) GetRankList(tc context.Context, raceType, front, isView int, userId int64) (dataList []jsonmodel.RankData,
	user jsonmodel.MyRankData, lastKey string, err error) {
	if c.so == nil {
		return nil, user, "", errors.New("no available gateway online")
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_RANK)
	m.SetFunctionID(F_ID_RANK_LIST)

	data := make(map[string]interface{})
	data["race_type"] = raceType
	data["front"] = front
	data["is_view"] = isView
	data["user_id"] = userId

	body, _ := utils.EncodeJson(&data)

	m.SetBody(body)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, user, "", fmt.Errorf("send msg err:%v", err.Error())
	}

	//
	mRespData := &jsonmodel.RspRankList{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, user, "", fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, user, "", errors.New(mRespData.Msg)
	}
	return mRespData.RankList, mRespData.User, mRespData.LastKey, err
}

// GetRankListByKey 获取排行榜列表
func (c *SdkClient) GetRankListByKey(tc context.Context, key string, size int) ([]jsonmodel.RankData, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}

	req := pCore.NewRequestMessage()
	req.SetToSvrType(ST_RANK)
	req.SetFunctionID(F_ID_RANK_KEY_KEY)

	data := make(map[string]interface{})
	data["key"] = key
	data["size"] = size

	body, _ := utils.EncodeJson(&data)

	req.SetBody(body)
	mResp, err := c.so.Send(tc, req)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}

	mRespData := &jsonmodel.RspRankAward{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, errors.New(mRespData.Msg)
	}
	return mRespData.RankList, err
}

// GetRoomCardIDs 批量获取房卡号
// roomid 房间id
// serverid 服务id（ip地址转int64）
// num 请求个数（1次不能超过500）
func (c *SdkClient) GetRoomCardIDs(tc context.Context, roomid, serverid int64, num int) ([]string, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	req := pCore.NewRequestMessage()
	req.SetToSvrType(ST_GAMEDISPATCH)
	req.SetFunctionID(F_ID_GAMEDISPATCH_GET_ROOMCARDIDS)

	data := jsonmodel.ReqGetRoomCardIDs{}
	data.RoomID = roomid
	data.ServerID = serverid
	data.Num = num
	body, _ := json.Marshal(data)
	req.SetBody(body)
	mResp, err := c.so.Send(tc, req)

	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RespGetRoomCardIDs{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}

	if len(mRespData.IDs) == 0 {
		return nil, errors.New("not available roomcardids")
	}

	return mRespData.IDs, nil
}

// ReleaseRoomCardIDs 批量获取房卡号
// roomid 房间id
// serverid 服务id（ip地址转int64）
// ids 释放的房卡号
func (c *SdkClient) ReleaseRoomCardIDs(tc context.Context, roomid, serverid int64, ids []string) error {
	if c.so == nil {
		return errors.New("no available gateway online")
	}
	req := pCore.NewRequestMessage()
	req.SetToSvrType(ST_GAMEDISPATCH)
	req.SetFunctionID(F_ID_GAMEDISPATCH_RELEASE_ROOMCARDIDS)

	data := jsonmodel.ReqReleaseRoomCardIDs{}
	data.RoomID = roomid
	data.ServerID = serverid
	data.IDs = ids
	body, _ := json.Marshal(data)
	req.SetBody(body)
	mResp, err := c.so.Send(tc, req)

	if err != nil {
		return fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RespCommon{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return fmt.Errorf("%v: (%d)%s", errRemoteServer, mRespData.Code, mRespData.Msg)
	}

	return nil
}

// 下单接口 请求大厅判断是否能付
// content 客户端下单请求参数json 字节数组
func (c *SdkClient) PayOrderCheck(tc context.Context, content []byte, svrType uint16, svrid uint32) (*jsonmodel.RespPay, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if len(content) <= 0 {
		return nil, errors.New("content is nil")
	}
	m := pCore.NewRequestMessage()
	if svrid != 0 {
		m.SetToSvrID(svrid)
		m.SetFunctionID(F_ID_PAY_CHECK_GAME)
	} else {
		m.SetFunctionID(F_ID_PAY_CHECK)
	}
	m.SetToSvrType(svrType)
	m.SetBody(content)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	//
	mRespData := &jsonmodel.RespPay{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, errors.New(mRespData.Msg)
	}
	return mRespData, nil
}

// 回调接口
// orderid 大厅订单号
// transaction 平台订单号
func (c *SdkClient) PayCallBack(tc context.Context, orderid, transactionid string, svrType uint16, svrid uint32) (*jsonmodel.RespPay, error) {
	if c.so == nil {
		return nil, errors.New("no available gateway online")
	}
	if orderid == "" || transactionid == "" {
		return nil, errors.New("orderid or transactionid is nil")
	}

	m := pCore.NewRequestMessage()
	if svrid != 0 {
		m.SetToSvrID(svrid)
		m.SetFunctionID(F_ID_PAY_CALLBACK_GAME)
	} else {
		m.SetFunctionID(F_ID_PAY_CALLBACK)
	}
	m.SetToSvrType(svrType)
	data := &jsonmodel.CallBackPay{OrderId: orderid, TransactionID: transactionid}
	mDataBuf, _ := json.Marshal(data)
	m.SetBody(mDataBuf)
	mResp, err := c.so.Send(tc, m)
	if err != nil {
		return nil, fmt.Errorf("send msg err:%v", err.Error())
	}
	mRespData := &jsonmodel.RespPay{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil {
		return nil, fmt.Errorf("decode resp err:%v", err.Error())
	}
	if mRespData.Code != 0 {
		return nil, errors.New(mRespData.Msg)
	}
	return mRespData, nil
}

// SendKFKMessage 异步发送kafka消息
// topic 发送的主题
// msg 发送的消息
// key 取key[0]进行哈希，确定消息的分区
func (c *SdkClient) SendKFKMessage(topic string, msg []byte, key ...string) error {
	return srv.KafkaProducer().SendMessage(topic, msg, key...)
}

/*
GetCommonRankList
@Desc	获得通用排行榜列表信息
@Param	gameName 	string			游戏名
@Param	raceType 	string			排行类型
@Param	condition 	int64			前多少名	<=0不限制
@Return	res			interface{}		排名数据
*/
func (s *SdkClient) GetCommonRankList(c context.Context, gameName string, raceType string, condition int64) (res interface{}, err error) {
	if s.so == nil {
		// logs.Errorf("GetCommonRankList sdkClient is nil")
		return nil, errors.New("GetCommonRankList sdkClient is nil")
	}

	if gameName == "" || raceType == "" {
		str := fmt.Sprintf("some args is nil,gameName:%v,raceType:%v,condition:%v", gameName, raceType, condition)
		// logs.Errorf(str)
		return nil, errors.New(str)
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_RANK)
	m.SetFunctionID(F_ID_RANK_GET_LIST)

	// 准备数据
	args := &jsonmodel.ReqGetCommonRank{
		GameName:  gameName,
		RaceType:  raceType,
		Condition: condition,
	}
	mDataBuf, _ := json.Marshal(args)
	m.SetBody(mDataBuf)

	// 调用rank服获取排行信息
	mResp, err := s.so.Send(c, m)
	if err != nil {
		// logs.Errorf("get commonRankList sdkClient send to rank is failed,args:%+v,err:%v", args, err)
		return nil, err
	}

	// rsp解析
	mRespData := &jsonmodel.RspCommonRankData{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil || mRespData.Code != 0 {
		// logs.Errorf("get commonRankList Unmarshal rsp body is failed,args:%+v,code,err:%v", args, mRespData.Code, err)
		return nil, err
	}

	// logs.Infof("get commonRankList success,args:%+v,mRespData:%+v", args, mRespData)
	return mRespData.Data, nil
}

/*
SetCommonUserRank
@Desc	设置排行榜
@Param	args 	*jsonmodel.ReqSetCommonRank
@Return			error
*/
func (s *SdkClient) SetCommonUserRank(c context.Context, args *jsonmodel.ReqSetCommonRank) error {
	if s.so == nil {
		// logs.Errorf("SetCommonUserRank sdkClient is nil")
		return errors.New("SetCommonUserRank sdkClient is nil")
	}

	if args == nil {
		// logs.Errorf("SetCommonUserRank args is nil,args:%+v", args)
		return errors.New("args is nil")
	}

	// 发送数据
	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_RANK)
	m.SetFunctionID(F_ID_RANK_SET_USER_RANK)

	// 调用rank设置排行信息
	mDataBuf, _ := json.Marshal(args)
	m.SetBody(mDataBuf)
	mResp, err := s.so.Send(c, m)
	if err != nil {
		// logs.Errorf("SetCommonUserRank sdkClient send to rank is failed,args:%+v,err:%v", args, err)
		return err
	}

	// 解析body
	mRespData := &jsonmodel.RspSetCommonRank{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil || mRespData.Code != 0 {
		// logs.Errorf("SetCommonUserRank Unmarshal rsp body is failed,args:%+v,code,err:%v", args, mRespData.Code, err)
		return err
	}

	// success
	// logs.Infof("SetCommonUserRank success,args:%+v,mRespData:%+v", args, mRespData)
	return nil
}

/*
ClearCommonRankList
@Desc	清除通用排行榜列表信息
@Param	gameName 	string			游戏名
@Param	raceType 	string			排行类型
*/
func (s *SdkClient) ClearCommonRankList(c context.Context, gameName string, raceType string) error {
	if s.so == nil {
		// logs.Errorf("ClearCommonRankList sdkClient is nil")
		return errors.New("ClearCommonRankList sdkClient is nil")
	}

	if gameName == "" || raceType == "" {
		str := fmt.Sprintf("some args is nil,gameName:%v,raceType:%v", gameName, raceType)
		// logs.Errorf(str)
		return errors.New(str)
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_RANK)
	m.SetFunctionID(F_ID_RANK_CLEAR_USER_RANK)

	// 准备数据
	args := &jsonmodel.ReqClearCommonRank{
		GameName: gameName,
		RaceType: raceType,
	}
	mDataBuf, _ := json.Marshal(args)
	m.SetBody(mDataBuf)

	// 调用rank服获取排行信息
	mResp, err := s.so.Send(c, m)
	if err != nil {
		// logs.Errorf("ClearCommonRankList sdkClient send to rank is failed,args:%+v,err:%v", args, err)
		return err
	}

	// rsp解析
	mRespData := &jsonmodel.RspCommonRankData{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil || mRespData.Code != 0 {
		// logs.Errorf("ClearCommonRankList Unmarshal rsp body is failed,args:%+v,code,err:%v", args, mRespData.Code, err)
		return err
	}

	// logs.Infof("ClearCommonRankList success,args:%+v,mRespData:%+v", args, mRespData)
	return nil
}

/*
GetSpecifiedRankList
@Desc	获得指定排行榜列表信息
@Param	gameName 	string		游戏名
@Param	raceType 	string		排行类型
@Param	list 		[][]int32	指定排名
@Param	userID 		string		玩家自己的ID
*/
func (s *SdkClient) GetSpecifiedRankList(c context.Context, gameName string,
	raceType string, userID string, list [][]int32) (res *jsonmodel.SpecifiedRankData, err error) {
	if s.so == nil {
		logs.Errorf("GetSpecifiedRankList sdkClient is nil")
		return nil, errors.New("GetSpecifiedRankList sdkClient is nil")
	}

	if gameName == "" || raceType == "" {
		str := fmt.Sprintf("some args is nil,gameName:%v,raceType:%v", gameName, raceType)
		logs.Errorf(str)
		return nil, errors.New(str)
	}

	m := pCore.NewRequestMessage()
	m.SetToSvrType(ST_RANK)
	m.SetFunctionID(F_ID_RANK_SPECIFIED_USER_RANK)

	// 准备数据
	args := &jsonmodel.ReqGetSpecifiedRank{
		GameName: gameName,
		RaceType: raceType,
		UserID:   userID,
		List:     list,
	}
	mDataBuf, _ := json.Marshal(args)
	m.SetBody(mDataBuf)

	// 调用rank服获取排行信息
	mResp, err := s.so.Send(c, m)
	if err != nil {
		logs.Errorf("GetSpecifiedRankList sdkClient send to rank is failed,args:%+v,err:%v", args, err)
		return nil, err
	}

	// rsp解析
	// mRespData := &jsonmodel.RspCommonRankData{}
	mRespData := &jsonmodel.RspSpecifiedRankData{}
	err = json.Unmarshal(mResp.GetBody(), mRespData)
	if err != nil || mRespData.Code != 0 {
		logs.Errorf("GetSpecifiedRankList Unmarshal rsp body is failed,args:%+v,code,err:%v", args, mRespData.Code, err)
		return nil, err
	}

	logs.Infof("GetSpecifiedRankList success,args:%+v,mRespData:%+v", args, mRespData)
	return mRespData.Data, nil
}
