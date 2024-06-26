package types

const DefaultCap = 8 // map 默认容量

// 服务类型
const (
	GAMESERVER     = 5  // 游戏服
	DISPATCHSERVER = 7  // 游戏调度服
	MAILSERVER     = 13 // 邮件服
	RANKSERVER     = 14 // 比赛排名服
)

// 接口号
const (
	F_ID_ASSIGN_ROOM_CARD_IDS = 1 // 游戏调度服分配房卡列表接口号
	F_ID_REVOKE_ROOM_CARD_IDS = 2 // 游戏调度服回收房卡列表接口号
)

// 部署类型
const (
	DeployApp   DeployType = iota //App
	DeployVGame                   // 微信小游戏
)
