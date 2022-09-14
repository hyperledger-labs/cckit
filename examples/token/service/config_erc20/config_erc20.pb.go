// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: token/service/config_erc20/config_erc20.proto

package config_erc20

import (
	context "context"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *NameResponse) Reset() {
	*x = NameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_service_config_erc20_config_erc20_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NameResponse) ProtoMessage() {}

func (x *NameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_token_service_config_erc20_config_erc20_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NameResponse.ProtoReflect.Descriptor instead.
func (*NameResponse) Descriptor() ([]byte, []int) {
	return file_token_service_config_erc20_config_erc20_proto_rawDescGZIP(), []int{0}
}

func (x *NameResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type SymbolResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
}

func (x *SymbolResponse) Reset() {
	*x = SymbolResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_service_config_erc20_config_erc20_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SymbolResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SymbolResponse) ProtoMessage() {}

func (x *SymbolResponse) ProtoReflect() protoreflect.Message {
	mi := &file_token_service_config_erc20_config_erc20_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SymbolResponse.ProtoReflect.Descriptor instead.
func (*SymbolResponse) Descriptor() ([]byte, []int) {
	return file_token_service_config_erc20_config_erc20_proto_rawDescGZIP(), []int{1}
}

func (x *SymbolResponse) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

type DecimalsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Decimals uint32 `protobuf:"varint,1,opt,name=decimals,proto3" json:"decimals,omitempty"`
}

func (x *DecimalsResponse) Reset() {
	*x = DecimalsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_service_config_erc20_config_erc20_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DecimalsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecimalsResponse) ProtoMessage() {}

func (x *DecimalsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_token_service_config_erc20_config_erc20_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecimalsResponse.ProtoReflect.Descriptor instead.
func (*DecimalsResponse) Descriptor() ([]byte, []int) {
	return file_token_service_config_erc20_config_erc20_proto_rawDescGZIP(), []int{2}
}

func (x *DecimalsResponse) GetDecimals() uint32 {
	if x != nil {
		return x.Decimals
	}
	return 0
}

type TotalSupplyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalSupply uint64 `protobuf:"varint,1,opt,name=total_supply,json=totalSupply,proto3" json:"total_supply,omitempty"`
}

func (x *TotalSupplyResponse) Reset() {
	*x = TotalSupplyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_service_config_erc20_config_erc20_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TotalSupplyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TotalSupplyResponse) ProtoMessage() {}

func (x *TotalSupplyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_token_service_config_erc20_config_erc20_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TotalSupplyResponse.ProtoReflect.Descriptor instead.
func (*TotalSupplyResponse) Descriptor() ([]byte, []int) {
	return file_token_service_config_erc20_config_erc20_proto_rawDescGZIP(), []int{3}
}

func (x *TotalSupplyResponse) GetTotalSupply() uint64 {
	if x != nil {
		return x.TotalSupply
	}
	return 0
}

var File_token_service_config_erc20_config_erc20_proto protoreflect.FileDescriptor

var file_token_service_config_erc20_config_erc20_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x65, 0x72, 0x63, 0x32, 0x30, 0x2f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x5f, 0x65, 0x72, 0x63, 0x32, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x2b, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x65, 0x72, 0x63, 0x32, 0x30, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x65, 0x72, 0x63, 0x32, 0x30, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a, 0x0c, 0x4e, 0x61, 0x6d, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x28, 0x0a, 0x0e, 0x53,
	0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x22, 0x2e, 0x0a, 0x10, 0x44, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x63,
	0x69, 0x6d, 0x61, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x64, 0x65, 0x63,
	0x69, 0x6d, 0x61, 0x6c, 0x73, 0x22, 0x38, 0x0a, 0x13, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x75,
	0x70, 0x70, 0x6c, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x73, 0x75, 0x70, 0x70, 0x6c, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x75, 0x70, 0x70, 0x6c, 0x79, 0x32,
	0xf1, 0x03, 0x0a, 0x12, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x52, 0x43, 0x32, 0x30, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6b, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x39, 0x2e, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x73, 0x2e, 0x65, 0x72, 0x63, 0x32, 0x30, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x5f, 0x65, 0x72, 0x63, 0x32, 0x30, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x0d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x07, 0x12, 0x05, 0x2f, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x71, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x3b, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x73, 0x2e, 0x65, 0x72, 0x63, 0x32, 0x30, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x5f, 0x65, 0x72, 0x63, 0x32, 0x30, 0x2e, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x0f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x2f,
	0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x77, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x44, 0x65, 0x63,
	0x69, 0x6d, 0x61, 0x6c, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x3d, 0x2e,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x65, 0x72, 0x63, 0x32, 0x30, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x65, 0x72, 0x63, 0x32, 0x30, 0x2e, 0x44, 0x65, 0x63, 0x69,
	0x6d, 0x61, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0b, 0x12, 0x09, 0x2f, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x73, 0x12,
	0x81, 0x01, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x75, 0x70, 0x70,
	0x6c, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x40, 0x2e, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x65, 0x72, 0x63, 0x32, 0x30, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x5f, 0x65, 0x72, 0x63, 0x32, 0x30, 0x2e, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x75,
	0x70, 0x70, 0x6c, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x2d, 0x73, 0x75, 0x70,
	0x70, 0x6c, 0x79, 0x42, 0x47, 0x5a, 0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x68, 0x79, 0x70, 0x65, 0x72, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2d, 0x6c, 0x61,
	0x62, 0x73, 0x2f, 0x63, 0x63, 0x6b, 0x69, 0x74, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x73, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x65, 0x72, 0x63, 0x32, 0x30, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_token_service_config_erc20_config_erc20_proto_rawDescOnce sync.Once
	file_token_service_config_erc20_config_erc20_proto_rawDescData = file_token_service_config_erc20_config_erc20_proto_rawDesc
)

func file_token_service_config_erc20_config_erc20_proto_rawDescGZIP() []byte {
	file_token_service_config_erc20_config_erc20_proto_rawDescOnce.Do(func() {
		file_token_service_config_erc20_config_erc20_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_service_config_erc20_config_erc20_proto_rawDescData)
	})
	return file_token_service_config_erc20_config_erc20_proto_rawDescData
}

var file_token_service_config_erc20_config_erc20_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_token_service_config_erc20_config_erc20_proto_goTypes = []interface{}{
	(*NameResponse)(nil),        // 0: examples.erc20_service.service.config_erc20.NameResponse
	(*SymbolResponse)(nil),      // 1: examples.erc20_service.service.config_erc20.SymbolResponse
	(*DecimalsResponse)(nil),    // 2: examples.erc20_service.service.config_erc20.DecimalsResponse
	(*TotalSupplyResponse)(nil), // 3: examples.erc20_service.service.config_erc20.TotalSupplyResponse
	(*emptypb.Empty)(nil),       // 4: google.protobuf.Empty
}
var file_token_service_config_erc20_config_erc20_proto_depIdxs = []int32{
	4, // 0: examples.erc20_service.service.config_erc20.ConfigERC20Service.GetName:input_type -> google.protobuf.Empty
	4, // 1: examples.erc20_service.service.config_erc20.ConfigERC20Service.GetSymbol:input_type -> google.protobuf.Empty
	4, // 2: examples.erc20_service.service.config_erc20.ConfigERC20Service.GetDecimals:input_type -> google.protobuf.Empty
	4, // 3: examples.erc20_service.service.config_erc20.ConfigERC20Service.GetTotalSupply:input_type -> google.protobuf.Empty
	0, // 4: examples.erc20_service.service.config_erc20.ConfigERC20Service.GetName:output_type -> examples.erc20_service.service.config_erc20.NameResponse
	1, // 5: examples.erc20_service.service.config_erc20.ConfigERC20Service.GetSymbol:output_type -> examples.erc20_service.service.config_erc20.SymbolResponse
	2, // 6: examples.erc20_service.service.config_erc20.ConfigERC20Service.GetDecimals:output_type -> examples.erc20_service.service.config_erc20.DecimalsResponse
	3, // 7: examples.erc20_service.service.config_erc20.ConfigERC20Service.GetTotalSupply:output_type -> examples.erc20_service.service.config_erc20.TotalSupplyResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_token_service_config_erc20_config_erc20_proto_init() }
func file_token_service_config_erc20_config_erc20_proto_init() {
	if File_token_service_config_erc20_config_erc20_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_token_service_config_erc20_config_erc20_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NameResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_service_config_erc20_config_erc20_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SymbolResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_service_config_erc20_config_erc20_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DecimalsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_service_config_erc20_config_erc20_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TotalSupplyResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_token_service_config_erc20_config_erc20_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_token_service_config_erc20_config_erc20_proto_goTypes,
		DependencyIndexes: file_token_service_config_erc20_config_erc20_proto_depIdxs,
		MessageInfos:      file_token_service_config_erc20_config_erc20_proto_msgTypes,
	}.Build()
	File_token_service_config_erc20_config_erc20_proto = out.File
	file_token_service_config_erc20_config_erc20_proto_rawDesc = nil
	file_token_service_config_erc20_config_erc20_proto_goTypes = nil
	file_token_service_config_erc20_config_erc20_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ConfigERC20ServiceClient is the client API for ConfigERC20Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConfigERC20ServiceClient interface {
	// Returns the name of the token.
	GetName(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NameResponse, error)
	// Returns the symbol of the token, usually a shorter version of the name.
	GetSymbol(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SymbolResponse, error)
	// Returns the number of decimals used to get its user representation.
	// For example, if decimals equals 2, a balance of 505 tokens should be displayed to a user as 5,05 (505 / 10 ** 2).
	GetDecimals(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*DecimalsResponse, error)
	// Returns the amount of tokens in existence.
	GetTotalSupply(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TotalSupplyResponse, error)
}

type configERC20ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigERC20ServiceClient(cc grpc.ClientConnInterface) ConfigERC20ServiceClient {
	return &configERC20ServiceClient{cc}
}

func (c *configERC20ServiceClient) GetName(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NameResponse, error) {
	out := new(NameResponse)
	err := c.cc.Invoke(ctx, "/examples.erc20_service.service.config_erc20.ConfigERC20Service/GetName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configERC20ServiceClient) GetSymbol(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SymbolResponse, error) {
	out := new(SymbolResponse)
	err := c.cc.Invoke(ctx, "/examples.erc20_service.service.config_erc20.ConfigERC20Service/GetSymbol", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configERC20ServiceClient) GetDecimals(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*DecimalsResponse, error) {
	out := new(DecimalsResponse)
	err := c.cc.Invoke(ctx, "/examples.erc20_service.service.config_erc20.ConfigERC20Service/GetDecimals", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configERC20ServiceClient) GetTotalSupply(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TotalSupplyResponse, error) {
	out := new(TotalSupplyResponse)
	err := c.cc.Invoke(ctx, "/examples.erc20_service.service.config_erc20.ConfigERC20Service/GetTotalSupply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigERC20ServiceServer is the server API for ConfigERC20Service service.
type ConfigERC20ServiceServer interface {
	// Returns the name of the token.
	GetName(context.Context, *emptypb.Empty) (*NameResponse, error)
	// Returns the symbol of the token, usually a shorter version of the name.
	GetSymbol(context.Context, *emptypb.Empty) (*SymbolResponse, error)
	// Returns the number of decimals used to get its user representation.
	// For example, if decimals equals 2, a balance of 505 tokens should be displayed to a user as 5,05 (505 / 10 ** 2).
	GetDecimals(context.Context, *emptypb.Empty) (*DecimalsResponse, error)
	// Returns the amount of tokens in existence.
	GetTotalSupply(context.Context, *emptypb.Empty) (*TotalSupplyResponse, error)
}

// UnimplementedConfigERC20ServiceServer can be embedded to have forward compatible implementations.
type UnimplementedConfigERC20ServiceServer struct {
}

func (*UnimplementedConfigERC20ServiceServer) GetName(context.Context, *emptypb.Empty) (*NameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetName not implemented")
}
func (*UnimplementedConfigERC20ServiceServer) GetSymbol(context.Context, *emptypb.Empty) (*SymbolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSymbol not implemented")
}
func (*UnimplementedConfigERC20ServiceServer) GetDecimals(context.Context, *emptypb.Empty) (*DecimalsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDecimals not implemented")
}
func (*UnimplementedConfigERC20ServiceServer) GetTotalSupply(context.Context, *emptypb.Empty) (*TotalSupplyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTotalSupply not implemented")
}

func RegisterConfigERC20ServiceServer(s *grpc.Server, srv ConfigERC20ServiceServer) {
	s.RegisterService(&_ConfigERC20Service_serviceDesc, srv)
}

func _ConfigERC20Service_GetName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigERC20ServiceServer).GetName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.erc20_service.service.config_erc20.ConfigERC20Service/GetName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigERC20ServiceServer).GetName(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigERC20Service_GetSymbol_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigERC20ServiceServer).GetSymbol(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.erc20_service.service.config_erc20.ConfigERC20Service/GetSymbol",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigERC20ServiceServer).GetSymbol(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigERC20Service_GetDecimals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigERC20ServiceServer).GetDecimals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.erc20_service.service.config_erc20.ConfigERC20Service/GetDecimals",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigERC20ServiceServer).GetDecimals(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigERC20Service_GetTotalSupply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigERC20ServiceServer).GetTotalSupply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.erc20_service.service.config_erc20.ConfigERC20Service/GetTotalSupply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigERC20ServiceServer).GetTotalSupply(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _ConfigERC20Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "examples.erc20_service.service.config_erc20.ConfigERC20Service",
	HandlerType: (*ConfigERC20ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetName",
			Handler:    _ConfigERC20Service_GetName_Handler,
		},
		{
			MethodName: "GetSymbol",
			Handler:    _ConfigERC20Service_GetSymbol_Handler,
		},
		{
			MethodName: "GetDecimals",
			Handler:    _ConfigERC20Service_GetDecimals_Handler,
		},
		{
			MethodName: "GetTotalSupply",
			Handler:    _ConfigERC20Service_GetTotalSupply_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token/service/config_erc20/config_erc20.proto",
}
