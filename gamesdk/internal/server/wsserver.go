package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"git.yuetanggame.com/zfish/fishpkg/gamesdk/internal/utils"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/config"
	errs "git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/errors"
	"git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/game"
	"git.yuetanggame.com/zfish/fishpkg/logs"
	sdk "git.yuetanggame.com/zfish/fishpkg/servicesdk/core"
	"git.yuetanggame.com/zfish/fishpkg/sprotocol/core"
	jutils "git.yuetanggame.com/zfish/fishpkg/utils"
	"golang.org/x/net/websocket"
)

// StartWebSocket 初始化Websocket
// timeout 单位s
func StartWebSocket(port uint16) error {
	ch := make(chan error, 0)
	go func(ch chan<- error) {
		defer func() {
			if err := recover(); err != nil {
				logs.Error(debug.Stack())
			}
		}()

		http.Handle("/", websocket.Handler(serveWS))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			ch <- err
		}
	}(ch)

	select {
	case err := <-ch:
		return err
	case <-time.After(time.Second):
		close(ch)
		return nil
	}
}

// serveWS websocket连接处理
func serveWS(rawConn *websocket.Conn) {
	addr := rawConn.Request().RemoteAddr
	remoteIP := "0.0.0.0"
	var uid int64
	var err error

	if addr != "" {
		remoteIP = addr[0:strings.LastIndex(addr, ":")]
	}
	if rawConn.Request().Header.Get("X-Real-IP") != "" {
		remoteIP = rawConn.Request().Header.Get("X-Real-IP")
	}

	params := rawConn.Request().URL.Query()

	// 获取头部并签名校验
	cliSign := params.Get("X-SIGN")
	if cliSign == "" {
		logs.Errorf("client %s accept failed: no sign", remoteIP)
		_ = rawConn.WriteClose(int(errs.ErrCode(errs.ErrAuthorized)))
		return
	}

	uid, err = strconv.ParseInt(params.Get("X-UID"), 10, 64)
	if err != nil || uid == 0 {
		logs.Errorf("client %s accept failed: uid required", remoteIP)
		_ = rawConn.WriteClose(int(errs.ErrCode(errs.ErrAuthorized)))
		return
	}

	// 设置需要签名校验，则进行校验
	if config.Authorize() {
		stamp, _ := strconv.ParseInt(params.Get("X-TS"), 10, 64)

		// 请求时间必须在30分钟内，防止重放攻击
		ts := time.Unix(stamp, 0)
		if time.Since(ts) > 30*time.Minute || time.Since(ts) < -30*time.Minute {
			logs.Errorf("client %s accept failed: request expires", remoteIP)
			_ = rawConn.WriteClose(int(errs.ErrCode(errs.ErrAuthorized)))
			return
		}

		var buf bytes.Buffer
		buf.WriteString(params.Get("X-NONCE"))
		buf.WriteString(params.Get("X-TS"))
		buf.WriteString(params.Get("X-UID"))

		var info map[string]string
		info, err = sdk.Client().GetUserAttr(context.Background(), uid, []string{"token"}, nil)
		if err != nil {
			logs.Errorf("client %s accept failed:%s")
			_ = rawConn.WriteClose(int(errs.ErrCode(errs.ErrAuthorized)))
			return
		}

		buf.WriteString(info["token"])

		sign := utils.Sign(buf.String())

		if sign != cliSign {
			logs.Errorf("client %s accept failed: sign error", remoteIP)
			_ = rawConn.WriteClose(int(errs.ErrCode(errs.ErrAuthorized)))
			return
		}
	}

	rawConn.PayloadType = 2
	// 创建自定义socket对象

	socket := connPools.Get().(*core.Socket)
	socket.SetConn(rawConn)
	socket.SetRemoteIP(jutils.Ip2long(remoteIP))
	socket.SetIsWebsocket(true)

	reConn := false // 重连标记
	player := &game.Player{ID: uid, Conn: socket}

	// 签名ok
	// 若玩家已在游戏中，但连接不同，则先断开旧连接
	if p, err := game.GetPlayer(uid); err == nil {
		player.SetContext(p.GetContext())
		reConn = true

		originConn := p.Conn
		p.Conn = socket
		_ = game.UpdatePlayer(uid, p)

		if originConn != nil {
			// 把Context设成nil，防止触发player的 Lost
			originConn.SetContext(nil)
			originConn.SetFilter(GetServer())
			originConn.Close()
			logs.Waringf("player(%v) is already in the game, original connection has be closed", p)
		}
	} else if err := game.AddPlayer(player); err != nil {
		logs.Errorf("add player(%v) failed %v", p, err)
		_ = rawConn.WriteClose(int(errs.ErrCode(errs.DeskFull)))
		return
	}

	socket.SetContext(player)
	// 回调游戏
	if err := game.OnConnect(player, reConn); err != nil {
		game.RemovePlayer(uid)
	}
	logs.Infof("player(%v) authorize success", player)

	addConn(socket)
	// 启动读写
	socket.Start()
}
