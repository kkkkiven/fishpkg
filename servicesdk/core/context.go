package core

import (
	"context"
	"errors"
	"fmt"

	pb "git.yuetanggame.com/zfish/fishpkg/servicesdk/core/pb/core"
	p "git.yuetanggame.com/zfish/fishpkg/sprotocol/core"

	"github.com/golang/protobuf/proto"
)

type SDKContext struct {
	ctx context.Context
	so  *p.Socket
	msg *p.Message
}

func NewSDKContext(ctx context.Context, so *p.Socket, msg *p.Message) *SDKContext {
	sctx := &SDKContext{
		ctx: ctx,
		so:  so,
		msg: msg,
	}

	return sctx
}

func (this *SDKContext) GetMsg() *p.Message {
	return this.msg
}

func (this *SDKContext) GetBody() []byte {
	return this.msg.GetBody()
}

func (this *SDKContext) GetContext() context.Context {
	return this.ctx
}

func (this *SDKContext) GetFromSvrID() uint32 {
	return this.msg.GetFromSvrID()
}

func (this *SDKContext) GetFromSvrType() uint16 {
	return this.msg.GetFromSvrType()
}

func (this *SDKContext) SendRequest(svrType uint16, svrID uint32, handlerID uint16, body []byte) ([]byte, error) {
	req := p.NewRequestMessage()
	req.SetToSvrType(svrType)
	req.SetToSvrID(svrID)
	req.SetFunctionID(handlerID)
	req.SetBody(body)

	rsp, err := this.so.Send(this.ctx, req)
	if err != nil {
		return nil, err
	}

	return rsp.GetBody(), nil
}

func (this *SDKContext) SendResponse(body []byte) error {
	rsp := p.NewResponseMessage()
	rsp.SetRequestID(this.msg.GetRequestID())
	rsp.SetBody(body)
	_, err := this.so.Send(this.ctx, rsp)
	return err
}

func (this *SDKContext) SendNormal(svrType uint16, svrID uint32, handlerID uint16, body []byte) error {
	req := p.NewNormalMessage()
	req.SetToSvrType(svrType)
	req.SetToSvrID(svrID)
	req.SetFunctionID(handlerID)
	req.SetBody(body)

	_, err := this.so.Send(this.ctx, req)
	return err
}

func (this *SDKContext) SendBroadcast(svrType uint16, handlerID uint16, body []byte) error {
	req := p.NewNormalMessage()
	req.SetToSvrType(svrType)
	req.SetToSvrID(0xFFFFFFFF)
	req.SetFunctionID(handlerID)
	req.SetBody(body)

	_, err := this.so.Send(this.ctx, req)

	return err
}

func SendRequest(ctx context.Context, svrType uint16, svrID uint32, handlerID uint16, body []byte) ([]byte, error) {
	so := gwList.Roll()
	if so == nil {
		return nil, fmt.Errorf("no gateway is available")
	}

	req := p.NewRequestMessage()
	req.SetToSvrType(svrType)
	req.SetToSvrID(svrID)
	req.SetFunctionID(handlerID)
	req.SetBody(body)

	rsp, err := so.Send(ctx, req)
	if err != nil {
		return nil, err
	}

	return rsp.GetBody(), nil
}

func SendNormal(ctx context.Context, svrType uint16, svrID uint32, handlerID uint16, body []byte) error {
	so := gwList.Roll()
	if so == nil {
		return fmt.Errorf("no gateway is available")
	}

	req := p.NewNormalMessage()
	req.SetToSvrType(svrType)
	req.SetToSvrID(svrID)
	req.SetFunctionID(handlerID)
	req.SetBody(body)

	_, err := so.Send(ctx, req)
	return err
}

func SendBroadcast(ctx context.Context, svrType uint16, handlerID uint16, body []byte) error {
	so := gwList.Roll()
	if so == nil {
		return fmt.Errorf("no gateway is available")
	}

	req := p.NewNormalMessage()
	req.SetToSvrType(svrType)
	req.SetToSvrID(0xFFFFFFFF)
	req.SetFunctionID(handlerID)
	req.SetBody(body)

	_, err := so.Send(ctx, req)
	return err
}

func AddHandler(id uint16, fn func(*SDKContext)) {
	fp := func(ctx context.Context, so *p.Socket, msg *p.Message) {
		sctx := NewSDKContext(ctx, so, msg)
		fn(sctx)
	}

	if err := p.AddHandler(id, fp); err != nil {
		panic(fmt.Sprintf("Add handler err: %v", err.Error()))
	}
}

// SendAliLog 发送阿里日志
func SendAliLog(store, topic string, contents map[string]string) error {
	if store == "" {
		return errors.New("store can't be empty")
	}
	if topic == "" {
		return errors.New("topic can't be empty")
	}

	producer := srv.KafkaProducer()
	if producer == nil {
		return errors.New("please init first")
	}

	reqMsg := &pb.SLSMsg{
		Store:    store,
		Topic:    topic,
		Contents: contents,
	}
	msg, _ := proto.Marshal(reqMsg)

	if err := producer.SendMessage(srv.AliLogTopic(), msg); err != nil {
		return err
	}

	return nil
}

func SendBroadcastFish(ctx context.Context, svrType uint16, handlerID uint16, msgID uint32, body []byte) error {
	so := gwList.Roll()
	if so == nil {
		return fmt.Errorf("no gateway is available")
	}

	req := p.NewNormalMessage()
	req.SetToSvrType(svrType)
	req.SetToSvrID(0xFFFFFFFF)
	req.SetFunctionID(handlerID)
	req.SetBody(body)
	req.SetRequestID(msgID)

	_, err := so.Send(ctx, req)
	return err
}
