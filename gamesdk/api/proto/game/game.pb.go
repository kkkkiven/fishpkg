// Code generated by protoc-gen-go. DO NOT EDIT.
// source: game.proto

package game

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// 通用消息回应
type RespComm struct {
	Errcode              int32    `protobuf:"varint,1,opt,name=errcode,proto3" json:"errcode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RespComm) Reset()         { *m = RespComm{} }
func (m *RespComm) String() string { return proto.CompactTextString(m) }
func (*RespComm) ProtoMessage()    {}
func (*RespComm) Descriptor() ([]byte, []int) {
	return fileDescriptor_38fc58335341d769, []int{0}
}

func (m *RespComm) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RespComm.Unmarshal(m, b)
}
func (m *RespComm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RespComm.Marshal(b, m, deterministic)
}
func (m *RespComm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RespComm.Merge(m, src)
}
func (m *RespComm) XXX_Size() int {
	return xxx_messageInfo_RespComm.Size(m)
}
func (m *RespComm) XXX_DiscardUnknown() {
	xxx_messageInfo_RespComm.DiscardUnknown(m)
}

var xxx_messageInfo_RespComm proto.InternalMessageInfo

func (m *RespComm) GetErrcode() int32 {
	if m != nil {
		return m.Errcode
	}
	return 0
}

func (m *RespComm) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

// 0x1000心跳
type ReqHeartbeat struct {
	Timestamp            int64    `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqHeartbeat) Reset()         { *m = ReqHeartbeat{} }
func (m *ReqHeartbeat) String() string { return proto.CompactTextString(m) }
func (*ReqHeartbeat) ProtoMessage()    {}
func (*ReqHeartbeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_38fc58335341d769, []int{1}
}

func (m *ReqHeartbeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqHeartbeat.Unmarshal(m, b)
}
func (m *ReqHeartbeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqHeartbeat.Marshal(b, m, deterministic)
}
func (m *ReqHeartbeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqHeartbeat.Merge(m, src)
}
func (m *ReqHeartbeat) XXX_Size() int {
	return xxx_messageInfo_ReqHeartbeat.Size(m)
}
func (m *ReqHeartbeat) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqHeartbeat.DiscardUnknown(m)
}

var xxx_messageInfo_ReqHeartbeat proto.InternalMessageInfo

func (m *ReqHeartbeat) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

// 0x1001
type RespHeartbeat struct {
	Errcode              int32    `protobuf:"varint,1,opt,name=errcode,proto3" json:"errcode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Timestamp            int64    `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RespHeartbeat) Reset()         { *m = RespHeartbeat{} }
func (m *RespHeartbeat) String() string { return proto.CompactTextString(m) }
func (*RespHeartbeat) ProtoMessage()    {}
func (*RespHeartbeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_38fc58335341d769, []int{2}
}

func (m *RespHeartbeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RespHeartbeat.Unmarshal(m, b)
}
func (m *RespHeartbeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RespHeartbeat.Marshal(b, m, deterministic)
}
func (m *RespHeartbeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RespHeartbeat.Merge(m, src)
}
func (m *RespHeartbeat) XXX_Size() int {
	return xxx_messageInfo_RespHeartbeat.Size(m)
}
func (m *RespHeartbeat) XXX_DiscardUnknown() {
	xxx_messageInfo_RespHeartbeat.DiscardUnknown(m)
}

var xxx_messageInfo_RespHeartbeat proto.InternalMessageInfo

func (m *RespHeartbeat) GetErrcode() int32 {
	if m != nil {
		return m.Errcode
	}
	return 0
}

func (m *RespHeartbeat) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *RespHeartbeat) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

// 0x1002 授权验证请求
type ReqAuthorize struct {
	Uid                  int64    `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Timestamp            int64    `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Nonce                string   `protobuf:"bytes,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Sign                 string   `protobuf:"bytes,4,opt,name=sign,proto3" json:"sign,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqAuthorize) Reset()         { *m = ReqAuthorize{} }
func (m *ReqAuthorize) String() string { return proto.CompactTextString(m) }
func (*ReqAuthorize) ProtoMessage()    {}
func (*ReqAuthorize) Descriptor() ([]byte, []int) {
	return fileDescriptor_38fc58335341d769, []int{3}
}

func (m *ReqAuthorize) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqAuthorize.Unmarshal(m, b)
}
func (m *ReqAuthorize) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqAuthorize.Marshal(b, m, deterministic)
}
func (m *ReqAuthorize) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqAuthorize.Merge(m, src)
}
func (m *ReqAuthorize) XXX_Size() int {
	return xxx_messageInfo_ReqAuthorize.Size(m)
}
func (m *ReqAuthorize) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqAuthorize.DiscardUnknown(m)
}

var xxx_messageInfo_ReqAuthorize proto.InternalMessageInfo

func (m *ReqAuthorize) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *ReqAuthorize) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *ReqAuthorize) GetNonce() string {
	if m != nil {
		return m.Nonce
	}
	return ""
}

func (m *ReqAuthorize) GetSign() string {
	if m != nil {
		return m.Sign
	}
	return ""
}

// 0x1003 授权验证响应
type RespAuthorize struct {
	Errcode              int32    `protobuf:"varint,1,opt,name=errcode,proto3" json:"errcode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RespAuthorize) Reset()         { *m = RespAuthorize{} }
func (m *RespAuthorize) String() string { return proto.CompactTextString(m) }
func (*RespAuthorize) ProtoMessage()    {}
func (*RespAuthorize) Descriptor() ([]byte, []int) {
	return fileDescriptor_38fc58335341d769, []int{4}
}

func (m *RespAuthorize) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RespAuthorize.Unmarshal(m, b)
}
func (m *RespAuthorize) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RespAuthorize.Marshal(b, m, deterministic)
}
func (m *RespAuthorize) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RespAuthorize.Merge(m, src)
}
func (m *RespAuthorize) XXX_Size() int {
	return xxx_messageInfo_RespAuthorize.Size(m)
}
func (m *RespAuthorize) XXX_DiscardUnknown() {
	xxx_messageInfo_RespAuthorize.DiscardUnknown(m)
}

var xxx_messageInfo_RespAuthorize proto.InternalMessageInfo

func (m *RespAuthorize) GetErrcode() int32 {
	if m != nil {
		return m.Errcode
	}
	return 0
}

func (m *RespAuthorize) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

// 0x8004 玩家被踢通知
type NotifyKickOff struct {
	Uid                  int64             `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Msg                  string            `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Ext                  map[string][]byte `protobuf:"bytes,10,rep,name=ext,proto3" json:"ext,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *NotifyKickOff) Reset()         { *m = NotifyKickOff{} }
func (m *NotifyKickOff) String() string { return proto.CompactTextString(m) }
func (*NotifyKickOff) ProtoMessage()    {}
func (*NotifyKickOff) Descriptor() ([]byte, []int) {
	return fileDescriptor_38fc58335341d769, []int{5}
}

func (m *NotifyKickOff) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyKickOff.Unmarshal(m, b)
}
func (m *NotifyKickOff) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyKickOff.Marshal(b, m, deterministic)
}
func (m *NotifyKickOff) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyKickOff.Merge(m, src)
}
func (m *NotifyKickOff) XXX_Size() int {
	return xxx_messageInfo_NotifyKickOff.Size(m)
}
func (m *NotifyKickOff) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyKickOff.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyKickOff proto.InternalMessageInfo

func (m *NotifyKickOff) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *NotifyKickOff) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *NotifyKickOff) GetExt() map[string][]byte {
	if m != nil {
		return m.Ext
	}
	return nil
}

//MSG_ID_OTHER_LOGIN uint16 = 0x8005 //顶号
type NotifyOtherLogin struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyOtherLogin) Reset()         { *m = NotifyOtherLogin{} }
func (m *NotifyOtherLogin) String() string { return proto.CompactTextString(m) }
func (*NotifyOtherLogin) ProtoMessage()    {}
func (*NotifyOtherLogin) Descriptor() ([]byte, []int) {
	return fileDescriptor_38fc58335341d769, []int{6}
}

func (m *NotifyOtherLogin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyOtherLogin.Unmarshal(m, b)
}
func (m *NotifyOtherLogin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyOtherLogin.Marshal(b, m, deterministic)
}
func (m *NotifyOtherLogin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyOtherLogin.Merge(m, src)
}
func (m *NotifyOtherLogin) XXX_Size() int {
	return xxx_messageInfo_NotifyOtherLogin.Size(m)
}
func (m *NotifyOtherLogin) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyOtherLogin.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyOtherLogin proto.InternalMessageInfo

func init() {
	proto.RegisterType((*RespComm)(nil), "game.RespComm")
	proto.RegisterType((*ReqHeartbeat)(nil), "game.ReqHeartbeat")
	proto.RegisterType((*RespHeartbeat)(nil), "game.RespHeartbeat")
	proto.RegisterType((*ReqAuthorize)(nil), "game.ReqAuthorize")
	proto.RegisterType((*RespAuthorize)(nil), "game.RespAuthorize")
	proto.RegisterType((*NotifyKickOff)(nil), "game.NotifyKickOff")
	proto.RegisterMapType((map[string][]byte)(nil), "game.NotifyKickOff.ExtEntry")
	proto.RegisterType((*NotifyOtherLogin)(nil), "game.NotifyOtherLogin")
}

func init() { proto.RegisterFile("game.proto", fileDescriptor_38fc58335341d769) }

var fileDescriptor_38fc58335341d769 = []byte{
	// 311 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x41, 0x4b, 0xfb, 0x40,
	0x10, 0xc5, 0x49, 0xd3, 0xfe, 0xdb, 0xcc, 0xbf, 0x85, 0xb2, 0x78, 0x08, 0xd2, 0x43, 0xc9, 0x29,
	0x07, 0x59, 0x41, 0xa1, 0x88, 0x9e, 0x54, 0x0a, 0x82, 0x62, 0x61, 0x6f, 0x7a, 0x4b, 0xd3, 0x49,
	0xb2, 0xd4, 0xcd, 0xc6, 0xcd, 0x46, 0x1a, 0xbf, 0x8a, 0x5f, 0x56, 0x76, 0x63, 0x09, 0xa9, 0x5e,
	0x7a, 0x7b, 0x33, 0xcc, 0xfb, 0xbd, 0x99, 0x64, 0x01, 0xd2, 0x48, 0x20, 0x2d, 0x94, 0xd4, 0x92,
	0xf4, 0x8d, 0x0e, 0x16, 0x30, 0x62, 0x58, 0x16, 0xf7, 0x52, 0x08, 0xe2, 0xc3, 0x10, 0x95, 0x8a,
	0xe5, 0x06, 0x7d, 0x67, 0xee, 0x84, 0x03, 0xb6, 0x2f, 0xc9, 0x14, 0x5c, 0x51, 0xa6, 0x7e, 0x6f,
	0xee, 0x84, 0x1e, 0x33, 0x32, 0x38, 0x83, 0x31, 0xc3, 0xf7, 0x07, 0x8c, 0x94, 0x5e, 0x63, 0xa4,
	0xc9, 0x0c, 0x3c, 0xcd, 0x05, 0x96, 0x3a, 0x12, 0x85, 0x75, 0xbb, 0xac, 0x6d, 0x04, 0x2f, 0x30,
	0x31, 0x29, 0xed, 0xf8, 0x11, 0x51, 0x5d, 0xb4, 0x7b, 0x88, 0xce, 0xec, 0x22, 0xb7, 0x95, 0xce,
	0xa4, 0xe2, 0x9f, 0xd6, 0x5f, 0xf1, 0xcd, 0xcf, 0x0a, 0x46, 0x76, 0xfd, 0xbd, 0x03, 0x3f, 0x39,
	0x81, 0x41, 0x2e, 0xf3, 0x18, 0x2d, 0xd9, 0x63, 0x4d, 0x41, 0x08, 0xf4, 0x4b, 0x9e, 0xe6, 0x7e,
	0xdf, 0x36, 0xad, 0x0e, 0x6e, 0x9a, 0x23, 0xda, 0xa8, 0x63, 0xbe, 0xd7, 0x97, 0x03, 0x93, 0x67,
	0xa9, 0x79, 0x52, 0x3f, 0xf2, 0x78, 0xbb, 0x4a, 0x92, 0x3f, 0x16, 0xfd, 0x7d, 0x3a, 0x05, 0x17,
	0x77, 0xda, 0x87, 0xb9, 0x1b, 0xfe, 0xbf, 0x98, 0x51, 0xfb, 0xf7, 0x3a, 0x14, 0xba, 0xdc, 0xe9,
	0x65, 0xae, 0x55, 0xcd, 0xcc, 0xe0, 0xe9, 0x02, 0x46, 0xfb, 0x86, 0xa1, 0x6d, 0xb1, 0xb6, 0x7c,
	0x8f, 0x19, 0x69, 0x4e, 0xfd, 0x88, 0xde, 0x2a, 0xb4, 0x09, 0x63, 0xd6, 0x14, 0xd7, 0xbd, 0x2b,
	0x27, 0x20, 0x30, 0x6d, 0xb0, 0x2b, 0x9d, 0xa1, 0x7a, 0x92, 0x29, 0xcf, 0xef, 0xbc, 0xd7, 0x21,
	0xa5, 0xe7, 0x26, 0x72, 0xfd, 0xcf, 0xbe, 0x98, 0xcb, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3c,
	0x1d, 0x58, 0xdb, 0x3f, 0x02, 0x00, 0x00,
}
