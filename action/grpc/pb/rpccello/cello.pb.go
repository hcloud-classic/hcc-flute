// Copyright 2020 by codex.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v4.0.0
// source: cello.proto

package rpccello

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	rpcmsgType "hcc/flute/action/grpc/pb/rpcmsgType"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Symbols defined in public import of msgType.proto.

type Empty = rpcmsgType.Empty
type Node = rpcmsgType.Node
type NodeDetail = rpcmsgType.NodeDetail
type Server = rpcmsgType.Server
type ServerNode = rpcmsgType.ServerNode
type Quota = rpcmsgType.Quota
type VNC = rpcmsgType.VNC
type Volume = rpcmsgType.Volume
type VolumeAttachment = rpcmsgType.VolumeAttachment
type AdaptiveIPSetting = rpcmsgType.AdaptiveIPSetting
type AdaptiveIPAvailableIPList = rpcmsgType.AdaptiveIPAvailableIPList
type AdaptiveIPServer = rpcmsgType.AdaptiveIPServer
type Subnet = rpcmsgType.Subnet
type Series = rpcmsgType.Series
type MetricInfo = rpcmsgType.MetricInfo
type MonitoringData = rpcmsgType.MonitoringData
type NormalAction = rpcmsgType.NormalAction
type HccAction = rpcmsgType.HccAction
type Action = rpcmsgType.Action
type Control = rpcmsgType.Control
type Controls = rpcmsgType.Controls

type ReqCello struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cello int32 `protobuf:"varint,1,opt,name=cello,proto3" json:"cello,omitempty"`
}

func (x *ReqCello) Reset() {
	*x = ReqCello{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cello_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqCello) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqCello) ProtoMessage() {}

func (x *ReqCello) ProtoReflect() protoreflect.Message {
	mi := &file_cello_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqCello.ProtoReflect.Descriptor instead.
func (*ReqCello) Descriptor() ([]byte, []int) {
	return file_cello_proto_rawDescGZIP(), []int{0}
}

func (x *ReqCello) GetCello() int32 {
	if x != nil {
		return x.Cello
	}
	return 0
}

type ResCello struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cello int32 `protobuf:"varint,2,opt,name=cello,proto3" json:"cello,omitempty"`
}

func (x *ResCello) Reset() {
	*x = ResCello{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cello_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResCello) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResCello) ProtoMessage() {}

func (x *ResCello) ProtoReflect() protoreflect.Message {
	mi := &file_cello_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResCello.ProtoReflect.Descriptor instead.
func (*ResCello) Descriptor() ([]byte, []int) {
	return file_cello_proto_rawDescGZIP(), []int{1}
}

func (x *ResCello) GetCello() int32 {
	if x != nil {
		return x.Cello
	}
	return 0
}

var File_cello_proto protoreflect.FileDescriptor

var file_cello_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x52,
	0x70, 0x63, 0x43, 0x65, 0x6c, 0x6c, 0x6f, 0x1a, 0x0d, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x20, 0x0a, 0x08, 0x52, 0x65, 0x71, 0x43, 0x65, 0x6c,
	0x6c, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x22, 0x20, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x43,
	0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x32, 0x3a, 0x0a, 0x05, 0x43, 0x65,
	0x6c, 0x6c, 0x6f, 0x12, 0x31, 0x0a, 0x05, 0x43, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x12, 0x2e, 0x52,
	0x70, 0x63, 0x43, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x52, 0x65, 0x71, 0x43, 0x65, 0x6c, 0x6c, 0x6f,
	0x1a, 0x12, 0x2e, 0x52, 0x70, 0x63, 0x43, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x43,
	0x65, 0x6c, 0x6c, 0x6f, 0x22, 0x00, 0x42, 0x23, 0x5a, 0x21, 0x68, 0x63, 0x63, 0x2f, 0x66, 0x6c,
	0x75, 0x74, 0x65, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f,
	0x70, 0x62, 0x2f, 0x72, 0x70, 0x63, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x50, 0x00, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cello_proto_rawDescOnce sync.Once
	file_cello_proto_rawDescData = file_cello_proto_rawDesc
)

func file_cello_proto_rawDescGZIP() []byte {
	file_cello_proto_rawDescOnce.Do(func() {
		file_cello_proto_rawDescData = protoimpl.X.CompressGZIP(file_cello_proto_rawDescData)
	})
	return file_cello_proto_rawDescData
}

var file_cello_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_cello_proto_goTypes = []interface{}{
	(*ReqCello)(nil), // 0: RpcCello.ReqCello
	(*ResCello)(nil), // 1: RpcCello.ResCello
}
var file_cello_proto_depIdxs = []int32{
	0, // 0: RpcCello.Cello.Cello:input_type -> RpcCello.ReqCello
	1, // 1: RpcCello.Cello.Cello:output_type -> RpcCello.ResCello
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_cello_proto_init() }
func file_cello_proto_init() {
	if File_cello_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cello_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqCello); i {
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
		file_cello_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResCello); i {
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
			RawDescriptor: file_cello_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cello_proto_goTypes,
		DependencyIndexes: file_cello_proto_depIdxs,
		MessageInfos:      file_cello_proto_msgTypes,
	}.Build()
	File_cello_proto = out.File
	file_cello_proto_rawDesc = nil
	file_cello_proto_goTypes = nil
	file_cello_proto_depIdxs = nil
}
