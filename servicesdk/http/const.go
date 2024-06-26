package http

// ETCD
const (
	ETCD_KEY_EXPIRE           = 300
	ETCD_KEY_KEEPALIVE_PERIOD = 5

	DEFAULT_GATEWAY_DIR   = "/YunFan/framework/gw_http/"
	DEFAULT_SERVER_DIR    = "/YunFan/platform/http/"
	DEFAULT_DISCOVER_MODE = "ETCD"
)

// 服务类型
const (
	ST_GW_CORE       = 1 + iota // 核心网关
	ST_GW_HTTP                  // 外部网关
	ST_LOGIN                    // 登录服务
	ST_VERSION                  // 版本检查服务
	ST_GAME                     // 游戏服务
	ST_SMS                      // 短信服务
	ST_GAME_DISPATCH            // 游戏调度服务
	ST_USER_API                 // userapi接口
	ST_SHARE_API                // 分享接口
	ST_RECHARGE_API             // 充值接口
	ST_HALL_API                 // 大厅接口
)

// 链路日志TOPIC
const DEFAULT_TOPIC_TRACE = "t_tracer"

// 权重环境变量名
const DEFAULT_WEIGHT_ENV = "HTTPSVR_WEIGHT"
