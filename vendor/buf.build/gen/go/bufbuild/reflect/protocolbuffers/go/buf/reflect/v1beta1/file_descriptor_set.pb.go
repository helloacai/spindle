// Copyright 2022-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: buf/reflect/v1beta1/file_descriptor_set.proto

package reflectv1beta1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetFileDescriptorSetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the module that contains the schema of interest.
	//
	// Some servers may host multiple modules and thus require this field. Others may host a
	// single module and not support this field. The format of the module name depends on the
	// server implementation.
	//
	// For Buf Schema Registries, the module name is required. An "Invalid Argument" error
	// will occur if it is missing. Buf Schema Registries require the module name to be in
	// the following format (note that the domain name of the registry must be included):
	//
	//	buf.build/acme/weather
	//
	// If the given module is not known to the server, a "Not Found" error is returned. If
	// a module name is given but not supported by this server or if the module name is in
	// an incorrect format, an "Invalid Argument" error is returned.
	Module string `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	// The version of the module to use.
	//
	// Some servers may not support multiple versions and thus not support this field. If
	// the field is supported by the server but not provided by the client, the server will
	// respond with the latest version of the requested module and indicate the version in
	// the response. The format of the module version depends on the server implementation.
	//
	// For Buf Schema Registries, this field can be a commit. But it can also be a tag, a
	// draft name, or "main" (which is the same as omitting it, since it will also resolve
	// to the latest version).
	//
	// If specified but the requested module has no such version, a "Not Found" error is
	// returned.
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// Zero or more symbol names. The names may refer to packages, messages, enums,
	// services, methods, or extensions. All names must be fully-qualified but should
	// NOT start with a period. If any name is invalid, the request will fail with an
	// "Invalid Argument" error. If any name is unresolvable/not present in the
	// requested module, the request will fail with a "Failed Precondition" error.
	//
	// If no names are provided, the full schema for the module is returned.
	// Otherwise, the resulting schema contains only the named elements and all of
	// their dependencies. This is enough information for the caller to construct
	// a dynamic message for any requested message types or to dynamically invoke
	// an RPC for any requested methods or services. If a package is named, that is
	// equivalent to specifying the names of all messages, enums, extensions, and
	// services defined in that package.
	Symbols []string `protobuf:"bytes,3,rep,name=symbols,proto3" json:"symbols,omitempty"`
}

func (x *GetFileDescriptorSetRequest) Reset() {
	*x = GetFileDescriptorSetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_reflect_v1beta1_file_descriptor_set_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileDescriptorSetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileDescriptorSetRequest) ProtoMessage() {}

func (x *GetFileDescriptorSetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_reflect_v1beta1_file_descriptor_set_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileDescriptorSetRequest.ProtoReflect.Descriptor instead.
func (*GetFileDescriptorSetRequest) Descriptor() ([]byte, []int) {
	return file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDescGZIP(), []int{0}
}

func (x *GetFileDescriptorSetRequest) GetModule() string {
	if x != nil {
		return x.Module
	}
	return ""
}

func (x *GetFileDescriptorSetRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *GetFileDescriptorSetRequest) GetSymbols() []string {
	if x != nil {
		return x.Symbols
	}
	return nil
}

type GetFileDescriptorSetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The schema, which is a set of file descriptors that include the requested symbols
	// and their dependencies.
	//
	// The returned file descriptors will be topologically sorted.
	FileDescriptorSet *descriptorpb.FileDescriptorSet `protobuf:"bytes,1,opt,name=file_descriptor_set,json=fileDescriptorSet,proto3" json:"file_descriptor_set,omitempty"`
	// The version of the returned schema. May not be set, such as if the server does not
	// support multiple versions of schemas. May be different from the requested version,
	// such as if the requested version was a name or tag that is resolved to another form.
	//
	// For Buf Schema Registries, if the requested version is a tag, draft name, or "main",
	// the returned version will be the corresponding commit.
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetFileDescriptorSetResponse) Reset() {
	*x = GetFileDescriptorSetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_reflect_v1beta1_file_descriptor_set_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileDescriptorSetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileDescriptorSetResponse) ProtoMessage() {}

func (x *GetFileDescriptorSetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_reflect_v1beta1_file_descriptor_set_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileDescriptorSetResponse.ProtoReflect.Descriptor instead.
func (*GetFileDescriptorSetResponse) Descriptor() ([]byte, []int) {
	return file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDescGZIP(), []int{1}
}

func (x *GetFileDescriptorSetResponse) GetFileDescriptorSet() *descriptorpb.FileDescriptorSet {
	if x != nil {
		return x.FileDescriptorSet
	}
	return nil
}

func (x *GetFileDescriptorSetResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_buf_reflect_v1beta1_file_descriptor_set_proto protoreflect.FileDescriptor

var file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65, 0x66, 0x6c, 0x65, 0x63, 0x74, 0x2f, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x13, 0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x66, 0x6c, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x69, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x79, 0x6d, 0x62, 0x6f,
	0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c,
	0x73, 0x22, 0x8c, 0x01, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x52, 0x0a, 0x13, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x53, 0x65, 0x74, 0x52, 0x11, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x32, 0x9d, 0x01, 0x0a, 0x18, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x80, 0x01,
	0x0a, 0x14, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x12, 0x30, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x66,
	0x6c, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x53, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x72,
	0x65, 0x66, 0x6c, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x03, 0x90, 0x02, 0x01,
	0x42, 0x59, 0x5a, 0x57, 0x62, 0x75, 0x66, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x75, 0x66, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x72, 0x65,
	0x66, 0x6c, 0x65, 0x63, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75,
	0x66, 0x66, 0x65, 0x72, 0x73, 0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65, 0x66,
	0x6c, 0x65, 0x63, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x3b, 0x72, 0x65, 0x66,
	0x6c, 0x65, 0x63, 0x74, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDescOnce sync.Once
	file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDescData = file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDesc
)

func file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDescGZIP() []byte {
	file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDescOnce.Do(func() {
		file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDescData = protoimpl.X.CompressGZIP(file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDescData)
	})
	return file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDescData
}

var file_buf_reflect_v1beta1_file_descriptor_set_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_buf_reflect_v1beta1_file_descriptor_set_proto_goTypes = []interface{}{
	(*GetFileDescriptorSetRequest)(nil),    // 0: buf.reflect.v1beta1.GetFileDescriptorSetRequest
	(*GetFileDescriptorSetResponse)(nil),   // 1: buf.reflect.v1beta1.GetFileDescriptorSetResponse
	(*descriptorpb.FileDescriptorSet)(nil), // 2: google.protobuf.FileDescriptorSet
}
var file_buf_reflect_v1beta1_file_descriptor_set_proto_depIdxs = []int32{
	2, // 0: buf.reflect.v1beta1.GetFileDescriptorSetResponse.file_descriptor_set:type_name -> google.protobuf.FileDescriptorSet
	0, // 1: buf.reflect.v1beta1.FileDescriptorSetService.GetFileDescriptorSet:input_type -> buf.reflect.v1beta1.GetFileDescriptorSetRequest
	1, // 2: buf.reflect.v1beta1.FileDescriptorSetService.GetFileDescriptorSet:output_type -> buf.reflect.v1beta1.GetFileDescriptorSetResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_buf_reflect_v1beta1_file_descriptor_set_proto_init() }
func file_buf_reflect_v1beta1_file_descriptor_set_proto_init() {
	if File_buf_reflect_v1beta1_file_descriptor_set_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_buf_reflect_v1beta1_file_descriptor_set_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileDescriptorSetRequest); i {
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
		file_buf_reflect_v1beta1_file_descriptor_set_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileDescriptorSetResponse); i {
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
			RawDescriptor: file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_buf_reflect_v1beta1_file_descriptor_set_proto_goTypes,
		DependencyIndexes: file_buf_reflect_v1beta1_file_descriptor_set_proto_depIdxs,
		MessageInfos:      file_buf_reflect_v1beta1_file_descriptor_set_proto_msgTypes,
	}.Build()
	File_buf_reflect_v1beta1_file_descriptor_set_proto = out.File
	file_buf_reflect_v1beta1_file_descriptor_set_proto_rawDesc = nil
	file_buf_reflect_v1beta1_file_descriptor_set_proto_goTypes = nil
	file_buf_reflect_v1beta1_file_descriptor_set_proto_depIdxs = nil
}
