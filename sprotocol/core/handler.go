package core

import (
	"context"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"

	"git.yuetanggame.com/zfish/fishpkg/logs"
	pb "git.yuetanggame.com/zfish/fishpkg/sprotocol/core/spropb"
	"git.yuetanggame.com/zfish/fishpkg/sprotocol/tracer"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// Handler 消息处理函数
type Handler func(context.Context, *Socket, *Message)

// wrapper
type wrapper struct {
	handler Handler
	name    string // 函数名称，用于全链路跟踪展示被调函数名称
}

var handlerMap map[uint16]wrapper

func init() {
	handlerMap = make(map[uint16]wrapper)
}

func routerHandler(ctx context.Context, so *Socket, msg *Message) error {
	fid := msg.GetFunctionID()
	hw, ok := handlerMap[fid]
	if !ok {
		rspMsg := &pb.RspCommon{}
		rspMsg.Code = RC_HANDLER_NOT_FOUND
		rspMsg.Msg = M(RC_HANDLER_NOT_FOUND)
		content, _ := proto.Marshal(rspMsg)

		rsp := NewResponseMessage()
		rsp.SetRequestID(msg.GetRequestID())
		rsp.SetBody(content)

		so.Send(ctx, rsp)

		span := tracer.GetSpan(ctx)
		if span != nil {
			span.Tag("code", RC_HANDLER_NOT_FOUND)
			span.Tag("msg", M(RC_HANDLER_NOT_FOUND))
			span.End()
		}

		return errors.Errorf("handler[%v] not found", fid)
	}

	go doCall(ctx, so, msg, hw.handler)
	return nil
}

func doCall(ctx context.Context, so *Socket, msg *Message, handler Handler) {
	defer func() {
		if err := recover(); nil != err {
			rspMsg := &pb.RspCommon{}
			rspMsg.Code = RC_HANDLER_PANIC
			rspMsg.Msg = M(RC_HANDLER_PANIC)
			content, _ := proto.Marshal(rspMsg)

			rsp := NewResponseMessage()
			rsp.SetRequestID(msg.GetRequestID())
			rsp.SetBody(content)

			so.Send(ctx, rsp)

			span := tracer.GetSpan(ctx)
			if span != nil {
				span.Tag("code", RC_HANDLER_PANIC)
				span.Tag("msg", M(RC_HANDLER_PANIC))
				span.End()
			}

			logs.Errorf("- %v - Panic: invork handler[%v] err: %v\n%s", so.GetConn().RemoteAddr().String(), msg.GetFunctionID(), err, string(debug.Stack()))
		}
	}()

	h := chain(handler)
	h(ctx, so, msg)

	span := tracer.GetSpan(ctx)
	if span != nil {
		span.Tag("code", RC_OK)
		span.Tag("msg", M(RC_OK))
		span.End()
	}
}

// AddHandler 添加消息处理函数
func AddHandler(funcID uint16, handler Handler) error {
	if _, ok := handlerMap[funcID]; ok {
		return errors.Errorf("handler[%v] already exists", funcID)
	}

	nameFull := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	nameEnd := filepath.Ext(nameFull)
	funcName := strings.TrimPrefix(nameEnd, ".")

	hw := wrapper{
		handler: handler,
		name:    funcName,
	}

	handlerMap[funcID] = hw

	return nil
}

// getFuncName 获取函数id对应的函数名称
func getFuncName(id uint16) string {
	if hw, ok := handlerMap[id]; ok {
		return hw.name
	}
	return ""
}
