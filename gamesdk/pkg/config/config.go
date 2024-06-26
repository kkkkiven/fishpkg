package config

import (
	"sync/atomic"

	. "github.com/kkkkiven/fishpkg/gamesdk/pkg/types"
)

// EntityServer 服务配置
type EntityServer struct {
	GameName  string     `yaml:"game_name" json:"game_name"`
	Port      int        `yaml:"port" json:"port"`           // 监听端口
	WsPort    int        `yaml:"ws_port" json:"ws_port"`     // websocket 端口
	MaxConns  int32      `yaml:"max_conns" json:"max_conns"` // 最大连接数
	Deploy    DeployType `yaml:"deploy" json:"deploy"`       // 部署类型
	TimeOut   int        `yaml:"timeout" json:"time_out"`
	Authorize bool       `yaml:"authorize" json:"authorize"` // 是否验证连接,启用验证后，第一包必须是验证包，否则连接将被断开
}

var (
	confSever atomic.Value
)

func Init(cfgServer EntityServer) {
	confSever.Store(cfgServer)
}

func GetGameName() string {
	return confSever.Load().(EntityServer).GameName
}

func GetPort() int {
	return confSever.Load().(EntityServer).Port
}

func GetWsPort() int {
	return confSever.Load().(EntityServer).WsPort
}

func GetMaxConns() int32 {
	return confSever.Load().(EntityServer).MaxConns
}

func GetDeploy() DeployType {
	return confSever.Load().(EntityServer).Deploy
}

func GetTimeOut() int {
	return confSever.Load().(EntityServer).TimeOut
}

func Authorize() bool {
	return confSever.Load().(EntityServer).Authorize
}
