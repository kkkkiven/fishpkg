package jsonmodel

type RankData struct {
	UserId int64 `json:"user_id"`
	Score  int64 `json:"score"`
}

type MyRankData struct {
	UserId  int64 `json:"user_id"`
	Ranking int   `json:"ranking"`
	Score   int64 `json:"score"`
	Total   int   `json:"total"`
}

type RspRankList struct {
	Code     int32      `json:"code"`
	Msg      string     `json:"msg,omitempty"`
	RankList []RankData `json:"rank_list,omitempty"`
	User     MyRankData `json:"user,omitempty"`
	LastKey  string     `json:"last_key,omitempty"`
}

type RspRankAward struct {
	Code     int32      `json:"code"`
	Msg      string     `json:"msg,omitempty"`
	RankList []RankData `json:"rank_list,omitempty"`
}

// 通用排行 - 请求REQ设置排行榜
type ReqSetCommonRank struct {
	NeedSort      []*CommonNeedSort `json:"need_sort"`       // 排名的数据
	GameName      string            `json:"game_name"`       // 游戏名
	UserID        int64             `json:"user_id"`         // 用户ID
	NeedCacheData string            `json:"need_cache_data"` // 需要缓存的信息
	CacheExpTime  int64             `json:"cache_exp_time"`  // 缓存信息过期时间
}

// 通用排行 - 需要排的信息
type CommonNeedSort struct {
	RaceType  []string `json:"race_type"` // 排行榜类型
	Condition int64    `json:"condition"` // 进入多少名排名,注意：<=0 不限制
	Score     float64  `json:"score"`     // 积分
}

// 通用排行 - 返回RSP设置排行榜
type RspSetCommonRank struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

// 通用排行 - 请求获得通用排行信息
type ReqGetCommonRank struct {
	GameName  string `json:"game_name"` // 游戏前缀
	RaceType  string `json:"race_type"` // 需要获得的排行榜类型
	Condition int64  `json:"condition"` // 需要前多少名,注意：<=0 不限制
}

// ReqGetSpecifiedRank 通用排行 - 请求获得指定排行信息
type ReqGetSpecifiedRank struct {
	GameName string    `json:"game_name"` // 游戏前缀
	RaceType string    `json:"race_type"` // 需要获得的排行榜类型
	List     [][]int32 `json:"list"`      // 指定排名
	UserID   string    `json:"user_id"`   // 自己的userID
}

// ReqClearCommonRank 通用排行 - 清除排行榜
type ReqClearCommonRank struct {
	GameName string `json:"game_name"` // 游戏前缀
	RaceType string `json:"race_type"` // 需要清除的排行榜类型
}

// 通用排行 - 返回排行数据
type RspCommonRankData struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

// 通用排行 - 用户积分排名
type CommonSortData struct {
	Ranking int32   `json:"ranking"` // 排名
	UserID  int64   `json:"user_id"` // 用户id
	Score   float64 `json:"score"`   // 积分
}

// SpecifiedRankData 指定排名数据
type SpecifiedRankData struct {
	RankData  map[int64]*CommonSortData `json:"rank_data"`  // 用户排名信息 key=userID
	CacheData []string                  `json:"cache_data"` // 缓存数据
	OwnRank   int64                     `json:"own_rank"`   // 自己的排名 0=没上榜
}

// 通用排行 - 返回排行数据
type RspSpecifiedRankData struct {
	Code int32              `json:"code"`
	Msg  string             `json:"msg,omitempty"`
	Data *SpecifiedRankData `json:"data"`
}

type ReqSaveUserArenaData struct {
	UserID   int64  `json:"userid"`
	Score    int64  `json:"score"`    // 积分(包含加成积分)
	Addition int64  `json:"addition"` // 加成积分
	Ext      string `json:"ext"`      // 扩展参数
}

type RespSaveUserArenaData struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Rank int64  `json:"rank"`
}
