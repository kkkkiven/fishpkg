package core

import (
	"bytes"
	"errors"
	"io"
	"unsafe"

	"compress/zlib"
)

const (
	HL_FIX      = 0x08 // fix head length
	HL_REQUEST  = 0x0A // ref to request and normal message
	HL_ROUTE    = 0x0C // ref to route
	HL_TRACE    = 0x10 // ref to trace
	HL_RESPONSE = 0x08 // ref to response message

	MAX_BODY_LENGTH = 0xFFFF // max total length of first package
)

// message type
const (
	MT_NORMAL byte = iota
	MT_REQUEST
	MT_RESPONSE
)

// message flags
const (
	MF_RESERVE byte = 1 << iota
	MF_COMPRESS
	MF_ROUTER
	MF_TRACE
	MF_PACKAGE
)

type Message struct {
	context interface{} // context: store customer data

	msgLen  uint16 // message length of first package
	msgType byte   // message type: normal, request, response
	msgFlag byte   // flag: encode, compress, route, trace, package
	reqID   uint32 // request id

	funcID uint16 // function id, ref for request and normal message

	// route head
	fromSvrType uint16 // source server type
	fromSvrID   uint32 // source server id
	toSvrType   uint16 // remote server type
	toSvrID     uint32 // remote server id

	// trace head
	traceID uint64 // link trace id
	spanID  uint64 // link span id

	msgBody []byte
}

// NewMessage 消息构造函数
func NewMessage(msgType byte) *Message {
	msg := &Message{}
	msg.msgType = msgType
	return msg
}

// NewNormalMessage 创建普通消息
func NewNormalMessage() *Message {
	msg := NewMessage(MT_NORMAL)
	return msg
}

// NewRequestMessage 创建请求消息
func NewRequestMessage() *Message {
	msg := NewMessage(MT_REQUEST)
	return msg
}

// NewResponseMessage 创建响应消息
func NewResponseMessage() *Message {
	msg := NewMessage(MT_RESPONSE)
	return msg
}

// Decode 解码
func Decode(buf []byte) ([]byte, *Message, error) {
	pos := 0
	size := len(buf)

	if (size - pos) < HL_FIX { // buf长度小于固定头部长度，直接返回
		return buf, nil, nil
	}

	msgLen := get16bit(buf, pos)
	if size < int(msgLen) { // buf长度小于消息长度，直接返回
		return buf, nil, nil
	}

	msgType := get8bit(buf, pos+2)
	msgFlag := get8bit(buf, pos+3)
	reqID := get32bit(buf, pos+4)

	if msgFlag > 31 { // 掩码错误
		return nil, nil, errors.New("bad message flag")
	}

	var subLen uint32
	if (msgFlag & MF_PACKAGE) != 0 {
		if size < MAX_BODY_LENGTH+4 { // 存在附件包且buf长度不足无法解析附件包长度，直接返回
			return buf, nil, nil
		}
		subLen = get32bit(buf, pos+MAX_BODY_LENGTH)

		if size < int(MAX_BODY_LENGTH+4+subLen) { // 存在附件包且buf长度小于总长度，直接返回
			return buf, nil, nil
		}
	}

	msg := NewMessage(msgType)
	msg.msgType = msgType
	msg.msgFlag = msgFlag
	msg.reqID = reqID

	switch msg.msgType {
	case MT_NORMAL:
		msg.funcID = get16bit(buf, pos+8)
		pos += 10
	case MT_REQUEST:
		msg.funcID = get16bit(buf, pos+8)
		pos += 10
	case MT_RESPONSE:
		pos += 8
	default: // 消息类型不匹配，丢弃buf ？？？
		return nil, nil, errors.New("bad message type")
	}

	if (msg.msgFlag & MF_ROUTER) != 0 {
		msg.fromSvrType = get16bit(buf, pos)
		msg.toSvrType = get16bit(buf, pos+2)
		msg.fromSvrID = get32bit(buf, pos+4)
		msg.toSvrID = get32bit(buf, pos+8)
		pos += 12
	}

	if (msg.msgFlag & MF_TRACE) != 0 {
		msg.traceID = get64bit(buf, pos)
		msg.spanID = get64bit(buf, pos+8)
		pos += 16
	}

	if pos > int(msgLen) {
		return buf[pos:], nil, errors.New("bad message length")
	}

	var (
		err     error
		msgBody []byte
	)

	if subLen == 0 {
		msgBody = buf[pos:msgLen]
		pos = int(msgLen)
	} else {
		msgBody = buf[pos:msgLen]
		pos = int(msgLen) + 4
		msgBody = append(msgBody, buf[pos:pos+int(subLen)]...)
		pos += int(subLen)
	}

	if (msg.msgFlag & MF_COMPRESS) != 0 { // 解压缩
		msgBody, err = zlibUnCompress(msgBody)
		if err != nil {
			return buf[pos:], nil, err
		}
	}

	msg.msgBody = msgBody
	return buf[pos:], msg, nil
}

// Encode 编码
func (this *Message) Encode() []byte {
	var msgLen uint16

	if this.msgType == MT_RESPONSE {
		msgLen = HL_RESPONSE
	} else {
		msgLen = HL_REQUEST
	}

	if (this.msgFlag & MF_ROUTER) != 0 {
		msgLen += HL_ROUTE
	}

	if (this.msgFlag & MF_TRACE) != 0 {
		msgLen += HL_TRACE
	}

	var msgBody []byte
	if (this.msgFlag & MF_COMPRESS) != 0 /*&& len(this.msgBody) > 2048*/ { // 压缩
		msgBody = zlibCompress(this.msgBody)
	} else {
		// this.msgFlag &= ^MF_COMPRESS
		msgBody = this.msgBody
	}

	bodyLen := len(msgBody)
	allLen := int(msgLen) + bodyLen
	if allLen > MAX_BODY_LENGTH {
		this.msgFlag |= MF_PACKAGE
		msgLen = MAX_BODY_LENGTH
	} else {
		msgLen = uint16(allLen)
	}

	buf := make([]byte, 0, allLen)
	buf = put16bit(buf, msgLen)
	buf = put8bit(buf, this.msgType)
	buf = put8bit(buf, this.msgFlag)
	buf = put32bit(buf, this.reqID)

	switch this.msgType {
	case MT_NORMAL:
		fallthrough
	case MT_REQUEST:
		buf = put16bit(buf, this.funcID)
	case MT_RESPONSE:
	default:
		return nil
	}

	if (this.msgFlag & MF_ROUTER) != 0 {
		buf = put16bit(buf, this.fromSvrType)
		buf = put16bit(buf, this.toSvrType)
		buf = put32bit(buf, this.fromSvrID)
		buf = put32bit(buf, this.toSvrID)
	}

	if (this.msgFlag & MF_TRACE) != 0 {
		buf = put64bit(buf, this.traceID)
		buf = put64bit(buf, this.spanID)
	}

	if allLen > MAX_BODY_LENGTH { // 附件包
		fstLen := MAX_BODY_LENGTH + bodyLen - allLen
		sndLen := bodyLen - fstLen
		buf = append(buf, msgBody[:fstLen]...)
		buf = put32bit(buf, uint32(sndLen))
		buf = append(buf, msgBody[fstLen:]...)
	} else {
		buf = append(buf, msgBody...)
	}

	return buf
}

// zlibCompress 进行zlib压缩
func zlibCompress(src []byte) []byte {
	return zlibP.Compress(src)
}

// zlibUnCompress 进行zlib解压缩
func zlibUnCompress(src []byte) ([]byte, error) {
	b := bytes.NewReader(src)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	io.Copy(&out, r)
	return out.Bytes(), nil
}

// GetContext 获取消息context
func (this *Message) GetContext() interface{} {
	return this.context
}

// SetContext 设置消息context
func (this *Message) SetContext(ctx interface{}) {
	this.context = ctx
}

// GetMessageType 获取消息类型
func (this *Message) GetMessageType() byte {
	return this.msgType
}

// SetMessageType 设置消息类型
func (this *Message) SetMessageType(t byte) {
	this.msgType = t
}

// GetMessageFlag 获取消息标志
func (this *Message) GetMessageFlag() byte {
	return this.msgFlag
}

// SetMessageFlag 设置消息标志
func (this *Message) SetMessageFlag(flag byte) {
	this.msgFlag = flag
}

// SetRequestID 设置请求id
func (this *Message) SetRequestID(id uint32) {
	this.reqID = id
}

// GetRequestID 获取请求id
func (this *Message) GetRequestID() uint32 {
	return this.reqID
}

// GetFunctionID 获取接口id
func (this *Message) GetFunctionID() uint16 {
	return this.funcID
}

// SetFunctionID 设置接口id
func (this *Message) SetFunctionID(id uint16) {
	this.funcID = id
}

// GetToSvrType 获取目的服务类型
func (this *Message) GetToSvrType() uint16 {
	return this.toSvrType
}

// SetToSvrType 设置目的服务类型
func (this *Message) SetToSvrType(t uint16) {
	this.msgFlag |= MF_ROUTER
	this.toSvrType = t
}

// GetFromSvrID 获取源服务id
func (this *Message) GetFromSvrID() uint32 {
	return this.fromSvrID
}

// SetFromSvrID 设置源服务id
func (this *Message) SetFromSvrID(id uint32) {
	this.msgFlag |= MF_ROUTER
	this.fromSvrID = id
}

// GetFromSvrType 获取源服务类型
func (this *Message) GetFromSvrType() uint16 {
	return this.fromSvrType
}

// SetFromSvrType 设置源服务类型
func (this *Message) SetFromSvrType(t uint16) {
	this.msgFlag |= MF_ROUTER
	this.fromSvrType = t
}

// GetToSvrID 获取目的服务id
func (this *Message) GetToSvrID() uint32 {
	return this.toSvrID
}

// SetToSvrID 设置目的服务id
func (this *Message) SetToSvrID(id uint32) {
	this.msgFlag |= MF_ROUTER
	this.toSvrID = id
}

// GetTraceID 获取trace id
func (this *Message) GetTraceID() uint64 {
	return this.traceID
}

// SetTraceID 设置trace id
func (this *Message) SetTraceID(id uint64) {
	this.msgFlag |= MF_TRACE
	this.traceID = id
}

// GetSpanID 获取span id
func (this *Message) GetSpanID() uint64 {
	return this.spanID
}

// SetSpanID 设置span id
func (this *Message) SetSpanID(id uint64) {
	this.msgFlag |= MF_TRACE
	this.spanID = id
}

// GetBody 获取消息体
func (this *Message) GetBody() []byte {
	return this.msgBody
}

// SetBody 设置消息体
func (this *Message) SetBody(body []byte) {
	this.msgBody = body
	this.msgLen = uint16(len(this.msgBody))
}

// GetStringBody 获取消息体
func (this *Message) GetStringBody() string {
	if this.msgBody == nil {
		return ""
	}

	return *(*string)(unsafe.Pointer(&this.msgBody))
}

// SetStringBody 设置消息体
func (this *Message) SetStringBody(body string) {
	this.msgBody = *(*[]byte)(unsafe.Pointer(&body))
	this.msgLen = uint16(len(this.msgBody))
}

// Reset 重置消息
func (this *Message) Reset() {
	this.funcID = 0
	this.fromSvrType = 0
	this.fromSvrID = 0
	this.toSvrType = 0
	this.toSvrID = 0
	this.traceID = 0
	this.spanID = 0
	this.msgBody = nil
	this.msgLen = 0
}

// 加入8字节
func put8bit(buf []byte, n byte) []byte {
	return append(buf, n)
}

// 加入16字节
func put16bit(buf []byte, n uint16) []byte {
	var by [2]byte

	by[0] = byte((n >> 8) & 0xff)
	by[1] = byte(n & 0xff)

	return append(buf, by[:]...)
}

// 加入32字节
func put32bit(buf []byte, n uint32) []byte {
	var by [4]byte

	by[0] = byte((n >> 24) & 0xff)
	by[1] = byte((n >> 16) & 0xff)
	by[2] = byte((n >> 8) & 0xff)
	by[3] = byte(n & 0xff)

	return append(buf, by[:]...)
}

// 加入64字节
func put64bit(buf []byte, n uint64) []byte {
	var by [8]byte

	by[0] = byte((n >> 56) & 0xff)
	by[1] = byte((n >> 48) & 0xff)
	by[2] = byte((n >> 40) & 0xff)
	by[3] = byte((n >> 32) & 0xff)
	by[4] = byte((n >> 24) & 0xff)
	by[5] = byte((n >> 16) & 0xff)
	by[6] = byte((n >> 8) & 0xff)
	by[7] = byte(n & 0xff)

	return append(buf, by[:]...)
}

// 获取8bit
func get8bit(buf []byte, start int) byte {
	return buf[start]
}

// 获取16bit
func get16bit(buf []byte, start int) uint16 {
	var ret uint16

	ret = uint16(buf[start]) << 8
	ret |= uint16(buf[start+1])

	return ret
}

// 获取32big
func get32bit(buf []byte, start int) uint32 {
	var ret uint32

	ret = uint32(buf[start]) << 24
	ret |= uint32(buf[start+1]) << 16
	ret |= uint32(buf[start+2]) << 8
	ret |= uint32(buf[start+3])

	return ret
}

// 获取64bit
func get64bit(buf []byte, start int) uint64 {
	var ret uint64

	ret = uint64(buf[start]) << 56
	ret |= uint64(buf[start+1]) << 48
	ret |= uint64(buf[start+2]) << 40
	ret |= uint64(buf[start+3]) << 32
	ret |= uint64(buf[start+4]) << 24
	ret |= uint64(buf[start+5]) << 16
	ret |= uint64(buf[start+6]) << 8
	ret |= uint64(buf[start+7])

	return ret
}

// 加入8字节
func Put8bit(buf []byte, n byte) []byte {
	return put8bit(buf, n)
}

// 加入16字节
func Put16bit(buf []byte, n uint16) []byte {
	return put16bit(buf, n)
}

// 加入32字节
func Put32bit(buf []byte, n uint32) []byte {
	return put32bit(buf, n)
}

// 加入64字节
func Put64bit(buf []byte, n uint64) []byte {
	return put64bit(buf, n)
}
