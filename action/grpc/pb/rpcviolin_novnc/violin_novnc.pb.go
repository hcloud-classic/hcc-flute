// Copyright 2015 gRPC authors.
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
// 	protoc-gen-go v1.25.0
// 	protoc        v4.0.0
// source: violin_novnc.proto

package rpcviolin_novnc

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
type HccError = rpcmsgType.HccError
type Node = rpcmsgType.Node
type NodeDetail = rpcmsgType.NodeDetail
type Server = rpcmsgType.Server
type ServerNode = rpcmsgType.ServerNode
type Quota = rpcmsgType.Quota
type VNC = rpcmsgType.VNC
type Volume = rpcmsgType.Volume
type VolumeAttachment = rpcmsgType.VolumeAttachment
type Pool = rpcmsgType.Pool
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
type ScheduledNodes = rpcmsgType.ScheduledNodes
type ScheduleServer = rpcmsgType.ScheduleServer
type ServerAction = rpcmsgType.ServerAction

type ReqNoVNC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vncs []*rpcmsgType.VNC `protobuf:"bytes,1,rep,name=vncs,proto3" json:"vncs,omitempty"`
}

func (x *ReqNoVNC) Reset() {
	*x = ReqNoVNC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_violin_novnc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqNoVNC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqNoVNC) ProtoMessage() {}

func (x *ReqNoVNC) ProtoReflect() protoreflect.Message {
	mi := &file_violin_novnc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqNoVNC.ProtoReflect.Descriptor instead.
func (*ReqNoVNC) Descriptor() ([]byte, []int) {
	return file_violin_novnc_proto_rawDescGZIP(), []int{0}
}

func (x *ReqNoVNC) GetVncs() []*rpcmsgType.VNC {
	if x != nil {
		return x.Vncs
	}
	return nil
}

type ResNoVNC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vncs          []*rpcmsgType.VNC      `protobuf:"bytes,1,rep,name=vncs,proto3" json:"vncs,omitempty"`
	HccErrorStack []*rpcmsgType.HccError `protobuf:"bytes,2,rep,name=hccErrorStack,proto3" json:"hccErrorStack,omitempty"`
}

func (x *ResNoVNC) Reset() {
	*x = ResNoVNC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_violin_novnc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResNoVNC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResNoVNC) ProtoMessage() {}

func (x *ResNoVNC) ProtoReflect() protoreflect.Message {
	mi := &file_violin_novnc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResNoVNC.ProtoReflect.Descriptor instead.
func (*ResNoVNC) Descriptor() ([]byte, []int) {
	return file_violin_novnc_proto_rawDescGZIP(), []int{1}
}

func (x *ResNoVNC) GetVncs() []*rpcmsgType.VNC {
	if x != nil {
		return x.Vncs
	}
	return nil
}

func (x *ResNoVNC) GetHccErrorStack() []*rpcmsgType.HccError {
	if x != nil {
		return x.HccErrorStack
	}
	return nil
}

type ReqControlVNC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vnc *rpcmsgType.VNC `protobuf:"bytes,1,opt,name=vnc,proto3" json:"vnc,omitempty"`
}

func (x *ReqControlVNC) Reset() {
	*x = ReqControlVNC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_violin_novnc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqControlVNC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqControlVNC) ProtoMessage() {}

func (x *ReqControlVNC) ProtoReflect() protoreflect.Message {
	mi := &file_violin_novnc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqControlVNC.ProtoReflect.Descriptor instead.
func (*ReqControlVNC) Descriptor() ([]byte, []int) {
	return file_violin_novnc_proto_rawDescGZIP(), []int{2}
}

func (x *ReqControlVNC) GetVnc() *rpcmsgType.VNC {
	if x != nil {
		return x.Vnc
	}
	return nil
}

type ResControlVNC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Port          string                 `protobuf:"bytes,1,opt,name=port,proto3" json:"port,omitempty"`
	HccErrorStack []*rpcmsgType.HccError `protobuf:"bytes,2,rep,name=hccErrorStack,proto3" json:"hccErrorStack,omitempty"`
}

func (x *ResControlVNC) Reset() {
	*x = ResControlVNC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_violin_novnc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResControlVNC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResControlVNC) ProtoMessage() {}

func (x *ResControlVNC) ProtoReflect() protoreflect.Message {
	mi := &file_violin_novnc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResControlVNC.ProtoReflect.Descriptor instead.
func (*ResControlVNC) Descriptor() ([]byte, []int) {
	return file_violin_novnc_proto_rawDescGZIP(), []int{3}
}

func (x *ResControlVNC) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

func (x *ResControlVNC) GetHccErrorStack() []*rpcmsgType.HccError {
	if x != nil {
		return x.HccErrorStack
	}
	return nil
}

var File_violin_novnc_proto protoreflect.FileDescriptor

var file_violin_novnc_proto_rawDesc = []byte{
	0x0a, 0x12, 0x76, 0x69, 0x6f, 0x6c, 0x69, 0x6e, 0x5f, 0x6e, 0x6f, 0x76, 0x6e, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x52, 0x70, 0x63, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x1a, 0x0d,
	0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2c, 0x0a,
	0x08, 0x52, 0x65, 0x71, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x12, 0x20, 0x0a, 0x04, 0x76, 0x6e, 0x63,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70,
	0x65, 0x2e, 0x56, 0x4e, 0x43, 0x52, 0x04, 0x76, 0x6e, 0x63, 0x73, 0x22, 0x65, 0x0a, 0x08, 0x52,
	0x65, 0x73, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x12, 0x20, 0x0a, 0x04, 0x76, 0x6e, 0x63, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x2e,
	0x56, 0x4e, 0x43, 0x52, 0x04, 0x76, 0x6e, 0x63, 0x73, 0x12, 0x37, 0x0a, 0x0d, 0x68, 0x63, 0x63,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x2e, 0x48, 0x63, 0x63, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x52, 0x0d, 0x68, 0x63, 0x63, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x53, 0x74, 0x61,
	0x63, 0x6b, 0x22, 0x2f, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x56, 0x4e, 0x43, 0x12, 0x1e, 0x0a, 0x03, 0x76, 0x6e, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x2e, 0x56, 0x4e, 0x43, 0x52, 0x03,
	0x76, 0x6e, 0x63, 0x22, 0x5c, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x56, 0x4e, 0x43, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x37, 0x0a, 0x0d, 0x68, 0x63, 0x63, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x2e, 0x48, 0x63, 0x63, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x52, 0x0d, 0x68, 0x63, 0x63, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x63,
	0x6b, 0x32, 0x80, 0x01, 0x0a, 0x05, 0x6e, 0x6f, 0x76, 0x6e, 0x63, 0x12, 0x35, 0x0a, 0x09, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x56, 0x4e, 0x43, 0x12, 0x12, 0x2e, 0x52, 0x70, 0x63, 0x4e, 0x6f,
	0x56, 0x4e, 0x43, 0x2e, 0x52, 0x65, 0x71, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x1a, 0x12, 0x2e, 0x52,
	0x70, 0x63, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x2e, 0x52, 0x65, 0x73, 0x4e, 0x6f, 0x56, 0x4e, 0x43,
	0x22, 0x00, 0x12, 0x40, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x56, 0x4e, 0x43,
	0x12, 0x17, 0x2e, 0x52, 0x70, 0x63, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x2e, 0x52, 0x65, 0x71, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x56, 0x4e, 0x43, 0x1a, 0x17, 0x2e, 0x52, 0x70, 0x63, 0x4e,
	0x6f, 0x56, 0x4e, 0x43, 0x2e, 0x52, 0x65, 0x73, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x56,
	0x4e, 0x43, 0x22, 0x00, 0x42, 0x2a, 0x5a, 0x28, 0x68, 0x63, 0x63, 0x2f, 0x66, 0x6c, 0x75, 0x74,
	0x65, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62,
	0x2f, 0x72, 0x70, 0x63, 0x76, 0x69, 0x6f, 0x6c, 0x69, 0x6e, 0x5f, 0x6e, 0x6f, 0x76, 0x6e, 0x63,
	0x50, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_violin_novnc_proto_rawDescOnce sync.Once
	file_violin_novnc_proto_rawDescData = file_violin_novnc_proto_rawDesc
)

func file_violin_novnc_proto_rawDescGZIP() []byte {
	file_violin_novnc_proto_rawDescOnce.Do(func() {
		file_violin_novnc_proto_rawDescData = protoimpl.X.CompressGZIP(file_violin_novnc_proto_rawDescData)
	})
	return file_violin_novnc_proto_rawDescData
}

var file_violin_novnc_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_violin_novnc_proto_goTypes = []interface{}{
	(*ReqNoVNC)(nil),            // 0: RpcNoVNC.ReqNoVNC
	(*ResNoVNC)(nil),            // 1: RpcNoVNC.ResNoVNC
	(*ReqControlVNC)(nil),       // 2: RpcNoVNC.ReqControlVNC
	(*ResControlVNC)(nil),       // 3: RpcNoVNC.ResControlVNC
	(*rpcmsgType.VNC)(nil),      // 4: MsgType.VNC
	(*rpcmsgType.HccError)(nil), // 5: MsgType.HccError
}
var file_violin_novnc_proto_depIdxs = []int32{
	4, // 0: RpcNoVNC.ReqNoVNC.vncs:type_name -> MsgType.VNC
	4, // 1: RpcNoVNC.ResNoVNC.vncs:type_name -> MsgType.VNC
	5, // 2: RpcNoVNC.ResNoVNC.hccErrorStack:type_name -> MsgType.HccError
	4, // 3: RpcNoVNC.ReqControlVNC.vnc:type_name -> MsgType.VNC
	5, // 4: RpcNoVNC.ResControlVNC.hccErrorStack:type_name -> MsgType.HccError
	0, // 5: RpcNoVNC.novnc.CreateVNC:input_type -> RpcNoVNC.ReqNoVNC
	2, // 6: RpcNoVNC.novnc.ControlVNC:input_type -> RpcNoVNC.ReqControlVNC
	1, // 7: RpcNoVNC.novnc.CreateVNC:output_type -> RpcNoVNC.ResNoVNC
	3, // 8: RpcNoVNC.novnc.ControlVNC:output_type -> RpcNoVNC.ResControlVNC
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_violin_novnc_proto_init() }
func file_violin_novnc_proto_init() {
	if File_violin_novnc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_violin_novnc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqNoVNC); i {
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
		file_violin_novnc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResNoVNC); i {
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
		file_violin_novnc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqControlVNC); i {
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
		file_violin_novnc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResControlVNC); i {
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
			RawDescriptor: file_violin_novnc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_violin_novnc_proto_goTypes,
		DependencyIndexes: file_violin_novnc_proto_depIdxs,
		MessageInfos:      file_violin_novnc_proto_msgTypes,
	}.Build()
	File_violin_novnc_proto = out.File
	file_violin_novnc_proto_rawDesc = nil
	file_violin_novnc_proto_goTypes = nil
	file_violin_novnc_proto_depIdxs = nil
}
