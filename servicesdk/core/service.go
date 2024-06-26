package core

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/kkkkiven/fishpkg/logs"
	. "github.com/kkkkiven/fishpkg/servicesdk/core/pb/core"
	"github.com/kkkkiven/fishpkg/servicesdk/pkg/etcd"
	"github.com/kkkkiven/fishpkg/servicesdk/pkg/kafka"
	u "github.com/kkkkiven/fishpkg/servicesdk/pkg/utils"
	p "github.com/kkkkiven/fishpkg/sprotocol/core"
	t "github.com/kkkkiven/fishpkg/sprotocol/tracer"
	"github.com/kkkkiven/fishpkg/utils"
)

// EntityEtcd etcd配置项
type EntityEtcd struct {
	Addrs []string `yaml:"addrs" json:"addr"`
	User  string   `yaml:"user" json:"-"`
	Pass  string   `yaml:"pass" json:"-"`
}

// EntityKafka kafka配置项
type EntityKafka struct {
	TracerTopic string   `yaml:"tracer_topic" json:"tracer_topic"`
	AliLogTopic string   `yaml:"alilog_topic" json:"alilog_topic"`
	Brokers     []string `yaml:"brokers" json:"brokers"`
}

type CConfig struct {
	Type         uint16       `yaml:"type" json:"type"`
	Name         string       `yaml:"name" json:"name"`
	DiscoverMode string       `yaml:"discover_mode" json:"discover_mode"`
	Secret       string       `yaml:"secret" json:"-"`
	ProberAddr   string       `yaml:"prober_addr" json:"prober_addr"`
	GatewayAddr  []string     `yaml:"gateway_addr" json:"gateway_addr"`
	TraceRate    int          `yaml:"trace_rate" json:"trace_rate"`
	Timeout      int64        `yaml:"timeout" json:"timeout"`
	Pack         bool         `yaml:"pack" json:"pack"`
	GatewayDir   string       `yaml:"gateway_dir" json:"gateway_dir"`
	ServiceDir   string       `yaml:"service_dir" json:"service_dir"`
	Etcd         *EntityEtcd  `yaml:"etcd" json:"etcd"`
	Kafka        *EntityKafka `yaml:"kafka" json:"kafka"`
}

type _Service struct {
	sync.RWMutex

	id           uint32
	typ          uint16
	name         string
	secret       string
	weight       int
	ip           string
	proberAddr   string
	gatewayAddr  []string
	discoverMode string
	traceRate    int
	timeout      int64
	pack         bool

	// ETCD相关
	gatewayDir string
	serviceDir string
	etcdConf   *EntityEtcd
	etcdConn   *etcd.Client

	// kafka
	kafkaConf     *EntityKafka
	kafkaProducer *kafka.AsyncProducer
}

type option func(*_Service)

func SetType(t uint16) option {
	return func(s *_Service) {
		s.typ = t
	}
}

func SetName(n string) option {
	return func(s *_Service) {
		s.name = n
	}
}

func SetSecret(secret string) option {
	return func(s *_Service) {
		s.secret = secret
	}
}

func SetProberAddr(p string) option {
	return func(s *_Service) {
		s.proberAddr = p
	}
}

func SetTraceRate(r int) option {
	return func(s *_Service) {
		s.traceRate = r
	}
}

func SetTimeout(t int64) option {
	return func(s *_Service) {
		s.timeout = t
	}
}

func SetPack(p bool) option {
	return func(s *_Service) {
		s.pack = p
	}
}

func SetGatewayDir(d string) option {
	return func(s *_Service) {
		s.gatewayDir = d
	}
}

func SetServiceDir(d string) option {
	return func(s *_Service) {
		s.serviceDir = d
	}
}

func SetEtcdConf(c *EntityEtcd) option {
	return func(s *_Service) {
		if c == nil || len(c.Addrs) == 0 {
			return
		}

		s.etcdConf = new(EntityEtcd)
		s.etcdConf.Addrs = append(s.etcdConf.Addrs, c.Addrs...)
		s.etcdConf.User = c.User
		s.etcdConf.Pass = c.Pass
	}
}

func SetKafkaConf(c *EntityKafka) option {
	return func(s *_Service) {
		if c == nil || len(c.Brokers) == 0 {
			return
		}

		s.kafkaConf = new(EntityKafka)
		s.kafkaConf.TracerTopic = c.TracerTopic
		s.kafkaConf.AliLogTopic = c.AliLogTopic
		s.kafkaConf.Brokers = append(s.kafkaConf.Brokers, c.Brokers...)

		if s.kafkaConf.TracerTopic == "" {
			s.kafkaConf.TracerTopic = DEFAULT_TOPIC_TRACER
		}

		if s.kafkaConf.AliLogTopic == "" {
			s.kafkaConf.AliLogTopic = DEFAULT_TOPIC_ALILOG
		}
	}
}

func SetDiscoverMode(m string) option {
	return func(s *_Service) {
		s.discoverMode = m
	}
}

func SetGatewayAddr(ips []string) option {
	return func(s *_Service) {
		s.gatewayAddr = nil
		s.gatewayAddr = append(s.gatewayAddr, ips...)
	}
}

func (s *_Service) Id() uint32 {
	s.RLock()
	defer s.RUnlock()

	return s.id
}

func (s *_Service) Type() uint16 {
	s.RLock()
	defer s.RUnlock()

	return s.typ
}

func (s *_Service) Name() string {
	s.RLock()
	defer s.RUnlock()

	return s.name
}

func (s *_Service) Secret() string {
	s.RLock()
	defer s.RUnlock()

	return s.secret
}

func (s *_Service) Weight() int {
	s.RLock()
	defer s.RUnlock()

	return s.weight
}

func (s *_Service) Ip() string {
	s.RLock()
	defer s.RUnlock()

	return s.ip
}

func (s *_Service) ProberAddr() string {
	s.RLock()
	defer s.RUnlock()

	return s.proberAddr
}

func (s *_Service) TraceRate() int {
	s.RLock()
	defer s.RUnlock()

	return s.traceRate
}

func (s *_Service) Timeout() int64 {
	s.RLock()
	defer s.RUnlock()

	return s.timeout
}

func (s *_Service) Pack() bool {
	s.RLock()
	defer s.RUnlock()

	return s.pack
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

func (s *_Service) EtcdConn() *etcd.Client {
	s.RLock()
	defer s.RUnlock()

	return s.etcdConn
}

func (s *_Service) KafkaProducer() *kafka.AsyncProducer {
	s.RLock()
	defer s.RUnlock()

	return s.kafkaProducer
}

func (s *_Service) AliLogTopic() string {
	s.RLock()
	defer s.RUnlock()

	return s.kafkaConf.AliLogTopic
}

func (s *_Service) TracerTopic() string {
	s.RLock()
	defer s.RUnlock()

	return s.kafkaConf.TracerTopic
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

func GetService() *_Service {
	return srv
}

var srv *_Service = &_Service{}

func Init(c *CConfig) error {
	srv.Lock()
	defer srv.Unlock()

	srv.typ = c.Type
	if srv.typ == 0 {
		return errors.New("bad service type")
	}

	srv.name = c.Name
	if srv.name == "" {
		return errors.New("bad service name")
	}

	srv.timeout = c.Timeout
	if srv.timeout == 0 {
		srv.timeout = DEFAULT_TIMEOUT
	}

	srv.gatewayDir = c.GatewayDir
	if srv.gatewayDir == "" {
		srv.gatewayDir = DEFAULT_GATEWAY_DIR
	}
	if srv.gatewayDir[len(srv.gatewayDir)-1] != '/' {
		srv.gatewayDir += "/"
	}

	srv.serviceDir = c.ServiceDir
	if srv.serviceDir == "" {
		srv.serviceDir = DEFAULT_SERVICE_DIR
	}
	if srv.serviceDir[len(srv.serviceDir)-1] != '/' {
		srv.serviceDir += "/"
	}

	srv.discoverMode = c.DiscoverMode
	if srv.discoverMode == "" {
		srv.discoverMode = DEFAULT_DISCOVER_MODE
	}

	srv.ip = u.GetLocalAddress(srv.proberAddr)
	if srv.ip == "" {
		return errors.New("get local address failed")
	}

	srv.gatewayAddr = c.GatewayAddr
	srv.proberAddr = c.ProberAddr
	srv.id = utils.Ip2long(srv.ip)
	srv.pack = c.Pack
	srv.secret = c.Secret
	srv.traceRate = c.TraceRate

	srv.weight = utils.Atoi(os.Getenv(DEFAULT_WEIGHT_ENV))
	logs.Infof("Node weight: [%v:%v]", DEFAULT_WEIGHT_ENV, srv.weight)

	var err error
	if strings.ToUpper(srv.discoverMode) == DEFAULT_DISCOVER_MODE {
		if c.Etcd == nil || len(c.Etcd.Addrs) == 0 {
			return errors.New("bad ETCD config")
		}
		srv.etcdConf = new(EntityEtcd)
		srv.etcdConf.User = c.Etcd.User
		srv.etcdConf.Pass = c.Etcd.Pass
		srv.etcdConf.Addrs = append(srv.etcdConf.Addrs, c.Etcd.Addrs...)
		if srv.etcdConn, err = etcd.New(srv.etcdConf.Addrs, srv.etcdConf.User, srv.etcdConf.Pass); err != nil {
			return err
		}
	}

	if c.Kafka == nil || len(c.Kafka.Brokers) == 0 {
		return errors.New("bad KAFKA config")
	}
	srv.kafkaConf = new(EntityKafka)

	srv.kafkaConf.AliLogTopic = c.Kafka.AliLogTopic
	if srv.kafkaConf.AliLogTopic == "" {
		srv.kafkaConf.AliLogTopic = DEFAULT_TOPIC_ALILOG
	}

	srv.kafkaConf.TracerTopic = c.Kafka.TracerTopic
	if srv.kafkaConf.TracerTopic == "" {
		srv.kafkaConf.TracerTopic = DEFAULT_TOPIC_TRACER
	}

	srv.kafkaConf.Brokers = append(srv.kafkaConf.Brokers, c.Kafka.Brokers...)

	onSuccess := func(msg *sarama.ProducerMessage) {
		value, err := msg.Value.Encode()
		if err != nil {
			logs.Errorf("Send kafka msg err: %s", err.Error())
		} else {
			logs.Tracef("Send kafka msg: %s", string(value))
		}
	}

	onError := func(err *sarama.ProducerError) {
		logs.Errorf("Send kafka msg err: %v", err.Error())
	}

	srv.kafkaProducer, err = kafka.NewAsyncProducer(srv.kafkaConf.Brokers, 10*time.Second, onSuccess, onError)
	if err != nil {
		return err
	}

	topic := srv.kafkaConf.TracerTopic
	fp := func(msg []byte, key string) error {
		return srv.kafkaProducer.SendMessage(topic, msg, key)
	}

	t.Init(srv.typ, srv.id, srv.name, srv.traceRate, fp)

	return nil
}

func InitOpts(opts ...option) error {
	srv.Lock()
	defer srv.Unlock()

	for _, opt := range opts {
		opt(srv)
	}

	if srv.typ == 0 {
		return errors.New("bad service type")
	}

	if srv.name == "" {
		return errors.New("bad service name")
	}

	if srv.timeout == 0 {
		srv.timeout = DEFAULT_TIMEOUT
	}

	if srv.gatewayDir == "" {
		srv.gatewayDir = DEFAULT_GATEWAY_DIR
	}
	if srv.gatewayDir[len(srv.gatewayDir)-1] != '/' {
		srv.gatewayDir += "/"
	}

	if srv.serviceDir == "" {
		srv.serviceDir = DEFAULT_SERVICE_DIR
	}
	if srv.serviceDir[len(srv.serviceDir)-1] != '/' {
		srv.serviceDir += "/"
	}

	if srv.discoverMode == "" {
		srv.discoverMode = DEFAULT_DISCOVER_MODE
	}

	srv.ip = u.GetLocalAddress(srv.proberAddr)
	if srv.ip == "" {
		return errors.New("get local address failed")
	}

	srv.id = utils.Ip2long(srv.ip)
	srv.weight = utils.Atoi(os.Getenv(DEFAULT_WEIGHT_ENV))

	var err error
	if strings.ToUpper(srv.discoverMode) == DEFAULT_DISCOVER_MODE {
		if srv.etcdConf == nil || len(srv.etcdConf.Addrs) == 0 {
			return errors.New("bad ETCD config")
		}

		if srv.etcdConn, err = etcd.New(srv.etcdConf.Addrs, srv.etcdConf.User, srv.etcdConf.Pass); err != nil {
			return err
		}
	}

	if srv.kafkaConf == nil || len(srv.kafkaConf.Brokers) == 0 {
		return errors.New("bad KAFKA config")
	}

	onSuccess := func(msg *sarama.ProducerMessage) {
		value, err := msg.Value.Encode()
		if err != nil {
			logs.Errorf("Send kafka msg err: %s", err.Error())
		} else {
			logs.Tracef("Send kafka msg: %s", string(value))
		}
	}

	onError := func(err *sarama.ProducerError) {
		logs.Errorf("Send kafka msg err: %v", err.Error())
	}

	if srv.kafkaProducer, err = kafka.NewAsyncProducer(srv.kafkaConf.Brokers, 10*time.Second, onSuccess, onError); err != nil {
		return err
	}

	topic := srv.kafkaConf.TracerTopic
	fp := func(msg []byte, key string) error {
		return srv.kafkaProducer.SendMessage(topic, msg, key)
	}

	t.Init(srv.typ, srv.id, srv.name, srv.traceRate, fp)

	return nil
}

func Update(rate int) {
	t.Update(rate)
}

func Start() error {
	if strings.ToUpper(srv.DiscoverMode()) == DEFAULT_DISCOVER_MODE {
		if srv.EtcdConn() == nil {
			return errors.New("please init first")
		}

		if err := publish(); err != nil {
			return err
		}

		return fetchGateway()
	}

	return dialGateway(srv.GatewayAddr())
}

func Stop() {
	if srv.EtcdConn() == nil {
		return
	}

	revoke()
	srv.EtcdConn().Close()

	return
}

func register(so *p.Socket) error {
	regReq := &RegMsg{}
	regReq.Id = srv.Id()
	regReq.Type = uint32(srv.Type())
	regReq.Weight = int32(srv.Weight())
	regReq.Name = srv.Name()
	regReq.Secret = srv.Secret()

	body, _ := proto.Marshal(regReq)

	reqMsg := p.NewRequestMessage()
	reqMsg.SetToSvrType(ST_GW_CORE)
	reqMsg.SetFromSvrID(srv.Id())
	reqMsg.SetFromSvrType(srv.Type())
	reqMsg.SetFunctionID(F_ID_REGISTER)
	reqMsg.SetBody(body)

	rspMsg, err := so.Send(nil, reqMsg)
	if err != nil {
		return err
	}

	regRsp := &RspMsg{}
	if err := proto.Unmarshal(rspMsg.GetBody(), regRsp); err != nil {
		return err
	}

	if regRsp.Code != p.RC_OK {
		return errors.New(regRsp.Msg)
	}

	logs.Infof("Register service success")
	return nil
}

// func update(so *p.Socket) error {
// 	upReq := &UpdateMsg{}
// 	upReq.Id = srv.Id()
// 	upReq.Type = uint32(srv.Type())
// 	upReq.Weight = int32(srv.Weight())

// 	body, _ := proto.Marshal(upReq)

// 	reqMsg := p.NewRequestMessage()
// 	reqMsg.SetToSvrType(ST_GW_CORE)
// 	reqMsg.SetFromSvrID(srv.Id())
// 	reqMsg.SetFromSvrType(srv.Type())
// 	reqMsg.SetFunctionID(F_ID_UPDATE)
// 	reqMsg.SetBody(body)

// 	rspMsg, err := so.Send(nil, reqMsg)
// 	if err != nil {
// 		return err
// 	}

// 	upRsp := &RspMsg{}
// 	if err := proto.Unmarshal(rspMsg.GetBody(), upRsp); err != nil {
// 		return err
// 	}

// 	if upRsp.Code != p.RC_OK {
// 		return errors.New(upRsp.Msg)
// 	}

// 	return nil
// }

type _Publish struct {
	Id     uint32 `json:"id"`
	Name   string `json:"name"`
	Type   uint16 `json:"type"`
	Ip     string `json:"ip"`
	Weight int    `json:"weight"`
}

func publish() error {
	pub := &_Publish{}
	pub.Id = srv.Id()
	pub.Name = srv.Name()
	pub.Type = srv.Type()
	pub.Weight = srv.Weight()
	pub.Ip = srv.Ip()
	key := fmt.Sprintf("%vcore_%v", srv.ServiceDir(), srv.Id())

	body, _ := json.Marshal(pub)
	if err := srv.EtcdConn().SetKvAndKeepAlive(key, string(body), DEFAULT_ETCD_EXPIRE, DEFAULT_ETCD_PERIOD); err != nil {
		return err
	}

	// go watch()

	logs.Debugf("Publish to etcd: [key=%v,value=%v]", key, string(body))
	return nil
}

func revoke() error {
	key := fmt.Sprintf("%vcore_%v", srv.ServiceDir(), srv.Id())

	return srv.EtcdConn().Del(key)
}

// func watch() {
// 	key := fmt.Sprintf("%vcore_%v", srv.ServiceDir(), srv.Id())
// 	rch := srv.EtcdConn().WatchKeyWithPrefix(key)

// 	for wresp := range rch {
// 		if wresp.Canceled {
// 			logs.Infof("Watch %v canceled", key)
// 			return
// 		}

// 		if err := wresp.Err(); err != nil {
// 			logs.Errorf("Watch %v err: %s", key, err.Error())
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
// 					if err := update(conn); err != nil {
// 						logs.Errorf("Update weight err: %s", err.Error())
// 					}
// 				}
// 			case "DELETE":
// 				logs.Infof("Delete ETCD [key:%s,value:%s]", string(ev.Kv.Key), string(ev.Kv.Value))
// 			}
// 		}
// 	}
// }
