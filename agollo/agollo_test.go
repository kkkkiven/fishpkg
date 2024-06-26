package agollo

import (
	`encoding/json`
	`os`
	`os/signal`
	`syscall`
	`testing`

	`git.yuetanggame.com/zfish/fishpkg/logs`
	`git.yuetanggame.com/zfish/fishpkg/servicesdk/core`
	"git.yuetanggame.com/zfish/fishpkg/servicesdk/http"
)

// Conf 阿波罗配置
type BConf struct {
	Host           *ConfHost              `yaml:"base.yaml"`       // 宿主基础配置
	ServiceSdkHTTP *http.HConfig          `yaml:"server.yaml"`     // servicesdk-http 配置
	ServiceSdkCore *core.CConfig          `yaml:"core.yaml"`       // servicesdk-core 配置
	NotifyMail     map[string]*NotifyMail `yaml:"notifyMail.yaml"` // 邮件配置
}

// NotifyMail 邮件通知配置
type NotifyMail struct {
	Title          string          `yaml:"title"`           // 邮件标题
	Content        string          `yaml:"content"`         // 邮件内容
	Sender         string          `yaml:"sender"`          // 发送者
	StartTime      string          `yaml:"start_time"`      // 邮件发送开始时间
	EndTime        string          `yaml:"end_time"`        // 邮件发送结束时间
	Rewards        map[int64]int64 `yaml:"rewards"`         // 邮件奖励
	ChannelValid   []int32         `yaml:"channel_valid"`   // 需要发送的渠道（-1=全部渠道）
	ChannelInvalid []int32         `yaml:"channel_invalid"` // 不需要发送的渠道
}

// OverflowWarningData 预警相关数据
type OverflowWarningData struct {
	PostUrl     string          `yaml:"post_url"`      // 请求地址
	VersionName string          `yaml:"version_name"`  // 版本昵称
	CoolingTime int64           `yaml:"cooling_time"`  // 预警冷却时间
	WhiteListIP []string        `yaml:"white_list_ip"` // 白名单IP
	WhiteListID []int64         `yaml:"white_list_id"` // 白名单ID
	PropsLimit  map[int32]int64 `yaml:"props_limit"`   // 道具预警数量
}

// ConfHost 宿主配置定义实体
type ConfHost struct {
	RunMode            string      `yaml:"runmode"`                // 服务运行模式
	EnvStatus          int         `yaml:"env_status"`             // 0 测试 1 审核 2正式
	IsFuFeiDeBug       bool        `yaml:"is_fufei_debug"`         // 是否免支付
	OldWeekCardBuyTime string      `yaml:"old_week_card_buy_time"` // 周卡旧的购买时间
	Old1GrowthGiftbag  string      `yaml:"old1_growth_giftbag"`    // 旧1成长基金至尊奖励
	Old2GrowthGiftbag  string      `yaml:"old2_growth_giftbag"`    // 旧2成长基金至尊奖励
	Logs               EntityLogs  `yaml:"logs"`                   // 日志
	MysqlMainDB        EntityMysql `yaml:"mysql_main"`             // 数据库 注意这个是main
	MysqlUserDB2       EntityMysql `yaml:"mysql_user"`             // 数据库
	MysqlLogDB         EntityMysql `yaml:"mysql_log"`              // 数据库
	Redis              EntityRedis `yaml:"redis"`                  // redis
	Redis8             EntityRedis `yaml:"redis8"`                 // redis：8
	Mongo              EntityMongo `yaml:"mongo"`                  // mongo
	KafkaTopicNotify   string      `yaml:"kafka_topic_notify"`     // topic notify
	KafkaProject       string      `yaml:"kafka_project"`          // kafka_project
	EtcdDir            string      `yaml:"etcd_dir"`               // etcd配置根目录
	BaiDuOcpcAkey      string      `yaml:"baidu_ocpc_akey"`        // 百度akey
	BrandName          string      `yaml:"brand_name"`             // 品牌
	// Server           cabsdk.Config      `yaml:"server"`             // cabsdk 配置
	PayMchid     string             `yaml:"pay_mchid"`      // 支付系统-mchid
	PayAppid     string             `yaml:"pay_appid"`      // 支付系统-appid
	FriendNum    int64              `yaml:"friend_num"`     // 好友上限数量
	AliOss       EntityAliOss       `yaml:"alioss"`         // 阿里云oss配置
	MysqlOrderDB EntityMysql        `yaml:"mysql_order"`    // 数据库
	ImServer     EntityImServer     `yaml:"imserver"`       // im配置
	RuixueServer EntityRuixueServer `yaml:"ruixueserver"`   // 瑞雪配置
	RuiXueCP     EntityRuiXueCP     `yaml:"ruixue_cp"`      // 瑞雪v2应用信息配置
	OrderAddrURL string             `yaml:"order_addr_url"` // 获取商城省市地区url地址（暂时不用-未配置）
	ChannelCfg   map[string][]int32 `yaml:"channel_cfg"`    // 渠道配置（暂时不用-未配置）
	AesKey       string             `yaml:"aes_key"`        // db解密需要的key
	BaiduAkey    map[int32]string   `yaml:"baidu_akey"`     // 百度akey

	PayNoNeedCert map[int32][]int32 `yaml:"pay_no_need_cert"` // 支付不需要实名认证
	IsApple       int32             `yaml:"is_apple"`         // 是否是苹果
	OrderMsg      int32             `yaml:"order_msg"`        // 订单信息0正式 2审核 3test

	Sync37240Secret string `yaml:"sync37240_secret"`  // 迁移37240
	Sync37240Url    string `yaml:"sync37240_url"`     // 迁移37240-获取用户信息-吉祥
	Sync37240UrlWL  string `yaml:"sync37240_url_wl"`  // 迁移37240-获取用户信息-微乐
	Sync37240Code   string `yaml:"sync37240_code"`    // 迁移37240-同步兑换码-吉祥
	Sync37240CodeWL string `yaml:"sync37240_code_wl"` // 迁移37240-同步兑换码-微乐

	ChargeMonthLimit     int64  `yaml:"charge_month_limit"`      // 月充值金额限制（单位分）
	ChargeDayLimit       int64  `yaml:"charge_day_limit"`        // 日充值金额限制（单位分）
	ThreeRewardsShowTime string `yaml:"three_rewards_show_time"` // 3日登陆什么时候注册之后显示奖励
	CompensateGunTime    string `yaml:"compensate_gun_time"`     // 炮倍补偿哪个时间之前注册

	OverflowWarning OverflowWarningData `yaml:"overflow_warning"` // 订单信息0正式 2审核 3test
	RXService       EntifyRXService     `yaml:"rxservice"`        // 瑞雪接口服务配置

	WechatMonthLimit   int64             `yaml:"wechat_month_limit"`   // 微信支付限额 - 月充值金额限制（单位分）
	WechatDayLimit     int64             `yaml:"wechat_day_limit"`     // 微信支付限额 - 日充值金额限制（单位分）
	WechatLimitStart   string            `yaml:"wechat_limit_start"`   // 微信支付限额 - 有效开始时间
	WechatLimitChannel map[int32][]int32 `yaml:"wechat_limit_channel"` // 微信支付限额 -
	RedisLock          bool              `yaml:"redis_lock"`           // 是否开启http请求的redis同步锁
	NotifyUrl          string            `yaml:"notify_url"`           // 给瑞雪的通知服回调地址
	RuixueSecret       string            `yaml:"ruixue_secret"`        // 瑞雪签名
}

// EntifyRXService 瑞雪接口服务配置实体
type EntifyRXService struct {
	Domain   string `yaml:"domain"`
	Key      string `yaml:"key"`
	Platform string `yaml:"platform"`
}

// EntityLogs logs实体
type EntityLogs struct {
	Dir      string `yaml:"dir"`      // 文件保存路径
	File     string `yaml:"file"`     // 文件名称,实际会保存为{filename}+{datetime}
	Level    int    `yaml:"level"`    // 日志等级
	SaveFile bool   `yaml:"savefile"` // 是否保存为文件
}

// EntityMysql 实体
type EntityMysql struct {
	Addr     string `yaml:"addr"`     // 数据库信息
	UserName string `yaml:"user"`     // 用户
	Password string `yaml:"password"` // 密码
	Db       string `yaml:"db"`       // 数据库
	Charset  string `yaml:"charset"`  // 字符集
	MaxOpen  int    `yaml:"max_open"` // 最大连接数
	MaxIdle  int    `yaml:"max_idle"` // 最大空闲数
}

// EntityRedis redis实体
type EntityRedis struct {
	Addr        string `yaml:"addr"`
	Pwd         string `yaml:"password"`
	DB          int    `yaml:"db"`
	PoolSize    int    `yaml:"pool_size"`
	PoolTimeout int    `yaml:"pool_timeout"`
}

// EntityMongo mongo 实体
type EntityMongo struct {
	Addrs     []string `yaml:"addrs"`
	UserName  string   `yaml:"user"`
	Password  string   `yaml:"password"`
	Db        string   `yaml:"db"`
	PoolLimit int      `yaml:"pool_limit"`
	Timeout   int      `yaml:"timeout"`
}

// EntityAliOss 阿里云oss配置
type EntityAliOss struct {
	Endpoint             string `yaml:"endpoint"`
	AccessKeyID          string `yaml:"access_key_id"`
	AccessKeySecret      string `yaml:"access_key_secret"`
	Bucket               string `yaml:"bucket"`
	Buckettow            string `yaml:"bucket2"`
	BucketUrl            string `yaml:"bucket_url"`
	Domain               string `yaml:"domain"`
	DomainUrl            string `yaml:"domain_url"`
	HttpHead             string `yaml:"http_head"`
	Document             string `yaml:"document"`
	GreenAccessKeyID     string `yaml:"green_access_key_id"` // 内容安全：鉴黄
	GreenAccessKeySecret string `yaml:"green_access_key_secret"`
}

// EntityImServer imsdk请求配置
type EntityImServer struct {
	AppID     string `yaml:"appid"`
	AppSecret string `yaml:"appsecret"`
	URL       string `yaml:"url"`
}

// EntityRuixueServer 瑞雪配置
type EntityRuixueServer struct {
	AppID     string `yaml:"appid"`
	AppSecret string `yaml:"appsecret"`
	API       string `yaml:"api"`
}

// 瑞雪v2应用配置
type EntityRuiXueCP struct {
	Domain string `yaml:"domain"`
	CPID   uint32 `yaml:"cp_id"`
	CPKey  string `yaml:"cp_key"`
}

func TestStartAndUnmarshal(t *testing.T) {
	var conf = &BConf{}
	StartAndUnmarshalOnChange(conf, func(event *ChangeEvent, err error) {
		if err != nil {
			logs.Errorf("Update config err: %s", err.Error())
			return
		}
		str, _ := json.Marshal(conf)
		logs.Debugf("Update config : %v", string(str))
	})
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	_ = <-signals
}
