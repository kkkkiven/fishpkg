package http

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/kkkkiven/fishpkg/logs"
	"github.com/kkkkiven/fishpkg/servicesdk/http/helper"
	"github.com/kkkkiven/fishpkg/servicesdk/pkg/etcd"
	p "github.com/kkkkiven/fishpkg/sprotocol/http"
	"github.com/kkkkiven/fishpkg/utils"

	json "github.com/json-iterator/go"
)

// HConfig
type HConfig struct {
	Name         string      `yaml:"name"`          // 服务名称
	Timeout      int64       `yaml:"timeout"`       // 超时时间
	DiscoverMode string      `yaml:"discover_mode"` // 服务发现方式：static 或者 etcd
	GatewayDir   string      `yaml:"gateway_dir"`   // 网关ETCD目录
	ServiceDir   string      `yaml:"service_dir"`   // 服务ETCD目录
	GatewayAddr  []string    `yaml:"gateway_addr"`  // 网关地址
	ProberAddr   string      `yaml:"prober_addr"`   // 探测地址
	Etcd         *EntityEtcd `yaml:"etcd"`          // ETCD配置
}

// EntityEtcd
type EntityEtcd struct {
	Addrs []string `yaml:"addrs"`
	User  string   `yaml:"user"`
	Pass  string   `yaml:"pass"`
}

type _Service struct {
	sync.RWMutex

	// 服务相关
	id           uint32
	ip           string
	name         string
	weight       int
	secret       string
	discoverMode string
	proberAddr   string
	gatewayAddr  []string
	iface        []_Iface
	timeout      int64

	// ETCD相关
	gatewayDir string
	serviceDir string
	etcdConf   *EntityEtcd
	etcdConn   *etcd.Client
}

var srv *_Service = &_Service{}

func GetService() *_Service {
	return srv
}

type option func(s *_Service)

func SetName(n string) option {
	return func(s *_Service) {
		s.name = n
	}
}

func SetWeight(w int) option {
	return func(s *_Service) {
		s.weight = w
	}
}

func SetSecret(scrt string) option {
	return func(s *_Service) {
		s.secret = scrt
	}
}

func SetTimeout(t int64) option {
	return func(s *_Service) {
		s.timeout = t
	}
}

func SetDiscoverMode(m string) option {
	return func(s *_Service) {
		s.discoverMode = m
	}
}

func SetProberAddr(addr string) option {
	return func(s *_Service) {
		s.proberAddr = addr
	}
}

func SetGatewayAddr(ips []string) option {
	return func(s *_Service) {
		s.gatewayAddr = nil
		s.gatewayAddr = append(s.gatewayAddr, ips...)
	}
}

func SetGatewayDir(gdir string) option {
	return func(s *_Service) {
		s.gatewayDir = gdir
	}
}

func SetServiceDir(sdir string) option {
	return func(s *_Service) {
		s.serviceDir = sdir
	}
}

func SetEtcdConf(c *EntityEtcd) option {
	return func(s *_Service) {
		if c == nil || len(c.Addrs) == 0 {
			return
		}
		s.etcdConf = &EntityEtcd{}
		s.etcdConf.User = c.User
		s.etcdConf.Pass = c.Pass
		s.etcdConf.Addrs = make([]string, len(c.Addrs))
		copy(s.etcdConf.Addrs, c.Addrs)
	}
}

func (s *_Service) Id() uint32 {
	s.RLock()
	defer s.RUnlock()

	return s.id
}

func (s *_Service) Name() string {
	s.RLock()
	defer s.RUnlock()

	return s.name
}

func (s *_Service) Weight() int {
	s.RLock()
	defer s.RUnlock()

	return s.weight
}

func (s *_Service) Secret() string {
	s.RLock()
	defer s.RUnlock()

	return s.secret
}

func (s *_Service) GatewayDir() string {
	s.RLock()
	defer s.RUnlock()

	return s.gatewayDir
}

func (s *_Service) ServiceDir() string {
	s.RLock()
	defer s.RUnlock()

	return s.serviceDir
}

func (s *_Service) Iface() []_Iface {
	s.RLock()
	defer s.RUnlock()

	return s.iface
}

func (s *_Service) Timeout() int64 {
	s.RLock()
	defer s.RUnlock()

	return s.timeout
}

func (s *_Service) EtcdConf() *EntityEtcd {
	s.RLock()
	defer s.RUnlock()

	return s.etcdConf
}

func (s *_Service) EtcdConn() *etcd.Client {
	s.RLock()
	defer s.RUnlock()

	return s.etcdConn
}

func (s *_Service) DiscoverMode() string {
	s.RLock()
	defer s.RUnlock()

	return s.discoverMode
}

func (s *_Service) GatewayAddr() []string {
	s.RLock()
	defer s.RUnlock()

	return s.gatewayAddr
}

func (s *_Service) AddIface(iface _Iface) {
	s.Lock()
	defer s.Unlock()

	s.iface = append(s.iface, iface)
}

func Init(cfg *HConfig) error {
	if cfg == nil {
		return errors.New("the config is nil")
	}

	if cfg.Name == "" {
		return errors.New("the server name is empty")
	}

	srv.Lock()
	defer srv.Unlock()

	srv.discoverMode = cfg.DiscoverMode
	if srv.discoverMode == "" {
		srv.discoverMode = DEFAULT_DISCOVER_MODE
	}

	if strings.ToUpper(srv.discoverMode) == DEFAULT_DISCOVER_MODE {
		if cfg.Etcd == nil || len(cfg.Etcd.Addrs) == 0 {
			return errors.New("bad etcd config")
		}

		conn, err := etcd.New(cfg.Etcd.Addrs, cfg.Etcd.User, cfg.Etcd.Pass)
		if err != nil {
			return errors.New(fmt.Sprintf("new etcd err: %v", err.Error()))
		}

		srv.etcdConf = &EntityEtcd{}
		srv.etcdConf.User = cfg.Etcd.User
		srv.etcdConf.Pass = cfg.Etcd.Pass
		srv.etcdConf.Addrs = make([]string, len(cfg.Etcd.Addrs))
		copy(srv.etcdConf.Addrs, cfg.Etcd.Addrs)
		srv.etcdConn = conn
	}

	srv.ip = helper.GetLocalAddress("")
	if srv.ip == "" {
		logs.Waring("get local address failed")
	}

	srv.gatewayAddr = cfg.GatewayAddr
	srv.name = cfg.Name
	srv.timeout = cfg.Timeout
	srv.id = utils.Ip2long(srv.ip)
	srv.weight = utils.Atoi(os.Getenv(DEFAULT_WEIGHT_ENV))
	logs.Infof("Node weight: [%v:%v]", DEFAULT_WEIGHT_ENV, srv.weight)

	srv.gatewayDir = cfg.GatewayDir
	if srv.gatewayDir == "" {
		srv.gatewayDir = DEFAULT_GATEWAY_DIR
	}
	if srv.gatewayDir[len(srv.gatewayDir)-1] != '/' {
		srv.gatewayDir += "/"
	}

	srv.serviceDir = cfg.ServiceDir
	if srv.serviceDir == "" {
		srv.serviceDir = DEFAULT_SERVER_DIR
	}
	if srv.serviceDir[len(srv.serviceDir)-1] != '/' {
		srv.serviceDir += "/"
	}

	return nil
}

func InitOpt(opts ...option) error {
	srv.Lock()
	defer srv.Unlock()

	for _, opt := range opts {
		opt(srv)
	}

	if srv.discoverMode == "" {
		srv.discoverMode = DEFAULT_DISCOVER_MODE
	}

	if strings.ToUpper(srv.discoverMode) == DEFAULT_DISCOVER_MODE {
		if srv.etcdConf == nil || len(srv.etcdConf.Addrs) == 0 {
			return errors.New("bad etcd config")
		}

		conn, err := etcd.New(srv.etcdConf.Addrs, srv.etcdConf.User, srv.etcdConf.Pass)
		if err != nil {
			return errors.New(fmt.Sprintf("new etcd err: %v", err.Error()))
		}
		srv.etcdConn = conn
	}

	if srv.name == "" {
		return errors.New("the server name is empty")
	}

	ip := helper.GetLocalAddress("")
	if ip == "" {
		logs.Waring("get local address failed")
	}

	srv.id = utils.Ip2long(ip)
	srv.weight = utils.Atoi(os.Getenv(DEFAULT_WEIGHT_ENV))

	if srv.gatewayDir == "" {
		srv.gatewayDir = DEFAULT_GATEWAY_DIR
	}
	if srv.gatewayDir[len(srv.gatewayDir)-1] != '/' {
		srv.gatewayDir += "/"
	}

	if srv.gatewayDir == "" {
		srv.serviceDir = DEFAULT_SERVER_DIR
	}
	if srv.serviceDir[len(srv.serviceDir)-1] != '/' {
		srv.serviceDir += "/"
	}

	return nil
}

func Run() error {
	if strings.ToUpper(srv.DiscoverMode()) == DEFAULT_DISCOVER_MODE {
		if srv.EtcdConf() == nil {
			return errors.New("please init first")
		}

		if len(srv.Iface()) == 0 {
			return errors.New("please add handler first")
		}

		if err := srv.publish(); err != nil {
			return err
		}

		return fetchGateway()
	}

	return dialGateway(srv.GatewayAddr())
}

func Stop() error {
	return srv.revoke()
}

// ==========================================

const (
	RC_OK int = iota
	RC_REJECT
)

type _RegRequest struct {
	Name   string   `json:"name"`
	Ip     string   `json:"ip"`
	Weight int      `json:"weight"`
	Secret string   `json:"secret"`
	Iface  []_Iface `json:"iface"`
}

type _GWResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type _Iface struct {
	Path string `json:"path"`
	Desc string `json:"desc"`
}

type _UnRegRequest struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type _UpdateRequest struct {
	Name   string `json:"name"`
	Ip     string `json:"ip"`
	Weight int    `json:"weight"`
}

func (s *_Service) register(so *p.Socket) error {
	s.RLock()
	regReq := &_RegRequest{}
	regReq.Name = srv.name
	regReq.Weight = srv.weight
	regReq.Secret = srv.secret
	regReq.Iface = srv.iface
	regReq.Ip = srv.ip
	s.RUnlock()

	body, _ := json.Marshal(regReq)

	reqMsg := p.NewRequestMessage()
	reqMsg.SetBody(body)
	reqMsg.SetStringExData("/gw/register")

	rspMsg, err := so.Send(reqMsg)
	if err != nil {
		return err
	}

	regRsp := &_GWResponse{}
	if err := json.Unmarshal(rspMsg.GetBody(), regRsp); err != nil {
		return err
	}

	if regRsp.Code != RC_OK {
		return errors.New(regRsp.Msg)
	}

	return nil
}

// func (s *_Service) update(so *p.Socket) error {
// 	s.RLock()
// 	upReq := &_UpdateRequest{}
// 	upReq.Weight = srv.weight
// 	upReq.Ip = srv.ip
// 	s.RUnlock()

// 	body, _ := json.Marshal(upReq)

// 	reqMsg := p.NewRequestMessage()
// 	reqMsg.SetBody(body)
// 	reqMsg.SetStringExData("/gw/update")

// 	rspMsg, err := so.Send(reqMsg)
// 	if err != nil {
// 		return err
// 	}

// 	upRsp := &_GWResponse{}
// 	if err := json.Unmarshal(rspMsg.GetBody(), upRsp); err != nil {
// 		return err
// 	}

// 	if upRsp.Code != RC_OK {
// 		return errors.New(upRsp.Msg)
// 	}

// 	return nil
// }

// ================================

type _Publish struct {
	Name   string `json:"name"`
	Ip     string `json:"ip"`
	Weight int    `json:"weight"`
}

func (s *_Service) publish() error {
	s.RLock()
	pub := &_Publish{}
	pub.Name = s.name
	pub.Weight = s.weight
	pub.Ip = s.ip
	key := fmt.Sprintf("%vhttp_%v", s.serviceDir, s.id)
	etcdConn := s.etcdConn
	s.RUnlock()

	body, _ := json.Marshal(pub)
	if err := etcdConn.SetKvAndKeepAlive(key, string(body), ETCD_KEY_EXPIRE, ETCD_KEY_KEEPALIVE_PERIOD); err != nil {
		return err
	}

	// go s.watch()

	logs.Debugf("Publish to etcd: [key=%v,value=%v]", key, string(body))
	return nil
}

func (s *_Service) revoke() error {
	s.RLock()
	key := fmt.Sprintf("%vhttp_%v", s.serviceDir, s.id)
	etcdConn := s.etcdConn
	s.RUnlock()

	return etcdConn.Del(key)
}

// func (s *_Service) watch() {

// 	s.RLock()
// 	pfx := fmt.Sprintf("%vhttp_%v", s.serviceDir, s.id)
// 	rch := s.etcdConn.WatchKeyWithPrefix(pfx)
// 	s.RUnlock()

// 	for wresp := range rch {
// 		if wresp.Canceled {
// 			logs.Infof("Watch %v canceled", pfx)
// 			return
// 		}

// 		if err := wresp.Err(); err != nil {
// 			logs.Errorf("Watch %v err: %s", pfx, err.Error())
// 		}

// 		for _, ev := range wresp.Events {
// 			switch ev.Type.String() {
// 			case "PUT":
// 				logs.Infof("Add ETCD [key:%s,value:%s]", string(ev.Kv.Key), string(ev.Kv.Value))
// 				pub := &_Publish{}
// 				if err := json.Unmarshal(ev.Kv.Value, pub); err != nil {
// 					logs.Errorf(err.Error())
// 					continue
// 				}

// 				conns := gwList.GetAll()
// 				for _, conn := range conns {
// 					if conn == nil {
// 						continue
// 					}
// 					if err := s.update(conn); err != nil {
// 						logs.Errorf("Update weight err: %s", err.Error())
// 					}
// 				}
// 			case "DELETE":
// 				logs.Infof("Delete ETCD [key:%s,value:%s]", string(ev.Kv.Key), string(ev.Kv.Value))
// 			}
// 		}
// 	}
// }
