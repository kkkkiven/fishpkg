package http

import (
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kkkkiven/fishpkg/logs"

	"github.com/pkg/errors"
)

const (
	DEFAULT_DAEDLINE  = 10
	DEFAULT_PACK_SIZE = 1024 * 4
)

// 通知回调
type Notify interface {
	// 关闭事件
	OnClose(*Socket)

	// 超时事件
	OnTimeout(*Socket)
}

// 过滤器回调
type Filter interface {
	// 裸消息过滤
	OnRecv(*Socket, []byte) bool

	// 组包消息过滤
	OnMessage(*Socket, *Message) bool
}

type Socket struct {
	id            uint64
	conn          net.Conn
	mu            sync.RWMutex
	timeout       time.Duration
	maxPack       int
	lastWriteTime int64
	context       interface{}

	notify     Notify
	filter     Filter
	msgHandler func(*Socket, *Message)

	requestID     uint32
	requestQueue  map[uint32]chan *Message
	requestLocker sync.Locker
}

var socketID uint64

type option func(so *Socket)

func NewSocket(conn net.Conn, opts ...option) *Socket {
	so := &Socket{}
	so.conn = conn

	for _, opt := range opts {
		opt(so)
	}

	if so.id == 0 {
		so.id = atomic.AddUint64(&socketID, 1)
	}

	if so.timeout == 0 {
		so.timeout = DEFAULT_DAEDLINE
	}

	if so.maxPack == 0 {
		so.maxPack = DEFAULT_PACK_SIZE
	}

	so.requestQueue = make(map[uint32]chan *Message, 0)
	so.requestLocker = new(sync.Mutex)

	return so
}

func SetID(id uint64) option {
	return func(so *Socket) {
		so.id = id
	}
}

func SetTimeout(tm int64) option {
	return func(so *Socket) {
		so.timeout = time.Duration(tm)
	}
}

func SetMaxPack(size int) option {
	return func(so *Socket) {
		so.maxPack = size
	}
}

func SetNotify(notify Notify) option {
	return func(so *Socket) {
		so.notify = notify
	}
}

func SetFilter(filter Filter) option {
	return func(so *Socket) {
		so.filter = filter
	}
}

func SetContext(ctx interface{}) option {
	return func(so *Socket) {
		so.context = ctx
	}
}

func SetMsgHandler(handler func(*Socket, *Message)) option {
	return func(so *Socket) {
		so.msgHandler = handler
	}
}

func (this *Socket) GetID() uint64 {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return this.id
}

func (this *Socket) SetID(id uint64) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.id = id
}

// SetContext 设置上下文
func (this *Socket) SetContext(ctx interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.context = ctx
}

// GetContext 获取上下文
func (this *Socket) GetContext() interface{} {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return this.context
}

// GetConn 获取句柄
func (this *Socket) GetConn() net.Conn {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return this.conn
}

// SetConn 设置句柄
func (this *Socket) SetConn(conn net.Conn) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.conn = conn
}

// SetTimeout 设置超时
func (this *Socket) SetTimeout(tmout int64) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.timeout = time.Duration(tmout)
}

// GetTimeout 获取超时
func (this *Socket) GetTimeout() int64 {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return int64(this.timeout)
}

// SetMaxPack 设置读缓存
func (this *Socket) SetMaxPack(size int) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.maxPack = size
}

// GetMaxPack 获取读缓存
func (this *Socket) GetMaxPack() int {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return this.maxPack
}

func (this *Socket) SetNotify(notify Notify) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.notify = notify
}

func (this *Socket) SetFilter(filter Filter) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.filter = filter
}

// SetMsgHandler 设置通用消息处理函数
func (this *Socket) SetMsgHandler(handler func(*Socket, *Message)) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.msgHandler = handler
}

// Start 启动
func (this *Socket) Start() {
	go this.Run()
}

// Run 运行
func (this *Socket) Run() {
	defer func() {
		if err := recover(); err != nil {
			logs.Errorf("- %v - Panic:%+v,stack:%s",
				this.GetConn().RemoteAddr().String(),
				err,
				string(debug.Stack()))
		}
	}()

	this.readLoop()
}

// readLoop 读协程
func (this *Socket) readLoop() {
	defer func() {
		if this.notify != nil {
			this.notify.OnClose(this)
		}
		this.Close()
	}()

	// buf := make([]byte, 0, this.maxPack)
	var buf []byte
	tmp := make([]byte, this.maxPack)
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
				logs.Errorf("- %v - Read err: %v", this.conn.RemoteAddr().String(), err.Error())
				break
			}
		}

		buf = append(buf, rd[:size]...)

		logs.Tracef("- %v - ->> READ(%v bytes): %v", this.conn.RemoteAddr().String(), size, rd[:size])

		// 基于文本协议的消息处理 或者 消息透传
		if this.filter != nil {
			if ok := this.filter.OnRecv(this, buf); ok {
				continue
			}
		}

		// 基于Length-Type-Value协议的消息处理
		for {
			var msg *Message
			buf, msg = Decode(buf)
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
	if MT_RESPONSE == msg.GetMessageType() || MT_PONG == msg.GetMessageType() {
		this.requestLocker.Lock()
		ch, ok := this.requestQueue[msg.GetRequestID()]
		this.requestLocker.Unlock()
		if ok {
			ch <- msg
		}
		return
	}

	if this.msgHandler != nil {
		this.msgHandler(this, msg)
		return
	}

	if err := routerHandler(this, msg); err != nil {
		logs.Errorf("- %v - Dispatch err: %v", this.conn.RemoteAddr().String(), err.Error())
	}

	return
}

// Post 发送数据，不等待响应
func (this *Socket) Post(data []byte) error {
	return this.post(data)
}

// post 发送数据
func (this *Socket) post(data []byte) error {
	this.mu.Lock()
	defer this.mu.Unlock()

	length := len(data)
	offset := 0
	for {
		tm := time.Now().Add(time.Second * this.timeout)
		this.conn.SetWriteDeadline(tm)
		size, err := this.conn.Write(data[offset:])
		if err != nil {
			if err, ok := err.(*net.OpError); ok && err.Timeout() {
				// this.checkTimeout()
				logs.Errorf("- %v - Send timeout", this.conn.RemoteAddr().String())
				return errors.WithMessage(err, "timeout")
			} else {
				logs.Errorf("- %v - Send err: %v", this.conn.RemoteAddr().String(), err.Error())
				return errors.WithMessage(err, "post failed")
			}
		}

		this.lastWriteTime = time.Now().Unix()

		logs.Tracef("- %v - <<- SEND(%v bytes): %v", this.conn.RemoteAddr().String(), size, data[offset:offset+size])

		offset += size

		if offset >= length {
			return nil
		}
	}
}

// Send 发送数据，超时等待响应
func (this *Socket) Send(msg *Message) (*Message, error) {
	return this.sendTimeout(msg, int64(this.timeout))
}

// SendTimeout 发送数据，超时等待响应
func (this *Socket) sendTimeout(msg *Message, tmout int64) (rsp *Message, err error) {
	if msg.GetMessageType() == MT_RESPONSE || msg.GetMessageType() == MT_NORMAL || msg.GetMessageType() == MT_PONG {
		err = this.post(msg.Encode())
		return
	}

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

	if err = this.post(msg.Encode()); err != nil {
		return
	}

	select {
	case rsp = <-waitCh:
	case <-time.After(time.Duration(tmout) * time.Second):
		err = errors.New("send timeout")
	}

	return
}

// Close 关闭连接
func (this *Socket) Close() {
	this.conn.Close()
}
