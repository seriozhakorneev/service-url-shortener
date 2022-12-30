// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: shortener_proto/shortener.proto

package shortener_proto

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

// The shortener message containing URL.
type ShortenerURLData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	URL string `protobuf:"bytes,1,opt,name=URL,proto3" json:"URL,omitempty"`
}

func (x *ShortenerURLData) Reset() {
	*x = ShortenerURLData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shortener_proto_shortener_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortenerURLData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenerURLData) ProtoMessage() {}

func (x *ShortenerURLData) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_shortener_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenerURLData.ProtoReflect.Descriptor instead.
func (*ShortenerURLData) Descriptor() ([]byte, []int) {
	return file_shortener_proto_shortener_proto_rawDescGZIP(), []int{0}
}

func (x *ShortenerURLData) GetURL() string {
	if x != nil {
		return x.URL
	}
	return ""
}

// The shortener message containing URL creation data.
type ShortenerCreateURLData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	URL string `protobuf:"bytes,1,opt,name=URL,proto3" json:"URL,omitempty"`
	// time to live in hours
	TTL int32 `protobuf:"varint,2,opt,name=TTL,proto3" json:"TTL,omitempty"`
}

func (x *ShortenerCreateURLData) Reset() {
	*x = ShortenerCreateURLData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shortener_proto_shortener_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortenerCreateURLData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenerCreateURLData) ProtoMessage() {}

func (x *ShortenerCreateURLData) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_shortener_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenerCreateURLData.ProtoReflect.Descriptor instead.
func (*ShortenerCreateURLData) Descriptor() ([]byte, []int) {
	return file_shortener_proto_shortener_proto_rawDescGZIP(), []int{1}
}

func (x *ShortenerCreateURLData) GetURL() string {
	if x != nil {
		return x.URL
	}
	return ""
}

func (x *ShortenerCreateURLData) GetTTL() int32 {
	if x != nil {
		return x.TTL
	}
	return 0
}

var File_shortener_proto_shortener_proto protoreflect.FileDescriptor

var file_shortener_proto_shortener_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x24, 0x0a, 0x10, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x55,
	0x52, 0x4c, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x52, 0x4c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x52, 0x4c, 0x22, 0x3c, 0x0a, 0x16, 0x53, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x52, 0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x55, 0x52, 0x4c, 0x12, 0x10, 0x0a, 0x03, 0x54, 0x54, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x03, 0x54, 0x54, 0x4c, 0x32, 0xb2, 0x01, 0x0a, 0x09, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x12, 0x56, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x27,
	0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x55, 0x52, 0x4c, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x21, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x55, 0x52, 0x4c, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x03,
	0x47, 0x65, 0x74, 0x12, 0x21, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x5f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x55,
	0x52, 0x4c, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x21, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x65, 0x72, 0x55, 0x52, 0x4c, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x42, 0x40, 0x5a, 0x3e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2d, 0x75, 0x72, 0x6c, 0x2d, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x6e,
	0x74, 0x72, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shortener_proto_shortener_proto_rawDescOnce sync.Once
	file_shortener_proto_shortener_proto_rawDescData = file_shortener_proto_shortener_proto_rawDesc
)

func file_shortener_proto_shortener_proto_rawDescGZIP() []byte {
	file_shortener_proto_shortener_proto_rawDescOnce.Do(func() {
		file_shortener_proto_shortener_proto_rawDescData = protoimpl.X.CompressGZIP(file_shortener_proto_shortener_proto_rawDescData)
	})
	return file_shortener_proto_shortener_proto_rawDescData
}

var file_shortener_proto_shortener_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_shortener_proto_shortener_proto_goTypes = []interface{}{
	(*ShortenerURLData)(nil),       // 0: shortener_proto.ShortenerURLData
	(*ShortenerCreateURLData)(nil), // 1: shortener_proto.ShortenerCreateURLData
}
var file_shortener_proto_shortener_proto_depIdxs = []int32{
	1, // 0: shortener_proto.Shortener.Create:input_type -> shortener_proto.ShortenerCreateURLData
	0, // 1: shortener_proto.Shortener.Get:input_type -> shortener_proto.ShortenerURLData
	0, // 2: shortener_proto.Shortener.Create:output_type -> shortener_proto.ShortenerURLData
	0, // 3: shortener_proto.Shortener.Get:output_type -> shortener_proto.ShortenerURLData
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_shortener_proto_shortener_proto_init() }
func file_shortener_proto_shortener_proto_init() {
	if File_shortener_proto_shortener_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shortener_proto_shortener_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortenerURLData); i {
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
		file_shortener_proto_shortener_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortenerCreateURLData); i {
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
			RawDescriptor: file_shortener_proto_shortener_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shortener_proto_shortener_proto_goTypes,
		DependencyIndexes: file_shortener_proto_shortener_proto_depIdxs,
		MessageInfos:      file_shortener_proto_shortener_proto_msgTypes,
	}.Build()
	File_shortener_proto_shortener_proto = out.File
	file_shortener_proto_shortener_proto_rawDesc = nil
	file_shortener_proto_shortener_proto_goTypes = nil
	file_shortener_proto_shortener_proto_depIdxs = nil
}
