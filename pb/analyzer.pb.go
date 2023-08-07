// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: pb/analyzer.proto

package pb

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

type ParseMode int32

const (
	ParseMode_Parsed   ParseMode = 0
	ParseMode_Unparsed ParseMode = 1
)

// Enum value maps for ParseMode.
var (
	ParseMode_name = map[int32]string{
		0: "Parsed",
		1: "Unparsed",
	}
	ParseMode_value = map[string]int32{
		"Parsed":   0,
		"Unparsed": 1,
	}
)

func (x ParseMode) Enum() *ParseMode {
	p := new(ParseMode)
	*p = x
	return p
}

func (x ParseMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ParseMode) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_analyzer_proto_enumTypes[0].Descriptor()
}

func (ParseMode) Type() protoreflect.EnumType {
	return &file_pb_analyzer_proto_enumTypes[0]
}

func (x ParseMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ParseMode.Descriptor instead.
func (ParseMode) EnumDescriptor() ([]byte, []int) {
	return file_pb_analyzer_proto_rawDescGZIP(), []int{0}
}

type UUID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *UUID) Reset() {
	*x = UUID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_analyzer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UUID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UUID) ProtoMessage() {}

func (x *UUID) ProtoReflect() protoreflect.Message {
	mi := &file_pb_analyzer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UUID.ProtoReflect.Descriptor instead.
func (*UUID) Descriptor() ([]byte, []int) {
	return file_pb_analyzer_proto_rawDescGZIP(), []int{0}
}

func (x *UUID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type LogInstance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       *UUID  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ServName string `protobuf:"bytes,2,opt,name=servName,proto3" json:"servName,omitempty"`
	Log      string `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
}

func (x *LogInstance) Reset() {
	*x = LogInstance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_analyzer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogInstance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogInstance) ProtoMessage() {}

func (x *LogInstance) ProtoReflect() protoreflect.Message {
	mi := &file_pb_analyzer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogInstance.ProtoReflect.Descriptor instead.
func (*LogInstance) Descriptor() ([]byte, []int) {
	return file_pb_analyzer_proto_rawDescGZIP(), []int{1}
}

func (x *LogInstance) GetId() *UUID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *LogInstance) GetServName() string {
	if x != nil {
		return x.ServName
	}
	return ""
}

func (x *LogInstance) GetLog() string {
	if x != nil {
		return x.Log
	}
	return ""
}

type AnalyzerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        *UUID          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ParseMode ParseMode      `protobuf:"varint,2,opt,name=parseMode,proto3,enum=ParseMode" json:"parseMode,omitempty"`
	Recent    []*LogInstance `protobuf:"bytes,3,rep,name=recent,proto3" json:"recent,omitempty"`
	History   []*LogInstance `protobuf:"bytes,4,rep,name=history,proto3" json:"history,omitempty"`
}

func (x *AnalyzerRequest) Reset() {
	*x = AnalyzerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_analyzer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzerRequest) ProtoMessage() {}

func (x *AnalyzerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_analyzer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzerRequest.ProtoReflect.Descriptor instead.
func (*AnalyzerRequest) Descriptor() ([]byte, []int) {
	return file_pb_analyzer_proto_rawDescGZIP(), []int{2}
}

func (x *AnalyzerRequest) GetId() *UUID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *AnalyzerRequest) GetParseMode() ParseMode {
	if x != nil {
		return x.ParseMode
	}
	return ParseMode_Parsed
}

func (x *AnalyzerRequest) GetRecent() []*LogInstance {
	if x != nil {
		return x.Recent
	}
	return nil
}

func (x *AnalyzerRequest) GetHistory() []*LogInstance {
	if x != nil {
		return x.History
	}
	return nil
}

type AnalyzerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       *UUID  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Rating   uint32 `protobuf:"varint,2,opt,name=rating,proto3" json:"rating,omitempty"`
	Review   string `protobuf:"bytes,3,opt,name=review,proto3" json:"review,omitempty"`
	Insight  string `protobuf:"bytes,4,opt,name=insight,proto3" json:"insight,omitempty"`
	Citation string `protobuf:"bytes,5,opt,name=citation,proto3" json:"citation,omitempty"`
}

func (x *AnalyzerResponse) Reset() {
	*x = AnalyzerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_analyzer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzerResponse) ProtoMessage() {}

func (x *AnalyzerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_analyzer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzerResponse.ProtoReflect.Descriptor instead.
func (*AnalyzerResponse) Descriptor() ([]byte, []int) {
	return file_pb_analyzer_proto_rawDescGZIP(), []int{3}
}

func (x *AnalyzerResponse) GetId() *UUID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *AnalyzerResponse) GetRating() uint32 {
	if x != nil {
		return x.Rating
	}
	return 0
}

func (x *AnalyzerResponse) GetReview() string {
	if x != nil {
		return x.Review
	}
	return ""
}

func (x *AnalyzerResponse) GetInsight() string {
	if x != nil {
		return x.Insight
	}
	return ""
}

func (x *AnalyzerResponse) GetCitation() string {
	if x != nil {
		return x.Citation
	}
	return ""
}

var File_pb_analyzer_proto protoreflect.FileDescriptor

var file_pb_analyzer_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x62, 0x2f, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x16, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x52, 0x0a, 0x0b, 0x4c,
	0x6f, 0x67, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x15, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6c, 0x6f, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6c, 0x6f, 0x67, 0x22,
	0xa0, 0x01, 0x0a, 0x0f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x28, 0x0a, 0x09, 0x70, 0x61,
	0x72, 0x73, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e,
	0x50, 0x61, 0x72, 0x73, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x09, 0x70, 0x61, 0x72, 0x73, 0x65,
	0x4d, 0x6f, 0x64, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x52, 0x06, 0x72, 0x65, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x26, 0x0a, 0x07, 0x68, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4c, 0x6f,
	0x67, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x07, 0x68, 0x69, 0x73, 0x74, 0x6f,
	0x72, 0x79, 0x22, 0x8f, 0x01, 0x0a, 0x10, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x15, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x12, 0x18,
	0x0a, 0x07, 0x69, 0x6e, 0x73, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x69, 0x6e, 0x73, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x69, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x69, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x25, 0x0a, 0x09, 0x50, 0x61, 0x72, 0x73, 0x65, 0x4d, 0x6f, 0x64,
	0x65, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x73, 0x65, 0x64, 0x10, 0x00, 0x12, 0x0c, 0x0a,
	0x08, 0x55, 0x6e, 0x70, 0x61, 0x72, 0x73, 0x65, 0x64, 0x10, 0x01, 0x32, 0x43, 0x0a, 0x08, 0x41,
	0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x0a, 0x61, 0x6e, 0x61, 0x6c, 0x79,
	0x7a, 0x65, 0x4c, 0x6f, 0x67, 0x12, 0x10, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01,
	0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_analyzer_proto_rawDescOnce sync.Once
	file_pb_analyzer_proto_rawDescData = file_pb_analyzer_proto_rawDesc
)

func file_pb_analyzer_proto_rawDescGZIP() []byte {
	file_pb_analyzer_proto_rawDescOnce.Do(func() {
		file_pb_analyzer_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_analyzer_proto_rawDescData)
	})
	return file_pb_analyzer_proto_rawDescData
}

var file_pb_analyzer_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pb_analyzer_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pb_analyzer_proto_goTypes = []interface{}{
	(ParseMode)(0),           // 0: ParseMode
	(*UUID)(nil),             // 1: UUID
	(*LogInstance)(nil),      // 2: LogInstance
	(*AnalyzerRequest)(nil),  // 3: AnalyzerRequest
	(*AnalyzerResponse)(nil), // 4: AnalyzerResponse
}
var file_pb_analyzer_proto_depIdxs = []int32{
	1, // 0: LogInstance.id:type_name -> UUID
	1, // 1: AnalyzerRequest.id:type_name -> UUID
	0, // 2: AnalyzerRequest.parseMode:type_name -> ParseMode
	2, // 3: AnalyzerRequest.recent:type_name -> LogInstance
	2, // 4: AnalyzerRequest.history:type_name -> LogInstance
	1, // 5: AnalyzerResponse.id:type_name -> UUID
	3, // 6: Analyzer.analyzeLog:input_type -> AnalyzerRequest
	4, // 7: Analyzer.analyzeLog:output_type -> AnalyzerResponse
	7, // [7:8] is the sub-list for method output_type
	6, // [6:7] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_pb_analyzer_proto_init() }
func file_pb_analyzer_proto_init() {
	if File_pb_analyzer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_analyzer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UUID); i {
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
		file_pb_analyzer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogInstance); i {
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
		file_pb_analyzer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzerRequest); i {
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
		file_pb_analyzer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzerResponse); i {
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
			RawDescriptor: file_pb_analyzer_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_analyzer_proto_goTypes,
		DependencyIndexes: file_pb_analyzer_proto_depIdxs,
		EnumInfos:         file_pb_analyzer_proto_enumTypes,
		MessageInfos:      file_pb_analyzer_proto_msgTypes,
	}.Build()
	File_pb_analyzer_proto = out.File
	file_pb_analyzer_proto_rawDesc = nil
	file_pb_analyzer_proto_goTypes = nil
	file_pb_analyzer_proto_depIdxs = nil
}
