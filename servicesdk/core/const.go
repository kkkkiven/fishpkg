package core

const (
	// 核心网关ETCD目录
	DEFAULT_GATEWAY_DIR = "/YunFan/framework/gw_core/"
	// 服务ETCD目录
	DEFAULT_SERVICE_DIR = "/YunFan/platform/core/"
	// 阿里日志topic
	DEFAULT_TOPIC_ALILOG = "t_alilog"
	// 链路追踪topic
	DEFAULT_TOPIC_TRACER = "t_tracer"
	// 默认超时时间
	DEFAULT_TIMEOUT = 10
	// ETCD过期时间
	DEFAULT_ETCD_EXPIRE = 300
	// ETCD保活间隔
	DEFAULT_ETCD_PERIOD = 5
	// 权重环境变量名
	DEFAULT_WEIGHT_ENV = "CORESVR_WEIGHT"
	// 服务发现方式: ETCD or STATIC
	DEFAULT_DISCOVER_MODE = "ETCD"
)

// 服务类型
const (
	ST_GW_CORE      = 1 + iota // 核心网关
	ST_GW_HTTP                 // 外部网关
	ST_LOGIN                   // 登录服务
	ST_VERSION                 // 版本检查服务
	ST_GAME                    // 游戏服务
	ST_SMS                     // 短信服务
	ST_GAMEDISPATCH            // 游戏调度服务
	ST_USER_API                // userapi接口
	ST_SHARE_API               // 分享接口
	ST_RECHARGE_API            // 充值接口
	ST_HALL_API                // 大厅接口
	ST_NOTIFY                  // 通知服务
	ST_MAIL                    // 邮件服务
	ST_RANK                    // 排名服务
	ST_AGENTAPI                // 代理服务
	ST_BACKEND                 // 后台统计系统

	ST_HALL_XQ_API = 128 + iota // 大厅接口
)

const (
	// 核心网关ID
	F_ID_REGISTER = 0x0001
	F_ID_UPDATE   = 0x0002
	F_ID_PING     = 0x1000
	F_ID_PONG     = 0x1001
)

const (
	// USERAPI 接口id定义
	F_ID_RESERVED                  uint16 = 0  // 保留
	F_ID_GET_INFO                         = 1  // 获取用户基础属性信息
	F_ID_GET_PROP                         = 2  // 获取用户道具信息
	F_ID_GET_ALL                          = 3  // 获取用户基础属性与道具信息
	F_ID_SET_INFO                         = 4  // 设置用户基础属性信息
	F_ID_SET_PROP                         = 5  // 设置用户道具信息
	F_ID_SET_ALL                          = 6  // 设置用户基础属性与道具信息
	F_ID_SET_PASS                         = 7  // 设置用户密码
	F_ID_OP_PROP                          = 8  // 操作用户道具信息
	F_ID_ENTER_GAME                       = 9  // 用户进入游戏
	F_ID_LEAVE_GAME                       = 10 // 用户离开游戏
	F_ID_CHECK_USERNAME_EXISTED           = 11 // 检查用户名是否已经被注册
	F_ID_GET_GAME_INFO                    = 12 // 获取用户游戏属性
	F_ID_SET_GAME_INFO                    = 13 // 设置用户游戏属性
	F_ID_OP_GAME_INFO                     = 14 // 操作用户游戏属性
	F_ID_GET_ADVANCE_PROP                 = 15 // 获取高级道具
	F_ID_SET_ADVANCE_PROP                 = 16 // 设置高级道具
	F_ID_OP_ADVANCE_PROP                  = 17 // 操作高级道具
	F_ID_HAS_BAD_WROD                     = 18 // 敏感词检查
	F_ID_REPLACE_BAD_WORD                 = 19 // 敏感词过滤
	F_ID_UPDATE_PASS                      = 20 // 更新用户密码
	F_ID_OP_PROP_AND_ADVANCED_PROP        = 21 // 原子更新普通道具与高级道具
	F_ID_CHECK_IDCARD_EXISTED             = 22 // 检查idcard是否存在
	F_ID_GET_USERS_BY_PHONE_NUMBER        = 23 // 通过电话号码获取用户列表
	F_ID_ADD_DEPUTY_ACCOUT                = 24 // 添加辅助账号
	F_ID_DEL_DEPUTY_ACCOUT                = 25 // 删除辅助账号
	F_ID_GET_DEPUTY_ACCOUT                = 26 // 获取辅助账号
	F_ID_CHECK_NICKNAME_EXISTED           = 27 // 检查昵称是否已经被注册

	F_ID_GET_ADVANCED_PROP_EX = 28
	F_ID_GET_ALL_EX           = 29
	F_ID_SET_ADVANCED_PROP_EX = 30
	F_ID_SET_ALL_EX           = 31

	F_ID_GET_USERS_INFO           = 32 // 批量获取用户基础属性信息
	F_ID_ENTER_GAME_EX            = 33 // 用户进入游戏，带扩展高级道具数据
	F_ID_GET_USERS_INFO_EX        = 34 // 批量获取用户基础属性&游戏属性信息
	F_ID_GET_BATCH_INFO_AND_PROP  = 35 // 批量获取用户基础信息和道具信息
	F_ID_GET_BATCH_INFOS_AND_PROP = 36 // 批量获取用户基础信息和基础道具信息和游戏属性信息

)

const (
	// 短信接口id定义
	F_ID_SMS_SEND   = 1 + iota // 发送验证码
	F_ID_SMS_VERIFY            // 校验验证码
)

const (
	// 邮件接口ID定义
	F_ID_MAIL_SEND        = 1 + iota // 发送邮件
	F_ID_MAIL_AWARD                  // 获取邮件奖励
	F_ID_MAIL_STATUS                 // 更新邮件状态
	F_ID_MAIL_BACK_SEND              // 后台发送邮件
	F_ID_MAIL_BATCH_READ             // 一键已读
	F_ID_MAIL_BATCH_AWARD            // 一键领取
	F_ID_MAIL_BATCH_DEL              // 一键领取
	F_ID_MAIL_ALL_DETAIL             // 获取邮件奖励
	F_ID_MAIL_EVALUATION             // 邮件服务评价
	F_ID_MAIL_SEND_ALL               // 发送全服邮件
)

const (
	// 支付接口ID定义
	F_ID_PAY_CHECK    = 1 + iota // 下单校验
	F_ID_PAY_CALLBACK            // 回调发货

	F_ID_PAY_CHECK_GAME    = 18001 // 下单直接游戏服校验
	F_ID_PAY_CALLBACK_GAME = 18002 // 回调游戏服发货
)

const (
	F_ID_GAME_KICKOFF = 1 + iota // 将玩家踢下线
	F_ID_GAME_GIVEUP  = 0x0B     // 玩家放弃当前比赛
)

const (
	// 排名服务接口ID定义
	F_ID_RANK_PUSH_SCORE    = 0x01 // 发送比赛积分
	F_ID_RANK_LIST          = 0x02 // 获取排行榜
	F_ID_RANK_KEY_KEY       = 0x03 // 获取排行榜通过key
	F_ID_RANK_TRANSFER_RANK = 0x04 // A排行榜用户分数转移到B排行榜
	F_ID_RANK_PUSH_RETURN   = 0x05 // 支持批量上传积分并返回排名

	F_ID_RANK_GET_LIST      = 0x07 // 获得排行榜列表（通用）
	F_ID_RANK_SET_USER_RANK = 0x08 // 设置玩家排名（通用）

	F_ID_SAVE_USER_ARENA_DATA     = 0x09 // 保存用户大奖赛数据
	F_ID_RANK_CLEAR_USER_RANK     = 0xA  // 清除排行榜列表（通用）
	F_ID_RANK_SPECIFIED_USER_RANK = 0xB  // 获得指定排行榜列表（通用）
)

const (
	// 调度服务接口ID定义
	F_ID_GAMEDISPATCH_GET_ROOMCARDIDS     = 1 + iota // 申请房卡号
	F_ID_GAMEDISPATCH_RELEASE_ROOMCARDIDS            // 玩家放弃当前比赛
)
