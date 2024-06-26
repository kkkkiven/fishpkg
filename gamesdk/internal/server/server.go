package server

import (
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/kkkkiven/fishpkg/gamesdk/internal/message"
	"github.com/kkkkiven/fishpkg/gamesdk/pkg/game"
	"github.com/kkkkiven/fishpkg/logs"
	tcp "github.com/kkkkiven/fishpkg/sprotocol/core"
	"golang.org/x/net/netutil"
)

type options struct {
	maxConns          int32         // 最大连接数
	connectionTimeOut time.Duration // 连接超时时长
}

var defaultOptions = options{
	maxConns:          100000,
	connectionTimeOut: 30, // 默认30s
}

// Option 设置Server参数
type Option func(*options)

// MaxConns 设置最大连接数
func MaxConns(n int32) Option {
	return func(o *options) {
		if n == 0 {
			o.maxConns = defaultOptions.maxConns

			return
		}
		o.maxConns = n
	}
}

// TimeOut 设置连接超时时长
// 默认30s
func TimeOut(d time.Duration) Option {
	return func(o *options) {
		if d == 0 {
			o.connectionTimeOut = defaultOptions.connectionTimeOut
			return
		}
		o.connectionTimeOut = d
	}
}

// Server ...
type Server struct {
	id       string
	opts     options
	mux      sync.Mutex
	listener net.Listener
	quit     chan struct{}
	done     chan struct{}
	cond     *sync.Cond     // conn close的信号量，热重启时使用
	wg       sync.WaitGroup // handle goroutines计数
	quitOnce sync.Once
	doneOnce sync.Once
}

var server = &Server{
	quit: make(chan struct{}),
	done: make(chan struct{}),
}

func GetServer() *Server {
	return server
}

var msgMap = make(map[uint16]Handler)

var connPools = &sync.Pool{
	New: func() interface{} {
		return tcp.NewSocket(nil, tcp.SetNotify(server),
			tcp.SetFilter(server),
			tcp.SetTimeout(int64(server.opts.connectionTimeOut)),
			tcp.SetEnablePack(true),
		)
	},
}

// 所有客户端连接,包括websocket
var conns = sync.Map{}

// Start ...
func Start(id string, listen net.Listener, opt ...Option) error {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	server.id = id
	server.opts = opts
	server.cond = sync.NewCond(&server.mux)

	server.mux.Lock()

	server.wg.Add(1)
	defer func() {
		server.wg.Done()
		select {
		// Stop or GracefulStop called; block until done and return nil.
		case <-server.quit:
			<-server.done
		default:
		}
	}()

	// 限制最大连接数
	server.listener = netutil.LimitListener(listen, int(server.opts.maxConns))
	server.mux.Unlock()

	var tempDelay time.Duration // Accept 失败时sleep的时长,最多1s
	for {
		// 接收到客户端请求
		rawConn, err := server.listener.Accept()
		if err != nil {
			if ne, ok := err.(interface {
				Temporary() bool
			}); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				logs.Waring(fmt.Sprintf("Accept error: %v; retrying in %v", err, tempDelay))

				timer := time.NewTimer(tempDelay)
				select {
				case <-timer.C:
				case <-server.quit:
					timer.Stop()
					return nil
				}
				continue
			}

			select {
			case <-server.quit:
				return nil
			default:
			}
			return err
		}

		tempDelay = 0

		socket := connPools.Get().(*tcp.Socket)
		socket.SetConn(rawConn)
		server.wg.Add(1)
		go func(socket *tcp.Socket) {
			defer func() {
				server.wg.Done()

				if err := recover(); err != nil {
					logs.Error(debug.Stack())
				}
			}()

			addConn(socket)
			socket.Start()
		}(socket)
	}
}

// OnClose ...
func (s *Server) OnClose(conn *tcp.Socket) {

	// 回调游戏处理
	ctx := conn.GetContext()
	if ctx != nil {
		player, _ := ctx.(*game.Player)
		game.OnPlayerLost(player)
		player.Conn = nil
	}
	removeConn(conn)
}

// OnTimeout ...
func (s *Server) OnTimeout(conn *tcp.Socket) {
	ctx := conn.GetContext()
	var player *game.Player
	if ctx != nil {
		player, _ = ctx.(*game.Player)
	}

	logs.Errorf("player(%v) time out", player)
	conn.SetFilter(s)
	conn.Close()
}

// OnRecv ...
func (s *Server) OnRecv(conn *tcp.Socket, data []byte) bool {
	return false
}

// OnMessage ...
func (s *Server) OnMessage(conn *tcp.Socket, msg *tcp.Message) bool {
	msgID := msg.GetFunctionID()
	ctx := conn.GetContext()

	// 如果启用了连接验证
	if msgID != message.MSGAuthorize {
		// 连接需要签名验证通过，才能发其他数据包，否则连接将被关闭
		if ctx == nil {
			logs.Errorf("close client %s: please authorized first", conn.GetRemoteIPStr())
			conn.SetFilter(s)
			conn.Close()
			return true
		}
	}

	switch msgID {
	case message.MSGAuthorize,
		message.MSGHeartbeat:
		return false
	default: // 其他消息
		player := ctx.(*game.Player)
		if handle, ok := msgMap[msgID]; ok {
			handle(player, msg.GetRequestID(), msg.GetBody())
			return true
		}
		game.OnMessage(player, msg.GetRequestID(), msgID, msg.GetBody())
		return true
	}
}

// 获取连接数
// s 普通socket个数
// ws websocket 个数
func GetConnCount() (s int32, ws int32) {
	conns.Range(func(key, value interface{}) bool {
		conn, _ := key.(*tcp.Socket)
		if conn != nil {
			if conn.GetIsWebsocket() {
				ws++
				return true
			}
			s++
		}
		return true
	})

	return s, ws
}

func GetServerID() string {
	return server.id
}

func addConn(conn *tcp.Socket) {
	conns.Store(conn, struct{}{})
}

func removeConn(conn *tcp.Socket) {
	conns.Delete(conn)
	server.cond.Broadcast()
	ctx := conn.GetContext()
	if ctx != nil {
		if player, ok := ctx.(*game.Player); ok {
			player.Conn = nil
		}
	}
	conn.SetContext(nil)
	conn.SetRemoteIP(0)
	conn.SetIsWebsocket(false)
	conn.SetFilter(GetServer())
	// connPools.Put(conn)
}

// Stop 停止服务
func Stop() {
	server.quitOnce.Do(func() {
		close(server.quit)
	})

	defer func() {
		server.wg.Wait()
		server.doneOnce.Do(func() {
			close(server.done)
		})
	}()

	server.mux.Lock()

	listen := server.listener
	server.listener = nil

	server.cond.Broadcast()
	server.mux.Unlock()

	_ = listen.Close()

	conns.Range(func(k, _ interface{}) bool {
		if conn, ok := k.(net.Conn); ok {
			_ = conn.Close()
			conns.Delete(k)
		}
		return true
	})
}

// GracefulStop 服务优雅退出
// 会等待所有连接的处理goroutine结束，才返回
func GracefulStop() {
	server.quitOnce.Do(func() {
		close(server.quit)
	})

	defer func() {
		server.doneOnce.Do(func() {
			close(server.done)
		})
	}()

	server.mux.Lock()
	_ = server.listener.Close()
	server.listener = nil

	// 等待所有正在处理的goroutine完成
	server.mux.Unlock()
	server.wg.Wait()

	server.mux.Lock()

	conns.Range(func(k, _ interface{}) bool {
		conn := k.(net.Conn)
		_ = conn.Close()
		conns.Delete(k)
		return true
	})

	for {
		count := 0
		conns.Range(func(_, _ interface{}) bool {
			count++
			return true
		})

		if count != 0 {
			server.cond.Wait()
		} else {
			break
		}
	}

	server.mux.Unlock()
}

// AddHandler ...
func AddHandler(msgID uint16, handle Handler) {
	msgMap[msgID] = handle
}

type Handler func(player *game.Player, requestID /*请求包ID*/ uint32, msg []byte)
