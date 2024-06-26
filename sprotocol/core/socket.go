package core

import (
	"context"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.yuetanggame.com/zfish/fishpkg/logs"
	pb "git.yuetanggame.com/zfish/fishpkg/sprotocol/core/spropb"
	"git.yuetanggame.com/zfish/fishpkg/sprotocol/tracer"
	"git.yuetanggame.com/zfish/fishpkg/utils"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	DEFAULT_DAEDLINE      = 10
	DEFAULT_POST_BUF_SIZE = 1024 * 4
	DEFAULT_GATEWAY_TYPE  = 1
)

// 通知回调
type Notify interface {
	// 关闭事件
	OnClose(*Socket)

	// 超时事件
	OnTimeout(*Socket)
}

// 过滤回调
type Filter interface {
	// 裸消息过滤
	OnRecv(*Socket, []byte) bool

	// 组包消息过滤
	OnMessage(*Socket, *Message) bool
}

// Socket
type Socket struct {
	id            uint64
	conn          net.Conn
	mu            sync.RWMutex
	timeout       time.Duration
	lastWriteTime int64
	context       interface{}

	notify     Notify
	filter     Filter
	msgHandler func(context.Context, *Socket, *Message)

	isWebsocket bool
	remoteIP    uint32

	enablePack bool

	postBuf     []byte
	postBufSize int

	requestID     uint32
	requestQueue  map[uint32]chan *Message
	requestLocker sync.Locker
}

var socketID uint64

type option func(so *Socket)

// NewSocket
func NewSocket(conn net.Conn, opts ...option) *Socket {
	so := &Socket{}
	so.conn = conn

	so.enablePack = true
	so.postBufSize = DEFAULT_POST_BUF_SIZE

	for _, opt := range opts {
		opt(so)
	}

	if so.id == 0 {
		so.id = atomic.AddUint64(&socketID, 1)
	}

	if so.timeout == 0 {
		so.timeout = DEFAULT_DAEDLINE
	}

	so.requestQueue = make(map[uint32]chan *Message, 0)
	so.requestLocker = new(sync.Mutex)
	so.postBuf = make([]byte, so.postBufSize)

	return so
}

// SetID 设置连接id
func SetID(id uint64) option {
	return func(so *Socket) {
		so.id = id
	}
}

// SetEnablePack
func SetEnablePack(b bool) option {
	return func(so *Socket) {
		so.enablePack = b
	}
}

// SetPostBufSize
func SetPostBufSize(size int) option {
	return func(so *Socket) {
		so.postBufSize = size
	}
}

// SetTimeout 设置超时时间
func SetTimeout(tm int64) option {
	return func(so *Socket) {
		so.timeout = time.Duration(tm)
	}
}

// SetNotify
func SetNotify(notify Notify) option {
	return func(so *Socket) {
		so.notify = notify
	}
}

// SetFilter
func SetFilter(filter Filter) option {
	return func(so *Socket) {
		so.filter = filter
	}
}

// SetContext 设置上下文
func SetContext(ctx interface{}) option {
	return func(so *Socket) {
		so.context = ctx
	}
}

// SetMsgHandler 设置通用消息处理函数
func SetMsgHandler(handler func(context.Context, *Socket, *Message)) option {
	return func(so *Socket) {
		so.msgHandler = handler
	}
}

// SetRemoteIP 设置远程ip
func SetRemoteIP(ip uint32) option {
	return func(so *Socket) {
		so.remoteIP = ip
	}
}

// SetIsWebsocket 设置连接属性为websocket
func SetIsWebsocket(b bool) option {
	return func(so *Socket) {
		so.isWebsocket = b
	}
}

// GetID 获取socket id
func (this *Socket) GetID() uint64 {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return this.id
}

// SetID 设置socket id
func (this *Socket) SetID(id uint64) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.id = id
}

// SetContext 设置socket 上下文
func (this *Socket) SetContext(ctx interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.context = ctx
}

// GetContext 获取socket 上下文
func (this *Socket) GetContext() interface{} {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return this.context
}

// GetConn 获取conn对象
func (this *Socket) GetConn() net.Conn {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return this.conn
}

// SetConn 设置conn对象
func (this *Socket) SetConn(conn net.Conn) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.conn = conn
}

// SetTimeout 设置超时时间
func (this *Socket) SetTimeout(tmout int64) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.timeout = time.Duration(tmout)
}

// GetTimeout 获取超时时间
func (this *Socket) GetTimeout() int64 {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return int64(this.timeout)
}

// SetNotify 设置通知回调
func (this *Socket) SetNotify(notify Notify) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.notify = notify
}

// SetFilter 设置过滤回调
func (this *Socket) SetFilter(filter Filter) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.filter = filter
}

// SetMsgHandler 设置通用消息处理函数
func (this *Socket) SetMsgHandler(handler func(context.Context, *Socket, *Message)) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.msgHandler = handler
}

// SetRemoteIP 设置远程ip
func (this *Socket) SetRemoteIP(ip uint32) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.remoteIP = ip
}

// GetRemoteIP 获取远程ip
func (this *Socket) GetRemoteIP() uint32 {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if this.remoteIP != 0 {
		return this.remoteIP
	}

	addr := this.conn.RemoteAddr().String()
	this.remoteIP = utils.Ip2long(addr[0:strings.LastIndex(addr, ":")])

	return this.remoteIP
}

// GetRemoteIPStr 获取远程ip字符串
func (this *Socket) GetRemoteIPStr() string {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if this.remoteIP != 0 {
		return utils.Long2ip(this.remoteIP)
	}

	addr := this.conn.RemoteAddr().String()
	return addr[0:strings.LastIndex(addr, ":")]
}

// SetIsWebsocket 设置连接属性为websocket
func (this *Socket) SetIsWebsocket(b bool) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.isWebsocket = b
}

// GetIsWebsocket 获取连接属性为websocket
func (this *Socket) GetIsWebsocket() bool {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return this.isWebsocket
}

// Start 启动
func (this *Socket) Start() {
	this.Run()
}

func (this *Socket) Run() {
	defer func() {
		if err := recover(); err != nil {
			logs.Errorf("- %v - Panic:%v\n%s", this.GetConn().RemoteAddr().String(), err, string(debug.Stack()))
		}
	}()

	this.readLoop()
}

// readLoop 读循环
func (this *Socket) readLoop() {
	defer func() {
		if this.notify != nil {
			this.notify.OnClose(this)
		}
		this.Close()
	}()

	var buf []byte
	tmp := make([]byte, 1024)

	for {
		rd := tmp[:]

		tm := time.Now().Add(time.Second * this.timeout)
		this.conn.SetReadDeadline(tm)
		size, err := this.conn.Read(rd)

		if size <= 0 && err != nil {
			if err, ok := err.(*net.OpError); ok && err.Timeout() {
				this.checkTimeout()
				continue
			} else {
				logs.Errorf("- %v - Read error: %v", this.conn.RemoteAddr().String(), err.Error())
				break
			}
		}

		if !this.isWebsocket && this.enablePack {
			UnPack(rd[:size])
			// logs.Tracef("- %v - ->> READ after unpack(%v bytes): %v", this.conn.RemoteAddr().String(), size, rd[:size])
		}

		logs.Tracef("- %v - ->> READ(%v bytes): %v", this.conn.RemoteAddr().String(), size, rd[:size])

		buf = append(buf, rd[:size]...)

		// 基于文本协议的消息处理 或者 消息透传
		if this.filter != nil {
			if ok := this.filter.OnRecv(this, buf); ok {
				continue
			}
		}

		// 基于Length-Type-Value协议的消息处理
		for {
			var (
				msg *Message
				err error
			)
			buf, msg, err = Decode(buf[:])
			if err != nil {
				logs.Errorf("- %v - Decode message err: %s", this.conn.RemoteAddr().String(), err.Error())
				return
			}
			if msg == nil {
				break
			}

			if this.filter != nil {
				if this.filter.OnMessage(this, msg) {
					continue
				}
			}

			this.dispatch(msg)
		}
	}
}

// checkTimeout 检查超时
func (this *Socket) checkTimeout() {
	this.mu.RLock()
	if (time.Now().Unix() - this.lastWriteTime) < int64(this.timeout) {
		this.mu.RUnlock()
		return
	}
	this.mu.RUnlock()

	// 调用超时回调函数
	if nil != this.notify {
		this.notify.OnTimeout(this)
	}
}

// dispatch 消息分发
func (this *Socket) dispatch(msg *Message) {
	if MT_RESPONSE == msg.GetMessageType() {
		this.requestLocker.Lock()
		if ch, ok := this.requestQueue[msg.GetRequestID()]; ok {
			ch <- msg
		}
		this.requestLocker.Unlock()

		return
	}

	var (
		span *tracer.Span
		ctx  context.Context
	)
	if (msg.GetMessageFlag() & MF_TRACE) != 0 {
		span, ctx = tracer.CreateSpan(int64(msg.GetTraceID()), int64(msg.GetSpanID()))
		span.SetRemoteEndpoint("", msg.GetFromSvrType(), msg.GetFromSvrID(), "", 0)
		span.Tag("funcId", msg.GetFunctionID())
	}

	if this.msgHandler != nil {
		this.msgHandler(ctx, this, msg)
		span.End()
		return
	}

	if err := routerHandler(ctx, this, msg); err != nil {
		logs.Errorf("- %v - Dispatch error: %v", this.conn.RemoteAddr().String(), err.Error())
	}

	return
}

// Post 发送数据，不等待响应
func (this *Socket) Post(data []byte) error {
	return this.post(data)
}

// post 底层发送数据
func (this *Socket) post(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	logs.Tracef("- %v - <<- SEND(%v bytes): %v", this.conn.RemoteAddr().String(), len(data), data[:])

	this.mu.Lock()
	defer this.mu.Unlock()

	offset := 0
	for {
		cnt := copy(this.postBuf, data[offset:])
		if !this.isWebsocket && this.enablePack {
			Pack(this.postBuf[:cnt])
		}

		pos := 0
		for {
			tm := time.Now().Add(time.Second * this.timeout)
			this.conn.SetWriteDeadline(tm)
			size, err := this.conn.Write(this.postBuf[pos:cnt])
			if err != nil {
				if err, ok := err.(*net.OpError); ok && err.Timeout() {
					// this.checkTimeout()
					logs.Errorf("- %v - Send timeout", this.conn.RemoteAddr().String())
					return errors.WithMessage(err, "timeout")
				} else {
					logs.Errorf("- %v - Send error: %v", this.conn.RemoteAddr().String(), err.Error())
					return errors.WithMessage(err, "post failed")
				}
			}

			this.lastWriteTime = time.Now().Unix()

			pos += size

			if pos >= cnt {
				break
			}
		}

		offset += cnt
		if offset >= len(data) {
			break
		}
	}

	return nil
}

// Send 发送数据，等待响应，并进行链路追踪
func (this *Socket) Send(ctx context.Context, msg *Message) (*Message, error) {
	return this.sendTimeout(ctx, msg, 10)
}

// sendTimeout 发送数据，超时等待响应，进行链路追踪
func (this *Socket) sendTimeout(ctx context.Context, msg *Message, tmout int64) (rsp *Message, err error) {
	if msg.GetMessageType() == MT_RESPONSE || msg.GetMessageType() == MT_NORMAL {
		err = this.post(msg.Encode())
		return
	}

	span, _ := tracer.CreateSubSpan(ctx)

	waitCh := make(chan *Message, 1)
	defer close(waitCh)

	this.requestLocker.Lock()
	this.requestID++
	reqID := this.requestID
	this.requestQueue[this.requestID] = waitCh
	this.requestLocker.Unlock()

	defer func() {
		this.requestLocker.Lock()
		delete(this.requestQueue, reqID)
		this.requestLocker.Unlock()
	}()

	msg.SetRequestID(reqID)

	if span != nil {
		span.SetRemoteEndpoint("", msg.GetToSvrType(), msg.GetToSvrID(), "", 0)
		span.Tag("funcId", msg.GetFunctionID())
		msg.SetTraceID(uint64(span.GetTraceID()))
		msg.SetSpanID(uint64(span.GetSpanID()))
	}

	if err = this.post(msg.Encode()); err != nil {
		if span != nil {
			span.Tag("code", RC_SYS_ERR)
			span.Tag("msg", err.Error())
			span.End()
		}

		return
	}

	select {
	case rsp = <-waitCh:
		if rsp.GetFromSvrType() == DEFAULT_GATEWAY_TYPE { // 网关返回的状态消息
			rspMsg := &pb.RspCommon{}
			if e := proto.Unmarshal(rsp.GetBody(), rspMsg); e == nil {
				if span != nil {
					span.Tag("code", rspMsg.Code)
					span.Tag("msg", rspMsg.Msg)
					span.End()
				}

				rsp, err = nil, errors.Errorf(rspMsg.Msg)
				return
			}
		}

		if span != nil {
			span.Tag("code", RC_OK)
			span.Tag("msg", M(RC_OK))
			span.End()
		}
		return
	case <-time.After(time.Duration(tmout) * time.Second):
		if span != nil {
			span.Tag("code", RC_TIMEOUT)
			span.Tag("msg", M(RC_TIMEOUT))
			span.End()
		}

		err = errors.Errorf("send timeout")
	}

	return
}

// Close 关闭
func (this *Socket) Close() {
	this.conn.Close()
}
