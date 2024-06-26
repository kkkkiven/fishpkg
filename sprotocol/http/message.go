package http

import (
	"unsafe"
)

const (
	PROTO_VERSION   = 100
	FIX_HEAD_LENGTH = 0x0D
	MAX_HEAD_LENGTH = 0x20
)

// message type
const (
	MT_NORMAL byte = iota
	MT_REQUEST
	MT_RESPONSE
	MT_PING
	MT_PONG
)

// message flags
const (
	MF_ENCODE byte = 1 << iota
	MF_COMPRESS
)

type Message struct {
	// fixed head
	m1      byte   // magic word one '#'
	m2      byte   // magic word two '@'
	version byte   // protocol version
	msgType byte   // message type: normal, request, response, ping, pong
	msgFlag byte   // flag: encode, compress,
	reqID   uint32 // request id
	bodyLen uint32 // body length

	status byte // response status, ref for response message

	// route head
	exLen  uint16 // extended data length, ref for request message
	exData []byte // extended data, ref for request message

	body []byte
}

// NewMessage 构造消息
func NewMessage(msgType byte) *Message {
	msg := &Message{}
	msg.m1 = '#'
	msg.m2 = '@'
	msg.version = PROTO_VERSION
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

// NewPingMessage 创建ping消息
func NewPingMessage() *Message {
	msg := NewMessage(MT_PING)
	return msg
}

// NewPongMessage 创建pong消息
func NewPongMessage() *Message {
	msg := NewMessage(MT_PONG)
	return msg
}

// Decode 解码
func Decode(buf []byte) ([]byte, *Message) {
	pos := 0
	size := len(buf)

	for {
		if size-pos < 2 {
			return buf[pos:], nil
		}
		if buf[pos] == '#' && buf[pos+1] == '@' {
			buf, size, pos = buf[pos:], size-pos, 0
			break
		}
		pos++
	}

	if size-pos < FIX_HEAD_LENGTH {
		return buf, nil
	}

	msgType := get8bit(buf, pos+3)
	msgFlag := get8bit(buf, pos+4)
	msgReqID := get32bit(buf, pos+5)
	msgBodyLen := get32bit(buf, pos+9)

	pos = FIX_HEAD_LENGTH
	if size-pos < int(msgBodyLen) {
		return buf, nil
	}

	msg := NewMessage(msgType)
	msg.msgType = msgType
	msg.msgFlag = msgFlag
	msg.reqID = msgReqID
	msg.bodyLen = msgBodyLen

	switch msgType {
	case MT_NORMAL:
		fallthrough
	case MT_REQUEST:
		if size-pos < 2 {
			return buf, nil
		}
		msg.exLen = get16bit(buf, pos)
		pos += 2

		if size-pos < int(msg.exLen)+int(msg.bodyLen) {
			return buf, nil
		}
		msg.exData = buf[pos : pos+int(msg.exLen)]
		pos += int(msg.exLen)
	case MT_RESPONSE:
		if size-pos < 1 {
			return buf, nil
		}
		msg.status = get8bit(buf, pos)
		pos += 1
	case MT_PING:
	case MT_PONG:
	default:
		return buf[2:], nil
	}

	if pos+int(msg.bodyLen) > len(buf) {
		return buf, nil
	}

	msg.body = buf[pos : pos+int(msg.bodyLen)]

	return buf[pos+int(msg.bodyLen):], msg
}

// Encode 编码
func (this *Message) Encode() []byte {
	buf := make([]byte, 0, MAX_HEAD_LENGTH+uint32(this.exLen)+this.bodyLen)
	buf = put8bit(buf, this.m1)
	buf = put8bit(buf, this.m2)
	buf = put8bit(buf, this.version)
	buf = put8bit(buf, this.msgType)
	buf = put8bit(buf, this.msgFlag)
	buf = put32bit(buf, this.reqID)
	buf = put32bit(buf, this.bodyLen)
	switch this.msgType {
	case MT_NORMAL:
		fallthrough
	case MT_REQUEST:
		buf = put16bit(buf, this.exLen)
		if len(this.exData) > 0 {
			buf = append(buf, this.exData[:]...)
		}
	case MT_RESPONSE:
		buf = put8bit(buf, this.status)
	case MT_PING:
	case MT_PONG:
	default:
		return nil
	}

	if len(this.body) > 0 {
		buf = append(buf, this.body...)
	}

	return buf
}

// GetBodyLength 获取消息体长度
func (this *Message) GetBodyLength() uint32 {
	return this.bodyLen
}

// GetMessageType 获取消息类型
func (this *Message) GetMessageType() byte {
	return this.msgType
}

// SetMessageType 设置消息类型
func (this *Message) SetMessageType(typ byte) {
	this.msgType = typ
}

// GetMessageFlag 获取消息标志
func (this *Message) GetMessageFlag() byte {
	return this.msgFlag
}

// SetMessageFlag 设置消息标志
func (this *Message) SetMessageFlag(flag byte) {
	this.msgFlag = flag
}

// GetRequestID 获取请求id
func (this *Message) GetRequestID() uint32 {
	return this.reqID
}

// SetRequestID 设置请求id
func (this *Message) SetRequestID(id uint32) {
	this.reqID = id
}

// GetStatus 获取状态
func (this *Message) GetStatus() byte {
	return this.status
}

// SetStatus 设置状态
func (this *Message) SetStatus(code byte) {
	this.status = code
}

// GetExtendDataLen 获取附加数据长度
func (this *Message) GetExtendDataLen() uint16 {
	return this.exLen
}

// GetExData 获取附加数据
func (this *Message) GetExData() []byte {
	return this.exData
}

// SetExData 设置附加数据
func (this *Message) SetExData(data []byte) {
	this.exData = data
	this.exLen = uint16(len(data))
}

// GetStringExData 获取附加数据
func (this *Message) GetStringExData() string {
	if this.exLen == 0 {
		return ""
	}

	return *(*string)(unsafe.Pointer(&this.exData))
}

// SetStringExData 设置附加数据
func (this *Message) SetStringExData(data string) {
	this.exData = *(*[]byte)(unsafe.Pointer(&data))
	this.exLen = uint16(len(this.exData))
}

// GetBody 获取消息体
func (this *Message) GetBody() []byte {
	if this.bodyLen == 0 {
		return nil
	}

	return this.body
}

// SetBody 设置消息体
func (this *Message) SetBody(body []byte) {
	this.body = body
	this.bodyLen = uint32(len(this.body))
}

// GetStringBody 获取消息体
func (this *Message) GetStringBody() string {
	if this.bodyLen == 0 {
		return ""
	}

	return *(*string)(unsafe.Pointer(&this.body))
}

// SetStringBody 设置消息体
func (this *Message) SetStringBody(body string) {
	this.body = *(*[]byte)(unsafe.Pointer(&body))
	this.bodyLen = uint32(len(this.body))
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
