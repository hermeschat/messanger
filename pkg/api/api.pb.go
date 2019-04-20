// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/api/api.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type DestroySessionRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DestroySessionRequest) Reset()         { *m = DestroySessionRequest{} }
func (m *DestroySessionRequest) String() string { return proto.CompactTextString(m) }
func (*DestroySessionRequest) ProtoMessage()    {}
func (*DestroySessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e50ccc58c7b575d, []int{0}
}

func (m *DestroySessionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DestroySessionRequest.Unmarshal(m, b)
}
func (m *DestroySessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DestroySessionRequest.Marshal(b, m, deterministic)
}
func (m *DestroySessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DestroySessionRequest.Merge(m, src)
}
func (m *DestroySessionRequest) XXX_Size() int {
	return xxx_messageInfo_DestroySessionRequest.Size(m)
}
func (m *DestroySessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DestroySessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DestroySessionRequest proto.InternalMessageInfo

func (m *DestroySessionRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type DestroySessionResponse struct {
	StatusCode           string   `protobuf:"bytes,1,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DestroySessionResponse) Reset()         { *m = DestroySessionResponse{} }
func (m *DestroySessionResponse) String() string { return proto.CompactTextString(m) }
func (*DestroySessionResponse) ProtoMessage()    {}
func (*DestroySessionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e50ccc58c7b575d, []int{1}
}

func (m *DestroySessionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DestroySessionResponse.Unmarshal(m, b)
}
func (m *DestroySessionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DestroySessionResponse.Marshal(b, m, deterministic)
}
func (m *DestroySessionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DestroySessionResponse.Merge(m, src)
}
func (m *DestroySessionResponse) XXX_Size() int {
	return xxx_messageInfo_DestroySessionResponse.Size(m)
}
func (m *DestroySessionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DestroySessionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DestroySessionResponse proto.InternalMessageInfo

func (m *DestroySessionResponse) GetStatusCode() string {
	if m != nil {
		return m.StatusCode
	}
	return ""
}

type Signal struct {
	SignalType           string   `protobuf:"bytes,1,opt,name=signalType,proto3" json:"signalType,omitempty"`
	From                 string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Payload              string   `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Signal) Reset()         { *m = Signal{} }
func (m *Signal) String() string { return proto.CompactTextString(m) }
func (*Signal) ProtoMessage()    {}
func (*Signal) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e50ccc58c7b575d, []int{2}
}

func (m *Signal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Signal.Unmarshal(m, b)
}
func (m *Signal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Signal.Marshal(b, m, deterministic)
}
func (m *Signal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Signal.Merge(m, src)
}
func (m *Signal) XXX_Size() int {
	return xxx_messageInfo_Signal.Size(m)
}
func (m *Signal) XXX_DiscardUnknown() {
	xxx_messageInfo_Signal.DiscardUnknown(m)
}

var xxx_messageInfo_Signal proto.InternalMessageInfo

func (m *Signal) GetSignalType() string {
	if m != nil {
		return m.SignalType
	}
	return ""
}

func (m *Signal) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *Signal) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

type InstantMessage struct {
	MessageType          string   `protobuf:"bytes,1,opt,name=messageType,proto3" json:"messageType,omitempty"`
	From                 string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Channel              string   `protobuf:"bytes,3,opt,name=channel,proto3" json:"channel,omitempty"`
	To                   string   `protobuf:"bytes,4,opt,name=to,proto3" json:"to,omitempty"`
	Body                 string   `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstantMessage) Reset()         { *m = InstantMessage{} }
func (m *InstantMessage) String() string { return proto.CompactTextString(m) }
func (*InstantMessage) ProtoMessage()    {}
func (*InstantMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e50ccc58c7b575d, []int{3}
}

func (m *InstantMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstantMessage.Unmarshal(m, b)
}
func (m *InstantMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstantMessage.Marshal(b, m, deterministic)
}
func (m *InstantMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstantMessage.Merge(m, src)
}
func (m *InstantMessage) XXX_Size() int {
	return xxx_messageInfo_InstantMessage.Size(m)
}
func (m *InstantMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_InstantMessage.DiscardUnknown(m)
}

var xxx_messageInfo_InstantMessage proto.InternalMessageInfo

func (m *InstantMessage) GetMessageType() string {
	if m != nil {
		return m.MessageType
	}
	return ""
}

func (m *InstantMessage) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *InstantMessage) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *InstantMessage) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *InstantMessage) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

type CreateSessionRequest struct {
	ClientType           string   `protobuf:"bytes,1,opt,name=ClientType,proto3" json:"ClientType,omitempty"`
	UserID               string   `protobuf:"bytes,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	UserIP               string   `protobuf:"bytes,3,opt,name=UserIP,proto3" json:"UserIP,omitempty"`
	UserAgent            string   `protobuf:"bytes,4,opt,name=UserAgent,proto3" json:"UserAgent,omitempty"`
	ClientVersion        string   `protobuf:"bytes,5,opt,name=ClientVersion,proto3" json:"ClientVersion,omitempty"`
	Node                 string   `protobuf:"bytes,6,opt,name=Node,proto3" json:"Node,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateSessionRequest) Reset()         { *m = CreateSessionRequest{} }
func (m *CreateSessionRequest) String() string { return proto.CompactTextString(m) }
func (*CreateSessionRequest) ProtoMessage()    {}
func (*CreateSessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e50ccc58c7b575d, []int{4}
}

func (m *CreateSessionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateSessionRequest.Unmarshal(m, b)
}
func (m *CreateSessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateSessionRequest.Marshal(b, m, deterministic)
}
func (m *CreateSessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateSessionRequest.Merge(m, src)
}
func (m *CreateSessionRequest) XXX_Size() int {
	return xxx_messageInfo_CreateSessionRequest.Size(m)
}
func (m *CreateSessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateSessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateSessionRequest proto.InternalMessageInfo

func (m *CreateSessionRequest) GetClientType() string {
	if m != nil {
		return m.ClientType
	}
	return ""
}

func (m *CreateSessionRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *CreateSessionRequest) GetUserIP() string {
	if m != nil {
		return m.UserIP
	}
	return ""
}

func (m *CreateSessionRequest) GetUserAgent() string {
	if m != nil {
		return m.UserAgent
	}
	return ""
}

func (m *CreateSessionRequest) GetClientVersion() string {
	if m != nil {
		return m.ClientVersion
	}
	return ""
}

func (m *CreateSessionRequest) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

type CreateSessionResponse struct {
	StatusCode           string   `protobuf:"bytes,1,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateSessionResponse) Reset()         { *m = CreateSessionResponse{} }
func (m *CreateSessionResponse) String() string { return proto.CompactTextString(m) }
func (*CreateSessionResponse) ProtoMessage()    {}
func (*CreateSessionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e50ccc58c7b575d, []int{5}
}

func (m *CreateSessionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateSessionResponse.Unmarshal(m, b)
}
func (m *CreateSessionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateSessionResponse.Marshal(b, m, deterministic)
}
func (m *CreateSessionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateSessionResponse.Merge(m, src)
}
func (m *CreateSessionResponse) XXX_Size() int {
	return xxx_messageInfo_CreateSessionResponse.Size(m)
}
func (m *CreateSessionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateSessionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateSessionResponse proto.InternalMessageInfo

func (m *CreateSessionResponse) GetStatusCode() string {
	if m != nil {
		return m.StatusCode
	}
	return ""
}

type Response struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e50ccc58c7b575d, []int{6}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Response) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*DestroySessionRequest)(nil), "DestroySessionRequest")
	proto.RegisterType((*DestroySessionResponse)(nil), "DestroySessionResponse")
	proto.RegisterType((*Signal)(nil), "Signal")
	proto.RegisterType((*InstantMessage)(nil), "InstantMessage")
	proto.RegisterType((*CreateSessionRequest)(nil), "CreateSessionRequest")
	proto.RegisterType((*CreateSessionResponse)(nil), "CreateSessionResponse")
	proto.RegisterType((*Response)(nil), "Response")
}

func init() { proto.RegisterFile("pkg/api/api.proto", fileDescriptor_7e50ccc58c7b575d) }

var fileDescriptor_7e50ccc58c7b575d = []byte{
	// 423 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xcf, 0xcb, 0xd3, 0x40,
	0x10, 0x25, 0xb5, 0x5f, 0xbe, 0x2f, 0x23, 0xad, 0xb8, 0xb4, 0x31, 0x14, 0xa9, 0x25, 0x78, 0xe8,
	0x69, 0x05, 0x7f, 0xa0, 0x27, 0xa1, 0xb4, 0x07, 0x8b, 0x58, 0xa4, 0xd5, 0xde, 0xb7, 0xcd, 0x18,
	0x83, 0xe9, 0x6e, 0xdc, 0xdd, 0x2a, 0xb9, 0xfb, 0x77, 0xf9, 0x67, 0x79, 0x96, 0x4d, 0xb6, 0xf9,
	0x51, 0x22, 0x78, 0x08, 0xbc, 0x79, 0x99, 0x99, 0x7d, 0xb3, 0x6f, 0x07, 0x1e, 0x66, 0xdf, 0xe2,
	0x67, 0x2c, 0x4b, 0xcc, 0x47, 0x33, 0x29, 0xb4, 0x08, 0x5f, 0xc1, 0x78, 0x85, 0x4a, 0x4b, 0x91,
	0xef, 0x50, 0xa9, 0x44, 0xf0, 0x2d, 0x7e, 0x3f, 0xa3, 0xd2, 0xe4, 0x31, 0x78, 0x96, 0x59, 0x47,
	0x81, 0x33, 0x73, 0xe6, 0xde, 0xb6, 0x26, 0xc2, 0x37, 0xe0, 0x5f, 0x97, 0xa9, 0x4c, 0x70, 0x85,
	0x64, 0x0a, 0xb0, 0xd3, 0x4c, 0x9f, 0xd5, 0x52, 0x44, 0x68, 0x0b, 0x1b, 0x4c, 0xb8, 0x07, 0x77,
	0x97, 0xc4, 0x9c, 0xa5, 0x26, 0x53, 0x15, 0xe8, 0x53, 0x9e, 0x55, 0x99, 0x35, 0x43, 0x08, 0xf4,
	0xbf, 0x48, 0x71, 0x0a, 0x7a, 0xc5, 0x9f, 0x02, 0x93, 0x00, 0x6e, 0x33, 0x96, 0xa7, 0x82, 0x45,
	0xc1, 0xbd, 0x82, 0xbe, 0x84, 0xe1, 0x2f, 0x07, 0x86, 0x6b, 0xae, 0x34, 0xe3, 0xfa, 0x03, 0x2a,
	0xc5, 0x62, 0x24, 0x33, 0xb8, 0x7f, 0x2a, 0x61, 0xe3, 0x84, 0x26, 0xf5, 0xaf, 0x23, 0x8e, 0x5f,
	0x19, 0xe7, 0x98, 0x5e, 0x8e, 0xb0, 0x21, 0x19, 0x42, 0x4f, 0x8b, 0xa0, 0x5f, 0x90, 0x3d, 0x2d,
	0x4c, 0xf5, 0x41, 0x44, 0x79, 0x70, 0x53, 0x56, 0x1b, 0x1c, 0xfe, 0x76, 0x60, 0xb4, 0x94, 0xc8,
	0x34, 0x5e, 0xdd, 0xe7, 0x14, 0x60, 0x99, 0x26, 0xc8, 0x75, 0x73, 0xda, 0x9a, 0x21, 0x3e, 0xb8,
	0x9f, 0x15, 0xca, 0xf5, 0xca, 0x8a, 0xb1, 0x51, 0xc5, 0x7f, 0xb4, 0x6a, 0x6c, 0x64, 0xfc, 0x31,
	0x68, 0x11, 0x23, 0xd7, 0x56, 0x53, 0x4d, 0x90, 0xa7, 0x30, 0x28, 0x7b, 0xef, 0x51, 0x1a, 0x15,
	0x56, 0x63, 0x9b, 0x34, 0x03, 0x6c, 0x8c, 0x4b, 0x6e, 0x39, 0x80, 0xc1, 0xe1, 0x6b, 0x18, 0x5f,
	0xe9, 0xff, 0x4f, 0x63, 0x5f, 0xc2, 0x5d, 0x95, 0x4b, 0xa0, 0x7f, 0xac, 0xb3, 0x0a, 0x4c, 0x46,
	0x70, 0x83, 0x52, 0x0a, 0x69, 0xe7, 0x2b, 0x83, 0xe7, 0x7f, 0x1c, 0x70, 0xdf, 0xa1, 0x3c, 0xa1,
	0x22, 0x4f, 0xc0, 0x7b, 0x8f, 0x98, 0x2d, 0xd2, 0xe4, 0x07, 0x92, 0x5b, 0x5a, 0xbe, 0x92, 0x89,
	0x47, 0xab, 0xae, 0x73, 0x80, 0x0d, 0xfe, 0xbc, 0xb8, 0xfb, 0x80, 0xb6, 0xed, 0x6e, 0x66, 0x4e,
	0xe1, 0x6e, 0x85, 0xa6, 0x8f, 0x8c, 0x3a, 0x3b, 0xbd, 0x85, 0x41, 0x6b, 0x48, 0x32, 0xa6, 0x5d,
	0xa6, 0x4d, 0x7c, 0xda, 0x7d, 0x17, 0x0b, 0x18, 0xb6, 0x9f, 0x3f, 0xf1, 0x69, 0xe7, 0x1a, 0x4d,
	0x1e, 0xd1, 0xee, 0x3d, 0x39, 0xb8, 0xc5, 0xfe, 0xbd, 0xf8, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xe8,
	0x9f, 0x82, 0x04, 0x94, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HermesClient is the client API for Hermes service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HermesClient interface {
	KeepAlive(ctx context.Context, in *Signal, opts ...grpc.CallOption) (*Response, error)
	NewMessage(ctx context.Context, in *InstantMessage, opts ...grpc.CallOption) (*Response, error)
	Deliverd(ctx context.Context, in *Signal, opts ...grpc.CallOption) (*Response, error)
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error)
	DestroySession(ctx context.Context, in *DestroySessionRequest, opts ...grpc.CallOption) (*DestroySessionResponse, error)
}

type hermesClient struct {
	cc *grpc.ClientConn
}

func NewHermesClient(cc *grpc.ClientConn) HermesClient {
	return &hermesClient{cc}
}

func (c *hermesClient) KeepAlive(ctx context.Context, in *Signal, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/Hermes/KeepAlive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hermesClient) NewMessage(ctx context.Context, in *InstantMessage, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/Hermes/NewMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hermesClient) Deliverd(ctx context.Context, in *Signal, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/Hermes/Deliverd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hermesClient) CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error) {
	out := new(CreateSessionResponse)
	err := c.cc.Invoke(ctx, "/Hermes/CreateSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hermesClient) DestroySession(ctx context.Context, in *DestroySessionRequest, opts ...grpc.CallOption) (*DestroySessionResponse, error) {
	out := new(DestroySessionResponse)
	err := c.cc.Invoke(ctx, "/Hermes/DestroySession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HermesServer is the server API for Hermes service.
type HermesServer interface {
	KeepAlive(context.Context, *Signal) (*Response, error)
	NewMessage(context.Context, *InstantMessage) (*Response, error)
	Deliverd(context.Context, *Signal) (*Response, error)
	CreateSession(context.Context, *CreateSessionRequest) (*CreateSessionResponse, error)
	DestroySession(context.Context, *DestroySessionRequest) (*DestroySessionResponse, error)
}

// UnimplementedHermesServer can be embedded to have forward compatible implementations.
type UnimplementedHermesServer struct {
}

func (*UnimplementedHermesServer) KeepAlive(ctx context.Context, req *Signal) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KeepAlive not implemented")
}
func (*UnimplementedHermesServer) NewMessage(ctx context.Context, req *InstantMessage) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewMessage not implemented")
}
func (*UnimplementedHermesServer) Deliverd(ctx context.Context, req *Signal) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deliverd not implemented")
}
func (*UnimplementedHermesServer) CreateSession(ctx context.Context, req *CreateSessionRequest) (*CreateSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (*UnimplementedHermesServer) DestroySession(ctx context.Context, req *DestroySessionRequest) (*DestroySessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DestroySession not implemented")
}

func RegisterHermesServer(s *grpc.Server, srv HermesServer) {
	s.RegisterService(&_Hermes_serviceDesc, srv)
}

func _Hermes_KeepAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Signal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HermesServer).KeepAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Hermes/KeepAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HermesServer).KeepAlive(ctx, req.(*Signal))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hermes_NewMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstantMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HermesServer).NewMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Hermes/NewMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HermesServer).NewMessage(ctx, req.(*InstantMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hermes_Deliverd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Signal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HermesServer).Deliverd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Hermes/Deliverd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HermesServer).Deliverd(ctx, req.(*Signal))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hermes_CreateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HermesServer).CreateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Hermes/CreateSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HermesServer).CreateSession(ctx, req.(*CreateSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hermes_DestroySession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DestroySessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HermesServer).DestroySession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Hermes/DestroySession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HermesServer).DestroySession(ctx, req.(*DestroySessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Hermes_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Hermes",
	HandlerType: (*HermesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "KeepAlive",
			Handler:    _Hermes_KeepAlive_Handler,
		},
		{
			MethodName: "NewMessage",
			Handler:    _Hermes_NewMessage_Handler,
		},
		{
			MethodName: "Deliverd",
			Handler:    _Hermes_Deliverd_Handler,
		},
		{
			MethodName: "CreateSession",
			Handler:    _Hermes_CreateSession_Handler,
		},
		{
			MethodName: "DestroySession",
			Handler:    _Hermes_DestroySession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/api.proto",
}
