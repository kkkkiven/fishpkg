package tracer

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"git.yuetanggame.com/zfish/fishpkg/logs"
	"git.yuetanggame.com/zfish/fishpkg/utils"

	"github.com/beinan/fastid"
)

const (
	KIND_CLIENT = "CLIENT"
	KIND_SERVER = "SERVER"
)

const (
	DEFAULT_TRACER_RATE = 5000
)

type _Manager struct {
	sync.RWMutex

	open       bool
	rate       int
	endpoint   *Endpoint
	aggregator func([]byte, string) error
}

var m *_Manager = &_Manager{rate: DEFAULT_TRACER_RATE}

func Init(serviceType uint16, serviceId uint32, serviceName string, rate int, fp func([]byte, string) error) {
	m.Lock()
	defer m.Unlock()

	m.rate = rate
	if rate < 0 {
		m.rate = DEFAULT_TRACER_RATE
	}

	if m.rate > 0 {
		m.open = true
	}

	m.endpoint = &Endpoint{}
	m.endpoint.ServiceType = serviceType
	m.endpoint.ServiceName = serviceName
	m.endpoint.ServiceId = serviceId
	m.endpoint.IPv4 = utils.Long2ip(serviceId)

	m.aggregator = fp
}

func Update(rate int) {
	m.Lock()
	defer m.Unlock()

	if rate < 0 {
		return
	}

	m.rate = rate

	if m.rate == 0 {
		m.open = false
	}

	return
}

type Span struct {
	TraceId        int64                  `json:"-"`
	TraceIdStr     string                 `json:"traceId"`
	ParentId       int64                  `json:"-"`
	ParentIdStr    string                 `json:"parentId,omitempty"`
	SpanId         int64                  `json:"-"`
	SpanIdStr      string                 `json:"id"`
	Name           string                 `json:"name"`
	Kind           string                 `json:"kind"`
	Timestamp      int64                  `json:"timestamp"`
	Duration       int64                  `json:"duration"`
	LocalEndpoint  *Endpoint              `json:"localEndpoint,omitempty"`
	RemoteEndPoint *Endpoint              `json:"remoteEndpoint,omitempty"`
	Annotations    []Annotation           `json:"annotations,omitempty"`
	Tags           map[string]interface{} `json:"tags,omitempty"`
}

type PropagateSpan struct {
	traceId int64
	id      int64
}

type Endpoint struct {
	ServiceName string `json:"serviceName"`
	ServiceType uint16 `json:"serviceType,omitempty"`
	ServiceId   uint32 `json:"serviceId,omitempty"`
	IPv4        string `json:"ipv4,omitempty"`
	Port        uint16 `json:"port,omitempty"`
}

type Annotation struct {
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

// CreateSpan 创建span, kind为SERVER
func CreateSpan(traceId, parentId int64) (*Span, context.Context) {
	m.RLock()
	defer m.RUnlock()

	return createSpan(traceId, parentId)
}

// CreatePropSpan 根据概率创建span，kind为SERVER
func CreateProbSpan(traceId, parentId int64) (*Span, context.Context) {
	m.RLock()
	defer m.RUnlock()

	if !m.open {
		return nil, context.WithValue(context.TODO(), ctxNoopKeyInstance, &ctxNoopKey{})
	}

	if rand.Intn(m.rate) != 0 {
		return nil, context.WithValue(context.TODO(), ctxNoopKeyInstance, &ctxNoopKey{})
	}

	return createSpan(traceId, parentId)
}

func createSpan(traceId, parentId int64) (*Span, context.Context) {
	if !m.open {
		return nil, context.WithValue(context.TODO(), ctxNoopKeyInstance, &ctxNoopKey{})
	}

	span := new(Span)
	span.TraceId = traceId
	span.ParentId = parentId
	span.SpanId = fastid.CommonConfig.GenInt64ID()

	span.Kind = KIND_SERVER
	span.Name = m.endpoint.ServiceName
	span.LocalEndpoint = m.endpoint
	span.Timestamp = time.Now().UnixNano() / 1e3

	return span, context.WithValue(context.TODO(), ctxKeyInstance, span)
}

// CreateSubSpan 创建span，kind为CLIENT
func CreateSubSpan(ctx context.Context) (*Span, context.Context) {
	if ctx == nil {
		return nil, context.WithValue(context.TODO(), ctxNoopKeyInstance, &ctxNoopKey{})
	}

	if _, ok := ctx.Value(ctxNoopKeyInstance).(*ctxNoopKey); ok {
		return nil, ctx
	}

	m.RLock()
	defer m.RUnlock()

	if !m.open {
		return nil, context.WithValue(context.TODO(), ctxNoopKeyInstance, &ctxNoopKey{})
	}

	span := &Span{}
	span.Kind = KIND_CLIENT
	span.Name = m.endpoint.ServiceName
	span.LocalEndpoint = m.endpoint
	span.Timestamp = time.Now().UnixNano() / 1e3
	span.SpanId = fastid.CommonConfig.GenInt64ID()

	// core服务父span创建子span
	parentSpan, ok := ctx.Value(ctxKeyInstance).(*Span)
	if ok {
		span.TraceId = parentSpan.TraceId
		span.ParentId = parentSpan.SpanId

		return span, context.WithValue(ctx, ctxKeyInstance, span)
	}

	// http服务父span创建子span
	propagateSpan, ok := ctx.Value(ctxPropagateKeyInstance).(*PropagateSpan)
	if ok {
		span.TraceId = propagateSpan.traceId
		span.ParentId = propagateSpan.id

		return span, context.WithValue(ctx, ctxKeyInstance, span)
	}

	// 创建根节点
	if rand.Intn(int(m.rate)) != 0 {
		return nil, context.WithValue(context.TODO(), ctxNoopKeyInstance, &ctxNoopKey{})
	}

	span.TraceId = fastid.CommonConfig.GenInt64ID()

	return span, context.WithValue(ctx, ctxKeyInstance, span)
}

func CreatePropagateSpan(traceId, parentId int64) (*PropagateSpan, context.Context) {
	m.RLock()
	defer m.RUnlock()

	if !m.open {
		return nil, context.WithValue(context.TODO(), ctxNoopKeyInstance, &ctxNoopKey{})
	}

	s := &PropagateSpan{
		traceId: traceId,
		id:      parentId,
	}

	return s, context.WithValue(context.TODO(), ctxPropagateKeyInstance, s)
}

func GetSpan(ctx context.Context) *Span {
	if ctx == nil {
		return nil
	}

	if _, ok := ctx.Value(ctxNoopKeyInstance).(*ctxNoopKey); ok {
		return nil
	}

	span, ok := ctx.Value(ctxKeyInstance).(*Span)
	if ok {
		return span
	}

	return nil
}

func (s *Span) GetTraceID() int64 {
	if s == nil {
		return 0
	}

	return s.TraceId
}

func (s *Span) GetSpanID() int64 {
	if s == nil {
		return 0
	}

	return s.SpanId
}

func (s *Span) Tag(key string, val interface{}) {
	if s == nil {
		return
	}

	if s.Tags == nil {
		s.Tags = make(map[string]interface{}, 0)
	}

	value := fmt.Sprintf("%v", val)

	s.Tags[key] = value
	return
}

func (s *Span) AddAnnotation(value string, timestamp int64) {
	if s == nil {
		return
	}

	an := Annotation{value, timestamp}
	s.Annotations = append(s.Annotations, an)
	return
}

func (s *Span) SetRemoteEndpoint(serviceName string, serviceType uint16, serviceId uint32, ipv4 string, port uint16) {
	if s == nil {
		return
	}

	if s.RemoteEndPoint == nil {
		s.RemoteEndPoint = new(Endpoint)
	}

	s.RemoteEndPoint.ServiceName = serviceName
	s.RemoteEndPoint.ServiceType = serviceType
	s.RemoteEndPoint.ServiceId = serviceId
	s.RemoteEndPoint.IPv4 = ipv4
	s.RemoteEndPoint.Port = port

	return
}

func (s *Span) SetLocalEndpoint(serviceName string, serviceType uint16, serviceId uint32, ipv4 string, port uint16) {
	if s == nil {
		return
	}

	if s.LocalEndpoint == nil {
		s.LocalEndpoint = new(Endpoint)
	}
	s.LocalEndpoint.ServiceName = serviceName
	s.LocalEndpoint.ServiceType = serviceType
	s.LocalEndpoint.ServiceId = serviceId
	s.LocalEndpoint.IPv4 = ipv4
	s.LocalEndpoint.Port = port

	return
}

func (s *Span) SetKind(k string) {
	if s == nil {
		return
	}

	s.Kind = k
}

func (s *Span) End() {
	if s == nil {
		return
	}

	var tids, pids, sids [8]byte
	binary.BigEndian.PutUint64(tids[:], uint64(s.TraceId))
	binary.BigEndian.PutUint64(pids[:], uint64(s.ParentId))
	binary.BigEndian.PutUint64(sids[:], uint64(s.SpanId))
	s.TraceIdStr = hex.EncodeToString(tids[:])
	s.ParentIdStr = hex.EncodeToString(pids[:])
	s.SpanIdStr = hex.EncodeToString(sids[:])
	s.Duration = time.Now().UnixNano()/1e3 - s.Timestamp

	if s.ParentId == 0 {
		s.ParentIdStr = ""
	}

	msg, _ := json.Marshal(s)

	if m.aggregator == nil {
		logs.Debug(string(msg))
	} else {
		if err := m.aggregator(msg, utils.I64toA(s.TraceId)); err != nil {
			logs.Error(err.Error())
		}
	}

	return
}

type ctxKey struct{}

var ctxKeyInstance = ctxKey{}

type ctxPropagateKey struct{}

var ctxPropagateKeyInstance = ctxPropagateKey{}

type ctxNoopKey struct{}

var ctxNoopKeyInstance = ctxNoopKey{}
