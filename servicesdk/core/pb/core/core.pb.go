// Code generated by protoc-gen-go. DO NOT EDIT.
// source: core.proto

package pb

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

type RspMsg struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RspMsg) Reset()         { *m = RspMsg{} }
func (m *RspMsg) String() string { return proto.CompactTextString(m) }
func (*RspMsg) ProtoMessage()    {}
func (*RspMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7e43720d1edc0fe, []int{0}
}

func (m *RspMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RspMsg.Unmarshal(m, b)
}
func (m *RspMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RspMsg.Marshal(b, m, deterministic)
}
func (m *RspMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RspMsg.Merge(m, src)
}
func (m *RspMsg) XXX_Size() int {
	return xxx_messageInfo_RspMsg.Size(m)
}
func (m *RspMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_RspMsg.DiscardUnknown(m)
}

var xxx_messageInfo_RspMsg proto.InternalMessageInfo

func (m *RspMsg) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *RspMsg) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type RegMsg struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 uint32   `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	Weight               int32    `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
	Secret               string   `protobuf:"bytes,4,opt,name=secret,proto3" json:"secret,omitempty"`
	Name                 string   `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegMsg) Reset()         { *m = RegMsg{} }
func (m *RegMsg) String() string { return proto.CompactTextString(m) }
func (*RegMsg) ProtoMessage()    {}
func (*RegMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7e43720d1edc0fe, []int{1}
}

func (m *RegMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegMsg.Unmarshal(m, b)
}
func (m *RegMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegMsg.Marshal(b, m, deterministic)
}
func (m *RegMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegMsg.Merge(m, src)
}
func (m *RegMsg) XXX_Size() int {
	return xxx_messageInfo_RegMsg.Size(m)
}
func (m *RegMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_RegMsg.DiscardUnknown(m)
}

var xxx_messageInfo_RegMsg proto.InternalMessageInfo

func (m *RegMsg) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *RegMsg) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *RegMsg) GetWeight() int32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *RegMsg) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

func (m *RegMsg) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type UpdateMsg struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 uint32   `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	Weight               int32    `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateMsg) Reset()         { *m = UpdateMsg{} }
func (m *UpdateMsg) String() string { return proto.CompactTextString(m) }
func (*UpdateMsg) ProtoMessage()    {}
func (*UpdateMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7e43720d1edc0fe, []int{2}
}

func (m *UpdateMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateMsg.Unmarshal(m, b)
}
func (m *UpdateMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateMsg.Marshal(b, m, deterministic)
}
func (m *UpdateMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateMsg.Merge(m, src)
}
func (m *UpdateMsg) XXX_Size() int {
	return xxx_messageInfo_UpdateMsg.Size(m)
}
func (m *UpdateMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateMsg.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateMsg proto.InternalMessageInfo

func (m *UpdateMsg) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateMsg) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *UpdateMsg) GetWeight() int32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

// SLSMsg sls日志消息
type SLSMsg struct {
	Store                string            `protobuf:"bytes,1,opt,name=store,proto3" json:"store,omitempty"`
	Topic                string            `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Contents             map[string]string `protobuf:"bytes,3,rep,name=contents,proto3" json:"contents,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SLSMsg) Reset()         { *m = SLSMsg{} }
func (m *SLSMsg) String() string { return proto.CompactTextString(m) }
func (*SLSMsg) ProtoMessage()    {}
func (*SLSMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7e43720d1edc0fe, []int{3}
}

func (m *SLSMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SLSMsg.Unmarshal(m, b)
}
func (m *SLSMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SLSMsg.Marshal(b, m, deterministic)
}
func (m *SLSMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SLSMsg.Merge(m, src)
}
func (m *SLSMsg) XXX_Size() int {
	return xxx_messageInfo_SLSMsg.Size(m)
}
func (m *SLSMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_SLSMsg.DiscardUnknown(m)
}

var xxx_messageInfo_SLSMsg proto.InternalMessageInfo

func (m *SLSMsg) GetStore() string {
	if m != nil {
		return m.Store
	}
	return ""
}

func (m *SLSMsg) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *SLSMsg) GetContents() map[string]string {
	if m != nil {
		return m.Contents
	}
	return nil
}

func init() {
	proto.RegisterType((*RspMsg)(nil), "pb.RspMsg")
	proto.RegisterType((*RegMsg)(nil), "pb.RegMsg")
	proto.RegisterType((*UpdateMsg)(nil), "pb.UpdateMsg")
	proto.RegisterType((*SLSMsg)(nil), "pb.SLSMsg")
	proto.RegisterMapType((map[string]string)(nil), "pb.SLSMsg.ContentsEntry")
}

func init() { proto.RegisterFile("core.proto", fileDescriptor_f7e43720d1edc0fe) }

var fileDescriptor_f7e43720d1edc0fe = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x51, 0xbf, 0x4a, 0xfe, 0x40,
	0x10, 0x24, 0xc9, 0x97, 0xf0, 0xcb, 0xfe, 0x88, 0xc8, 0x21, 0x72, 0x58, 0x85, 0x54, 0xa9, 0x52,
	0xa8, 0x85, 0x68, 0x29, 0x62, 0xa3, 0xcd, 0x7d, 0xf8, 0x00, 0xf9, 0xb3, 0xc4, 0xa0, 0x5f, 0xee,
	0xc8, 0xad, 0x4a, 0x9e, 0xc8, 0xd7, 0x94, 0xdd, 0x9c, 0x82, 0xb5, 0xdd, 0xcc, 0xec, 0xdc, 0x0c,
	0x7b, 0x0b, 0xd0, 0xdb, 0x05, 0x1b, 0xb7, 0x58, 0xb2, 0x2a, 0x76, 0x5d, 0xd5, 0x40, 0x66, 0xbc,
	0x7b, 0xf4, 0xa3, 0x52, 0xb0, 0xeb, 0xed, 0x80, 0x3a, 0x2a, 0xa3, 0x3a, 0x35, 0x82, 0xd5, 0x31,
	0x24, 0x07, 0x3f, 0xea, 0xb8, 0x8c, 0xea, 0xdc, 0x30, 0xac, 0x1c, 0x64, 0x06, 0x47, 0xf6, 0x1f,
	0x41, 0x3c, 0x0d, 0xe2, 0x2e, 0x4c, 0x3c, 0x0d, 0xfc, 0x9e, 0x56, 0x87, 0x62, 0x2e, 0x8c, 0x60,
	0x75, 0x0a, 0xd9, 0x07, 0x4e, 0xe3, 0x33, 0xe9, 0x44, 0x52, 0x03, 0x63, 0xdd, 0x63, 0xbf, 0x20,
	0xe9, 0x9d, 0x44, 0x07, 0xc6, 0x19, 0x73, 0x7b, 0x40, 0x9d, 0x8a, 0x2a, 0xb8, 0xba, 0x87, 0xfc,
	0xc9, 0x0d, 0x2d, 0xe1, 0x1f, 0x4b, 0xab, 0xcf, 0x08, 0xb2, 0xfd, 0xc3, 0x9e, 0x63, 0x4e, 0x20,
	0xf5, 0x64, 0x97, 0x6d, 0xd9, 0xdc, 0x6c, 0x84, 0x55, 0xb2, 0x6e, 0xea, 0xc3, 0xbe, 0x1b, 0x51,
	0x97, 0xf0, 0xaf, 0xb7, 0x33, 0xe1, 0x4c, 0x5e, 0x27, 0x65, 0x52, 0xff, 0x3f, 0xd7, 0x8d, 0xeb,
	0x9a, 0x2d, 0xa9, 0xb9, 0x0d, 0xa3, 0xbb, 0x99, 0x96, 0xd5, 0xfc, 0x38, 0xcf, 0x6e, 0xa0, 0xf8,
	0x35, 0xe2, 0xaf, 0x7c, 0xc1, 0x35, 0x14, 0x32, 0xe4, 0xba, 0xf7, 0xf6, 0xf5, 0x0d, 0xbf, 0xeb,
	0x84, 0x5c, 0xc7, 0x57, 0x51, 0x97, 0xc9, 0x7d, 0x2e, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x80,
	0xcf, 0x39, 0x48, 0xad, 0x01, 0x00, 0x00,
}
