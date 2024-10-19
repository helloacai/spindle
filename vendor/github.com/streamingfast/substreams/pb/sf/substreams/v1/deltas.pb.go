// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: sf/substreams/v1/deltas.proto

package pbsubstreams

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StoreDelta_Operation int32

const (
	StoreDelta_UNSET  StoreDelta_Operation = 0
	StoreDelta_CREATE StoreDelta_Operation = 1
	StoreDelta_UPDATE StoreDelta_Operation = 2
	StoreDelta_DELETE StoreDelta_Operation = 3
)

// Enum value maps for StoreDelta_Operation.
var (
	StoreDelta_Operation_name = map[int32]string{
		0: "UNSET",
		1: "CREATE",
		2: "UPDATE",
		3: "DELETE",
	}
	StoreDelta_Operation_value = map[string]int32{
		"UNSET":  0,
		"CREATE": 1,
		"UPDATE": 2,
		"DELETE": 3,
	}
)

func (x StoreDelta_Operation) Enum() *StoreDelta_Operation {
	p := new(StoreDelta_Operation)
	*p = x
	return p
}

func (x StoreDelta_Operation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StoreDelta_Operation) Descriptor() protoreflect.EnumDescriptor {
	return file_sf_substreams_v1_deltas_proto_enumTypes[0].Descriptor()
}

func (StoreDelta_Operation) Type() protoreflect.EnumType {
	return &file_sf_substreams_v1_deltas_proto_enumTypes[0]
}

func (x StoreDelta_Operation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StoreDelta_Operation.Descriptor instead.
func (StoreDelta_Operation) EnumDescriptor() ([]byte, []int) {
	return file_sf_substreams_v1_deltas_proto_rawDescGZIP(), []int{1, 0}
}

type StoreDeltas struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreDeltas []*StoreDelta `protobuf:"bytes,1,rep,name=store_deltas,json=storeDeltas,proto3" json:"store_deltas,omitempty"`
}

func (x *StoreDeltas) Reset() {
	*x = StoreDeltas{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sf_substreams_v1_deltas_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreDeltas) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreDeltas) ProtoMessage() {}

func (x *StoreDeltas) ProtoReflect() protoreflect.Message {
	mi := &file_sf_substreams_v1_deltas_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreDeltas.ProtoReflect.Descriptor instead.
func (*StoreDeltas) Descriptor() ([]byte, []int) {
	return file_sf_substreams_v1_deltas_proto_rawDescGZIP(), []int{0}
}

func (x *StoreDeltas) GetStoreDeltas() []*StoreDelta {
	if x != nil {
		return x.StoreDeltas
	}
	return nil
}

type StoreDelta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Operation StoreDelta_Operation `protobuf:"varint,1,opt,name=operation,proto3,enum=sf.substreams.v1.StoreDelta_Operation" json:"operation,omitempty"`
	Ordinal   uint64               `protobuf:"varint,2,opt,name=ordinal,proto3" json:"ordinal,omitempty"`
	Key       string               `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	OldValue  []byte               `protobuf:"bytes,4,opt,name=old_value,json=oldValue,proto3" json:"old_value,omitempty"`
	NewValue  []byte               `protobuf:"bytes,5,opt,name=new_value,json=newValue,proto3" json:"new_value,omitempty"`
}

func (x *StoreDelta) Reset() {
	*x = StoreDelta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sf_substreams_v1_deltas_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreDelta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreDelta) ProtoMessage() {}

func (x *StoreDelta) ProtoReflect() protoreflect.Message {
	mi := &file_sf_substreams_v1_deltas_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreDelta.ProtoReflect.Descriptor instead.
func (*StoreDelta) Descriptor() ([]byte, []int) {
	return file_sf_substreams_v1_deltas_proto_rawDescGZIP(), []int{1}
}

func (x *StoreDelta) GetOperation() StoreDelta_Operation {
	if x != nil {
		return x.Operation
	}
	return StoreDelta_UNSET
}

func (x *StoreDelta) GetOrdinal() uint64 {
	if x != nil {
		return x.Ordinal
	}
	return 0
}

func (x *StoreDelta) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *StoreDelta) GetOldValue() []byte {
	if x != nil {
		return x.OldValue
	}
	return nil
}

func (x *StoreDelta) GetNewValue() []byte {
	if x != nil {
		return x.NewValue
	}
	return nil
}

var File_sf_substreams_v1_deltas_proto protoreflect.FileDescriptor

var file_sf_substreams_v1_deltas_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x73, 0x66, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x2f,
	0x76, 0x31, 0x2f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x10, 0x73, 0x66, 0x2e, 0x73, 0x75, 0x62, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x76,
	0x31, 0x22, 0x4e, 0x0a, 0x0b, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x44, 0x65, 0x6c, 0x74, 0x61, 0x73,
	0x12, 0x3f, 0x0a, 0x0c, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x66, 0x2e, 0x73, 0x75, 0x62, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x44,
	0x65, 0x6c, 0x74, 0x61, 0x52, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x44, 0x65, 0x6c, 0x74, 0x61,
	0x73, 0x22, 0xf4, 0x01, 0x0a, 0x0a, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x44, 0x65, 0x6c, 0x74, 0x61,
	0x12, 0x44, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x73, 0x66, 0x2e, 0x73, 0x75, 0x62, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x44, 0x65, 0x6c, 0x74,
	0x61, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x6c,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x6f, 0x6c, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6f, 0x6c, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x6e, 0x65, 0x77, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x08, 0x6e, 0x65, 0x77, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3a, 0x0a, 0x09,
	0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x09, 0x0a, 0x05, 0x55, 0x4e, 0x53,
	0x45, 0x54, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x10, 0x01,
	0x12, 0x0a, 0x0a, 0x06, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06,
	0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x03, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67,
	0x66, 0x61, 0x73, 0x74, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x2f,
	0x70, 0x62, 0x2f, 0x73, 0x66, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73,
	0x2f, 0x76, 0x31, 0x3b, 0x70, 0x62, 0x73, 0x75, 0x62, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sf_substreams_v1_deltas_proto_rawDescOnce sync.Once
	file_sf_substreams_v1_deltas_proto_rawDescData = file_sf_substreams_v1_deltas_proto_rawDesc
)

func file_sf_substreams_v1_deltas_proto_rawDescGZIP() []byte {
	file_sf_substreams_v1_deltas_proto_rawDescOnce.Do(func() {
		file_sf_substreams_v1_deltas_proto_rawDescData = protoimpl.X.CompressGZIP(file_sf_substreams_v1_deltas_proto_rawDescData)
	})
	return file_sf_substreams_v1_deltas_proto_rawDescData
}

var file_sf_substreams_v1_deltas_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_sf_substreams_v1_deltas_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_sf_substreams_v1_deltas_proto_goTypes = []any{
	(StoreDelta_Operation)(0), // 0: sf.substreams.v1.StoreDelta.Operation
	(*StoreDeltas)(nil),       // 1: sf.substreams.v1.StoreDeltas
	(*StoreDelta)(nil),        // 2: sf.substreams.v1.StoreDelta
}
var file_sf_substreams_v1_deltas_proto_depIdxs = []int32{
	2, // 0: sf.substreams.v1.StoreDeltas.store_deltas:type_name -> sf.substreams.v1.StoreDelta
	0, // 1: sf.substreams.v1.StoreDelta.operation:type_name -> sf.substreams.v1.StoreDelta.Operation
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_sf_substreams_v1_deltas_proto_init() }
func file_sf_substreams_v1_deltas_proto_init() {
	if File_sf_substreams_v1_deltas_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sf_substreams_v1_deltas_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*StoreDeltas); i {
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
		file_sf_substreams_v1_deltas_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*StoreDelta); i {
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
			RawDescriptor: file_sf_substreams_v1_deltas_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sf_substreams_v1_deltas_proto_goTypes,
		DependencyIndexes: file_sf_substreams_v1_deltas_proto_depIdxs,
		EnumInfos:         file_sf_substreams_v1_deltas_proto_enumTypes,
		MessageInfos:      file_sf_substreams_v1_deltas_proto_msgTypes,
	}.Build()
	File_sf_substreams_v1_deltas_proto = out.File
	file_sf_substreams_v1_deltas_proto_rawDesc = nil
	file_sf_substreams_v1_deltas_proto_goTypes = nil
	file_sf_substreams_v1_deltas_proto_depIdxs = nil
}
