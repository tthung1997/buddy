// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: framework/grpc/proto/choice.proto

package choice

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Choice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Value           string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Weight          int32                  `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
	Color           string                 `protobuf:"bytes,4,opt,name=color,proto3" json:"color,omitempty"`
	UpdatedDateTime *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=updatedDateTime,proto3" json:"updatedDateTime,omitempty"`
}

func (x *Choice) Reset() {
	*x = Choice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_framework_grpc_proto_choice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Choice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Choice) ProtoMessage() {}

func (x *Choice) ProtoReflect() protoreflect.Message {
	mi := &file_framework_grpc_proto_choice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Choice.ProtoReflect.Descriptor instead.
func (*Choice) Descriptor() ([]byte, []int) {
	return file_framework_grpc_proto_choice_proto_rawDescGZIP(), []int{0}
}

func (x *Choice) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Choice) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Choice) GetWeight() int32 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *Choice) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *Choice) GetUpdatedDateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedDateTime
	}
	return nil
}

type ChoiceList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Choices         []*Choice              `protobuf:"bytes,2,rep,name=choices,proto3" json:"choices,omitempty"`
	UpdatedDateTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updatedDateTime,proto3" json:"updatedDateTime,omitempty"`
}

func (x *ChoiceList) Reset() {
	*x = ChoiceList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_framework_grpc_proto_choice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChoiceList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChoiceList) ProtoMessage() {}

func (x *ChoiceList) ProtoReflect() protoreflect.Message {
	mi := &file_framework_grpc_proto_choice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChoiceList.ProtoReflect.Descriptor instead.
func (*ChoiceList) Descriptor() ([]byte, []int) {
	return file_framework_grpc_proto_choice_proto_rawDescGZIP(), []int{1}
}

func (x *ChoiceList) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ChoiceList) GetChoices() []*Choice {
	if x != nil {
		return x.Choices
	}
	return nil
}

func (x *ChoiceList) GetUpdatedDateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedDateTime
	}
	return nil
}

type GetByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetByIdRequest) Reset() {
	*x = GetByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_framework_grpc_proto_choice_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIdRequest) ProtoMessage() {}

func (x *GetByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_framework_grpc_proto_choice_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIdRequest.ProtoReflect.Descriptor instead.
func (*GetByIdRequest) Descriptor() ([]byte, []int) {
	return file_framework_grpc_proto_choice_proto_rawDescGZIP(), []int{2}
}

func (x *GetByIdRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type UpsertResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error   string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpsertResponse) Reset() {
	*x = UpsertResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_framework_grpc_proto_choice_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertResponse) ProtoMessage() {}

func (x *UpsertResponse) ProtoReflect() protoreflect.Message {
	mi := &file_framework_grpc_proto_choice_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertResponse.ProtoReflect.Descriptor instead.
func (*UpsertResponse) Descriptor() ([]byte, []int) {
	return file_framework_grpc_proto_choice_proto_rawDescGZIP(), []int{3}
}

func (x *UpsertResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UpsertResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_framework_grpc_proto_choice_proto protoreflect.FileDescriptor

var file_framework_grpc_proto_choice_proto_rawDesc = []byte{
	0x0a, 0x21, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa2, 0x01, 0x0a,
	0x06, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x77,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x44, 0x0a, 0x0f, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x22, 0x8c, 0x01, 0x0a, 0x0a, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x28, 0x0a, 0x07, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x6f, 0x69, 0x63,
	0x65, 0x52, 0x07, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x12, 0x44, 0x0a, 0x0f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x22, 0x20, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x40, 0x0a, 0x0e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x32, 0x94, 0x01, 0x0a, 0x0d, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x68, 0x6f,
	0x69, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12, 0x16, 0x2e, 0x63, 0x68,
	0x6f, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x6f,
	0x69, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x10, 0x55, 0x70, 0x73,
	0x65, 0x72, 0x74, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x2e,
	0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x73, 0x65, 0x72,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e,
	0x2f, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_framework_grpc_proto_choice_proto_rawDescOnce sync.Once
	file_framework_grpc_proto_choice_proto_rawDescData = file_framework_grpc_proto_choice_proto_rawDesc
)

func file_framework_grpc_proto_choice_proto_rawDescGZIP() []byte {
	file_framework_grpc_proto_choice_proto_rawDescOnce.Do(func() {
		file_framework_grpc_proto_choice_proto_rawDescData = protoimpl.X.CompressGZIP(file_framework_grpc_proto_choice_proto_rawDescData)
	})
	return file_framework_grpc_proto_choice_proto_rawDescData
}

var file_framework_grpc_proto_choice_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_framework_grpc_proto_choice_proto_goTypes = []interface{}{
	(*Choice)(nil),                // 0: choice.Choice
	(*ChoiceList)(nil),            // 1: choice.ChoiceList
	(*GetByIdRequest)(nil),        // 2: choice.GetByIdRequest
	(*UpsertResponse)(nil),        // 3: choice.UpsertResponse
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_framework_grpc_proto_choice_proto_depIdxs = []int32{
	4, // 0: choice.Choice.updatedDateTime:type_name -> google.protobuf.Timestamp
	0, // 1: choice.ChoiceList.choices:type_name -> choice.Choice
	4, // 2: choice.ChoiceList.updatedDateTime:type_name -> google.protobuf.Timestamp
	2, // 3: choice.ChoiceService.GetChoiceListById:input_type -> choice.GetByIdRequest
	1, // 4: choice.ChoiceService.UpsertChoiceList:input_type -> choice.ChoiceList
	1, // 5: choice.ChoiceService.GetChoiceListById:output_type -> choice.ChoiceList
	3, // 6: choice.ChoiceService.UpsertChoiceList:output_type -> choice.UpsertResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_framework_grpc_proto_choice_proto_init() }
func file_framework_grpc_proto_choice_proto_init() {
	if File_framework_grpc_proto_choice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_framework_grpc_proto_choice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Choice); i {
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
		file_framework_grpc_proto_choice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChoiceList); i {
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
		file_framework_grpc_proto_choice_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByIdRequest); i {
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
		file_framework_grpc_proto_choice_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertResponse); i {
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
			RawDescriptor: file_framework_grpc_proto_choice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_framework_grpc_proto_choice_proto_goTypes,
		DependencyIndexes: file_framework_grpc_proto_choice_proto_depIdxs,
		MessageInfos:      file_framework_grpc_proto_choice_proto_msgTypes,
	}.Build()
	File_framework_grpc_proto_choice_proto = out.File
	file_framework_grpc_proto_choice_proto_rawDesc = nil
	file_framework_grpc_proto_choice_proto_goTypes = nil
	file_framework_grpc_proto_choice_proto_depIdxs = nil
}
