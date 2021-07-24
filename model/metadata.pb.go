// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: model/metadata.proto

package model

import (
	any "github.com/golang/protobuf/ptypes/any"
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

type Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UniqueId       string   `protobuf:"bytes,1,opt,name=unique_id,json=uniqueId,proto3" json:"unique_id,omitempty"`
	Code           uint64   `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	ContractIdList []string `protobuf:"bytes,3,rep,name=contract_id_list,json=contractIdList,proto3" json:"contract_id_list,omitempty"`
	// useful for http trigger
	Path   string            `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	Method string            `protobuf:"bytes,5,opt,name=method,proto3" json:"method,omitempty"`
	Query  map[string]string `protobuf:"bytes,6,rep,name=query,proto3" json:"query,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// attach it to response
	// important for http trigger
	Headers map[string]string `protobuf:"bytes,7,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// attach it to response
	// important for http trigger
	StatusCode uint32   `protobuf:"varint,8,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	Status     string   `protobuf:"bytes,9,opt,name=status,proto3" json:"status,omitempty"`
	Data       *any.Any `protobuf:"bytes,500,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_metadata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_model_metadata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_model_metadata_proto_rawDescGZIP(), []int{0}
}

func (x *Metadata) GetUniqueId() string {
	if x != nil {
		return x.UniqueId
	}
	return ""
}

func (x *Metadata) GetCode() uint64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Metadata) GetContractIdList() []string {
	if x != nil {
		return x.ContractIdList
	}
	return nil
}

func (x *Metadata) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Metadata) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Metadata) GetQuery() map[string]string {
	if x != nil {
		return x.Query
	}
	return nil
}

func (x *Metadata) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *Metadata) GetStatusCode() uint32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *Metadata) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Metadata) GetData() *any.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_model_metadata_proto protoreflect.FileDescriptor

var file_model_metadata_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x19, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61,
	0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd5, 0x03, 0x0a, 0x08, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65,
	0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x64, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x70, 0x61, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x30, 0x0a, 0x05,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x12, 0x36,
	0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x29, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0xf4, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x38, 0x0a, 0x0a, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d,
	0x6b, 0x61, 0x77, 0x73, 0x65, 0x72, 0x6d, 0x2f, 0x61, 0x62, 0x65, 0x73, 0x68, 0x2f, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_model_metadata_proto_rawDescOnce sync.Once
	file_model_metadata_proto_rawDescData = file_model_metadata_proto_rawDesc
)

func file_model_metadata_proto_rawDescGZIP() []byte {
	file_model_metadata_proto_rawDescOnce.Do(func() {
		file_model_metadata_proto_rawDescData = protoimpl.X.CompressGZIP(file_model_metadata_proto_rawDescData)
	})
	return file_model_metadata_proto_rawDescData
}

var file_model_metadata_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_model_metadata_proto_goTypes = []interface{}{
	(*Metadata)(nil), // 0: model.Metadata
	nil,              // 1: model.Metadata.QueryEntry
	nil,              // 2: model.Metadata.HeadersEntry
	(*any.Any)(nil),  // 3: google.protobuf.Any
}
var file_model_metadata_proto_depIdxs = []int32{
	1, // 0: model.Metadata.query:type_name -> model.Metadata.QueryEntry
	2, // 1: model.Metadata.headers:type_name -> model.Metadata.HeadersEntry
	3, // 2: model.Metadata.data:type_name -> google.protobuf.Any
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_model_metadata_proto_init() }
func file_model_metadata_proto_init() {
	if File_model_metadata_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_model_metadata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metadata); i {
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
			RawDescriptor: file_model_metadata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_model_metadata_proto_goTypes,
		DependencyIndexes: file_model_metadata_proto_depIdxs,
		MessageInfos:      file_model_metadata_proto_msgTypes,
	}.Build()
	File_model_metadata_proto = out.File
	file_model_metadata_proto_rawDesc = nil
	file_model_metadata_proto_goTypes = nil
	file_model_metadata_proto_depIdxs = nil
}
