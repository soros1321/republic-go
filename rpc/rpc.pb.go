// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc.proto

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	rpc.proto

It has these top-level messages:
	OrderFragment
	ResultFragment
	Address
	MultiAddress
	MultiAddresses
	Nothing
*/
package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// An OrderFragment is a message contains the details of an order fragment.
type OrderFragment struct {
	// Network data.
	To   *Address `protobuf:"bytes,1,opt,name=to" json:"to,omitempty"`
	From *Address `protobuf:"bytes,2,opt,name=from" json:"from,omitempty"`
	// Public data.
	Id           []byte `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	OrderId      []byte `protobuf:"bytes,4,opt,name=orderId,proto3" json:"orderId,omitempty"`
	OrderType    int64  `protobuf:"varint,5,opt,name=orderType" json:"orderType,omitempty"`
	OrderBuySell int64  `protobuf:"varint,6,opt,name=orderBuySell" json:"orderBuySell,omitempty"`
	// Private data.
	FstCodeShare   []byte `protobuf:"bytes,7,opt,name=fstCodeShare,proto3" json:"fstCodeShare,omitempty"`
	SndCodeShare   []byte `protobuf:"bytes,8,opt,name=sndCodeShare,proto3" json:"sndCodeShare,omitempty"`
	PriceShare     []byte `protobuf:"bytes,9,opt,name=priceShare,proto3" json:"priceShare,omitempty"`
	MaxVolumeShare []byte `protobuf:"bytes,10,opt,name=maxVolumeShare,proto3" json:"maxVolumeShare,omitempty"`
	MinVolumeShare []byte `protobuf:"bytes,11,opt,name=minVolumeShare,proto3" json:"minVolumeShare,omitempty"`
}

func (m *OrderFragment) Reset()                    { *m = OrderFragment{} }
func (m *OrderFragment) String() string            { return proto.CompactTextString(m) }
func (*OrderFragment) ProtoMessage()               {}
func (*OrderFragment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *OrderFragment) GetTo() *Address {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *OrderFragment) GetFrom() *Address {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *OrderFragment) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *OrderFragment) GetOrderId() []byte {
	if m != nil {
		return m.OrderId
	}
	return nil
}

func (m *OrderFragment) GetOrderType() int64 {
	if m != nil {
		return m.OrderType
	}
	return 0
}

func (m *OrderFragment) GetOrderBuySell() int64 {
	if m != nil {
		return m.OrderBuySell
	}
	return 0
}

func (m *OrderFragment) GetFstCodeShare() []byte {
	if m != nil {
		return m.FstCodeShare
	}
	return nil
}

func (m *OrderFragment) GetSndCodeShare() []byte {
	if m != nil {
		return m.SndCodeShare
	}
	return nil
}

func (m *OrderFragment) GetPriceShare() []byte {
	if m != nil {
		return m.PriceShare
	}
	return nil
}

func (m *OrderFragment) GetMaxVolumeShare() []byte {
	if m != nil {
		return m.MaxVolumeShare
	}
	return nil
}

func (m *OrderFragment) GetMinVolumeShare() []byte {
	if m != nil {
		return m.MinVolumeShare
	}
	return nil
}

// A ResultFragment message is the network representation of a
// compute.ResultFragment and the metadata needed to distribute it through the
// network.
type ResultFragment struct {
	// Network data.
	To   *Address `protobuf:"bytes,1,opt,name=to" json:"to,omitempty"`
	From *Address `protobuf:"bytes,2,opt,name=from" json:"from,omitempty"`
	// Public data.
	Id                  []byte `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	BuyOrderId          []byte `protobuf:"bytes,4,opt,name=buyOrderId,proto3" json:"buyOrderId,omitempty"`
	SellOrderId         []byte `protobuf:"bytes,5,opt,name=sellOrderId,proto3" json:"sellOrderId,omitempty"`
	BuyOrderFragmentId  []byte `protobuf:"bytes,6,opt,name=buyOrderFragmentId,proto3" json:"buyOrderFragmentId,omitempty"`
	SellOrderFragmentId []byte `protobuf:"bytes,7,opt,name=sellOrderFragmentId,proto3" json:"sellOrderFragmentId,omitempty"`
	// Private data.
	FstCodeShare   []byte `protobuf:"bytes,8,opt,name=fstCodeShare,proto3" json:"fstCodeShare,omitempty"`
	SndCodeShare   []byte `protobuf:"bytes,9,opt,name=sndCodeShare,proto3" json:"sndCodeShare,omitempty"`
	PriceShare     []byte `protobuf:"bytes,10,opt,name=priceShare,proto3" json:"priceShare,omitempty"`
	MaxVolumeShare []byte `protobuf:"bytes,11,opt,name=maxVolumeShare,proto3" json:"maxVolumeShare,omitempty"`
	MinVolumeShare []byte `protobuf:"bytes,12,opt,name=minVolumeShare,proto3" json:"minVolumeShare,omitempty"`
}

func (m *ResultFragment) Reset()                    { *m = ResultFragment{} }
func (m *ResultFragment) String() string            { return proto.CompactTextString(m) }
func (*ResultFragment) ProtoMessage()               {}
func (*ResultFragment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ResultFragment) GetTo() *Address {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *ResultFragment) GetFrom() *Address {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *ResultFragment) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *ResultFragment) GetBuyOrderId() []byte {
	if m != nil {
		return m.BuyOrderId
	}
	return nil
}

func (m *ResultFragment) GetSellOrderId() []byte {
	if m != nil {
		return m.SellOrderId
	}
	return nil
}

func (m *ResultFragment) GetBuyOrderFragmentId() []byte {
	if m != nil {
		return m.BuyOrderFragmentId
	}
	return nil
}

func (m *ResultFragment) GetSellOrderFragmentId() []byte {
	if m != nil {
		return m.SellOrderFragmentId
	}
	return nil
}

func (m *ResultFragment) GetFstCodeShare() []byte {
	if m != nil {
		return m.FstCodeShare
	}
	return nil
}

func (m *ResultFragment) GetSndCodeShare() []byte {
	if m != nil {
		return m.SndCodeShare
	}
	return nil
}

func (m *ResultFragment) GetPriceShare() []byte {
	if m != nil {
		return m.PriceShare
	}
	return nil
}

func (m *ResultFragment) GetMaxVolumeShare() []byte {
	if m != nil {
		return m.MaxVolumeShare
	}
	return nil
}

func (m *ResultFragment) GetMinVolumeShare() []byte {
	if m != nil {
		return m.MinVolumeShare
	}
	return nil
}

// An Address message is the network representation of an identity.Address.
type Address struct {
	Address string `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
}

func (m *Address) Reset()                    { *m = Address{} }
func (m *Address) String() string            { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()               {}
func (*Address) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Address) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

// A MultiAddress is the public multiaddress of a Node in the overlay network.
// It provides the Republic address of the Node, as well as the network
// address.
type MultiAddress struct {
	Multi string `protobuf:"bytes,1,opt,name=multi" json:"multi,omitempty"`
}

func (m *MultiAddress) Reset()                    { *m = MultiAddress{} }
func (m *MultiAddress) String() string            { return proto.CompactTextString(m) }
func (*MultiAddress) ProtoMessage()               {}
func (*MultiAddress) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *MultiAddress) GetMulti() string {
	if m != nil {
		return m.Multi
	}
	return ""
}

// MultiAddresses are public multiaddress of multiple Nodes in the overlay
// network.
type MultiAddresses struct {
	Multis []*MultiAddress `protobuf:"bytes,1,rep,name=multis" json:"multis,omitempty"`
}

func (m *MultiAddresses) Reset()                    { *m = MultiAddresses{} }
func (m *MultiAddresses) String() string            { return proto.CompactTextString(m) }
func (*MultiAddresses) ProtoMessage()               {}
func (*MultiAddresses) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *MultiAddresses) GetMultis() []*MultiAddress {
	if m != nil {
		return m.Multis
	}
	return nil
}

// Nothing is in this message. It is used to send nothing, or signal a
// successful response.
type Nothing struct {
}

func (m *Nothing) Reset()                    { *m = Nothing{} }
func (m *Nothing) String() string            { return proto.CompactTextString(m) }
func (*Nothing) ProtoMessage()               {}
func (*Nothing) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func init() {
	proto.RegisterType((*OrderFragment)(nil), "rpc.OrderFragment")
	proto.RegisterType((*ResultFragment)(nil), "rpc.ResultFragment")
	proto.RegisterType((*Address)(nil), "rpc.Address")
	proto.RegisterType((*MultiAddress)(nil), "rpc.MultiAddress")
	proto.RegisterType((*MultiAddresses)(nil), "rpc.MultiAddresses")
	proto.RegisterType((*Nothing)(nil), "rpc.Nothing")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Node service

type NodeClient interface {
	// Ping the connection and swap MultiAddresses.
	Ping(ctx context.Context, in *MultiAddress, opts ...grpc.CallOption) (*MultiAddress, error)
	// Get all peers connected to the Node.
	Peers(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*MultiAddresses, error)
	// Send an OrderFragment to some target Node.
	SendOrderFragment(ctx context.Context, in *OrderFragment, opts ...grpc.CallOption) (*MultiAddress, error)
	// Send a ResultFragment to some target Node, where the ResultFragment is the
	// result of a computation on two OrderFragments.
	SendResultFragment(ctx context.Context, in *OrderFragment, opts ...grpc.CallOption) (*MultiAddress, error)
}

type nodeClient struct {
	cc *grpc.ClientConn
}

func NewNodeClient(cc *grpc.ClientConn) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) Ping(ctx context.Context, in *MultiAddress, opts ...grpc.CallOption) (*MultiAddress, error) {
	out := new(MultiAddress)
	err := grpc.Invoke(ctx, "/rpc.Node/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Peers(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*MultiAddresses, error) {
	out := new(MultiAddresses)
	err := grpc.Invoke(ctx, "/rpc.Node/Peers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) SendOrderFragment(ctx context.Context, in *OrderFragment, opts ...grpc.CallOption) (*MultiAddress, error) {
	out := new(MultiAddress)
	err := grpc.Invoke(ctx, "/rpc.Node/SendOrderFragment", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) SendResultFragment(ctx context.Context, in *OrderFragment, opts ...grpc.CallOption) (*MultiAddress, error) {
	out := new(MultiAddress)
	err := grpc.Invoke(ctx, "/rpc.Node/SendResultFragment", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Node service

type NodeServer interface {
	// Ping the connection and swap MultiAddresses.
	Ping(context.Context, *MultiAddress) (*MultiAddress, error)
	// Get all peers connected to the Node.
	Peers(context.Context, *Nothing) (*MultiAddresses, error)
	// Send an OrderFragment to some target Node.
	SendOrderFragment(context.Context, *OrderFragment) (*MultiAddress, error)
	// Send a ResultFragment to some target Node, where the ResultFragment is the
	// result of a computation on two OrderFragments.
	SendResultFragment(context.Context, *OrderFragment) (*MultiAddress, error)
}

func RegisterNodeServer(s *grpc.Server, srv NodeServer) {
	s.RegisterService(&_Node_serviceDesc, srv)
}

func _Node_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MultiAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Node/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Ping(ctx, req.(*MultiAddress))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Peers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Peers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Node/Peers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Peers(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_SendOrderFragment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderFragment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).SendOrderFragment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Node/SendOrderFragment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).SendOrderFragment(ctx, req.(*OrderFragment))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_SendResultFragment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderFragment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).SendResultFragment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Node/SendResultFragment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).SendResultFragment(ctx, req.(*OrderFragment))
	}
	return interceptor(ctx, in, info, handler)
}

var _Node_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Node_Ping_Handler,
		},
		{
			MethodName: "Peers",
			Handler:    _Node_Peers_Handler,
		},
		{
			MethodName: "SendOrderFragment",
			Handler:    _Node_SendOrderFragment_Handler,
		},
		{
			MethodName: "SendResultFragment",
			Handler:    _Node_SendResultFragment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 455 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x94, 0xdf, 0x6a, 0xd4, 0x40,
	0x14, 0x87, 0x49, 0xb2, 0x7f, 0x9a, 0x93, 0xb8, 0xd0, 0x53, 0x2f, 0x42, 0x29, 0x25, 0x44, 0x91,
	0x55, 0x64, 0x91, 0xf5, 0xce, 0x5e, 0xa9, 0x20, 0xf4, 0xc2, 0xb6, 0x64, 0xc5, 0xfb, 0xed, 0xce,
	0x74, 0x1b, 0x48, 0x32, 0x61, 0x66, 0x02, 0xee, 0xbb, 0xf9, 0x22, 0x3e, 0x87, 0x2f, 0x20, 0x73,
	0x92, 0xd8, 0xc9, 0x6e, 0xc0, 0xbd, 0xf1, 0x2e, 0xe7, 0x3b, 0xdf, 0x99, 0x09, 0xf3, 0x1b, 0x06,
	0x7c, 0x59, 0x6d, 0x16, 0x95, 0x14, 0x5a, 0xa0, 0x27, 0xab, 0x4d, 0xf2, 0xdb, 0x85, 0x67, 0xb7,
	0x92, 0x71, 0xf9, 0x45, 0xae, 0xb7, 0x05, 0x2f, 0x35, 0x5e, 0x80, 0xab, 0x45, 0xe4, 0xc4, 0xce,
	0x3c, 0x58, 0x86, 0x0b, 0xa3, 0x7f, 0x64, 0x4c, 0x72, 0xa5, 0x52, 0x57, 0x0b, 0x8c, 0x61, 0xf4,
	0x20, 0x45, 0x11, 0xb9, 0x03, 0x7d, 0xea, 0xe0, 0x0c, 0xdc, 0x8c, 0x45, 0x5e, 0xec, 0xcc, 0xc3,
	0xd4, 0xcd, 0x18, 0x46, 0x30, 0x15, 0x66, 0x83, 0x6b, 0x16, 0x8d, 0x08, 0x76, 0x25, 0x5e, 0x80,
	0x4f, 0x9f, 0xdf, 0x76, 0x15, 0x8f, 0xc6, 0xb1, 0x33, 0xf7, 0xd2, 0x27, 0x80, 0x09, 0x84, 0x54,
	0x7c, 0xaa, 0x77, 0x2b, 0x9e, 0xe7, 0xd1, 0x84, 0x84, 0x1e, 0x33, 0xce, 0x83, 0xd2, 0x9f, 0x05,
	0xe3, 0xab, 0xc7, 0xb5, 0xe4, 0xd1, 0x94, 0x36, 0xe8, 0x31, 0xe3, 0xa8, 0x92, 0x3d, 0x39, 0x27,
	0x8d, 0x63, 0x33, 0xbc, 0x04, 0xa8, 0x64, 0xb6, 0x69, 0x0d, 0x9f, 0x0c, 0x8b, 0xe0, 0x2b, 0x98,
	0x15, 0xeb, 0x1f, 0xdf, 0x45, 0x5e, 0x17, 0xad, 0x03, 0xe4, 0xec, 0x51, 0xf2, 0xb2, 0xd2, 0xf6,
	0x82, 0xd6, 0xeb, 0xd1, 0xe4, 0xa7, 0x07, 0xb3, 0x94, 0xab, 0x3a, 0xd7, 0xff, 0xed, 0xd8, 0x2f,
	0x01, 0xee, 0xeb, 0xdd, 0x6d, 0xef, 0xe4, 0x2d, 0x82, 0x31, 0x04, 0x8a, 0xe7, 0x79, 0x27, 0x8c,
	0x49, 0xb0, 0x11, 0x2e, 0x00, 0x3b, 0xbf, 0xfb, 0xcb, 0x6b, 0x46, 0x31, 0x84, 0xe9, 0x40, 0x07,
	0xdf, 0xc1, 0xd9, 0xdf, 0x71, 0x6b, 0xa0, 0xc9, 0x64, 0xa8, 0x75, 0x10, 0xdf, 0xc9, 0x11, 0xf1,
	0xf9, 0xff, 0x8c, 0x0f, 0x8e, 0x88, 0x2f, 0x38, 0x32, 0xbe, 0x70, 0x30, 0xbe, 0x17, 0x30, 0x6d,
	0x0f, 0xdf, 0xdc, 0xee, 0x75, 0xf3, 0x49, 0xd9, 0xf9, 0x69, 0x57, 0x26, 0x2f, 0x21, 0xfc, 0x5a,
	0xe7, 0x3a, 0xeb, 0xcc, 0xe7, 0x30, 0x2e, 0x4c, 0xdd, 0x7a, 0x4d, 0x91, 0x5c, 0xc1, 0xcc, 0xb6,
	0xb8, 0xc2, 0xd7, 0x30, 0xa1, 0x96, 0x59, 0xd0, 0x9b, 0x07, 0xcb, 0x53, 0x0a, 0xdb, 0x96, 0xd2,
	0x56, 0x48, 0x7c, 0x98, 0xde, 0x08, 0xfd, 0x98, 0x95, 0xdb, 0xe5, 0x2f, 0x07, 0x46, 0x37, 0x82,
	0x71, 0x7c, 0x0b, 0xa3, 0xbb, 0xac, 0xdc, 0xe2, 0xe1, 0xd8, 0xf9, 0x21, 0xc2, 0x37, 0x30, 0xbe,
	0xe3, 0x5c, 0x2a, 0x6c, 0xae, 0x54, 0xbb, 0xda, 0xf9, 0xd9, 0x81, 0xc9, 0x15, 0x7e, 0x80, 0xd3,
	0x15, 0x2f, 0x59, 0xff, 0xb5, 0x40, 0x32, 0x7b, 0x6c, 0x68, 0x9f, 0x2b, 0x40, 0x33, 0xbb, 0x77,
	0xe7, 0x8f, 0x1b, 0xbe, 0x9f, 0xd0, 0x7b, 0xf5, 0xfe, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x13,
	0x83, 0xba, 0x7f, 0xbc, 0x04, 0x00, 0x00,
}
