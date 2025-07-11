// proto/url_service.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: proto/url_service.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// These replace your JSON structs
type CreateURLRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OriginalUrl   string                 `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateURLRequest) Reset() {
	*x = CreateURLRequest{}
	mi := &file_proto_url_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateURLRequest) ProtoMessage() {}

func (x *CreateURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_url_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateURLRequest.ProtoReflect.Descriptor instead.
func (*CreateURLRequest) Descriptor() ([]byte, []int) {
	return file_proto_url_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateURLRequest) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

func (x *CreateURLRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type CreateURLResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortCode     string                 `protobuf:"bytes,1,opt,name=short_code,json=shortCode,proto3" json:"short_code,omitempty"`
	ShortUrl      string                 `protobuf:"bytes,2,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	Success       bool                   `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	Error         string                 `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateURLResponse) Reset() {
	*x = CreateURLResponse{}
	mi := &file_proto_url_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateURLResponse) ProtoMessage() {}

func (x *CreateURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_url_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateURLResponse.ProtoReflect.Descriptor instead.
func (*CreateURLResponse) Descriptor() ([]byte, []int) {
	return file_proto_url_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateURLResponse) GetShortCode() string {
	if x != nil {
		return x.ShortCode
	}
	return ""
}

func (x *CreateURLResponse) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

func (x *CreateURLResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *CreateURLResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetURLRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortCode     string                 `protobuf:"bytes,1,opt,name=short_code,json=shortCode,proto3" json:"short_code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetURLRequest) Reset() {
	*x = GetURLRequest{}
	mi := &file_proto_url_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLRequest) ProtoMessage() {}

func (x *GetURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_url_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLRequest.ProtoReflect.Descriptor instead.
func (*GetURLRequest) Descriptor() ([]byte, []int) {
	return file_proto_url_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetURLRequest) GetShortCode() string {
	if x != nil {
		return x.ShortCode
	}
	return ""
}

type GetURLResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OriginalUrl   string                 `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	Found         bool                   `protobuf:"varint,2,opt,name=found,proto3" json:"found,omitempty"`
	Error         string                 `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetURLResponse) Reset() {
	*x = GetURLResponse{}
	mi := &file_proto_url_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLResponse) ProtoMessage() {}

func (x *GetURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_url_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLResponse.ProtoReflect.Descriptor instead.
func (*GetURLResponse) Descriptor() ([]byte, []int) {
	return file_proto_url_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetURLResponse) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

func (x *GetURLResponse) GetFound() bool {
	if x != nil {
		return x.Found
	}
	return false
}

func (x *GetURLResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type HealthRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthRequest) Reset() {
	*x = HealthRequest{}
	mi := &file_proto_url_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthRequest) ProtoMessage() {}

func (x *HealthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_url_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthRequest.ProtoReflect.Descriptor instead.
func (*HealthRequest) Descriptor() ([]byte, []int) {
	return file_proto_url_service_proto_rawDescGZIP(), []int{4}
}

type HealthResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Healthy       bool                   `protobuf:"varint,1,opt,name=healthy,proto3" json:"healthy,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	mi := &file_proto_url_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_url_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthResponse.ProtoReflect.Descriptor instead.
func (*HealthResponse) Descriptor() ([]byte, []int) {
	return file_proto_url_service_proto_rawDescGZIP(), []int{5}
}

func (x *HealthResponse) GetHealthy() bool {
	if x != nil {
		return x.Healthy
	}
	return false
}

var File_proto_url_service_proto protoreflect.FileDescriptor

const file_proto_url_service_proto_rawDesc = "" +
	"\n" +
	"\x17proto/url_service.proto\x12\n" +
	"urlservice\"N\n" +
	"\x10CreateURLRequest\x12!\n" +
	"\foriginal_url\x18\x01 \x01(\tR\voriginalUrl\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\tR\x06userId\"\x7f\n" +
	"\x11CreateURLResponse\x12\x1d\n" +
	"\n" +
	"short_code\x18\x01 \x01(\tR\tshortCode\x12\x1b\n" +
	"\tshort_url\x18\x02 \x01(\tR\bshortUrl\x12\x18\n" +
	"\asuccess\x18\x03 \x01(\bR\asuccess\x12\x14\n" +
	"\x05error\x18\x04 \x01(\tR\x05error\".\n" +
	"\rGetURLRequest\x12\x1d\n" +
	"\n" +
	"short_code\x18\x01 \x01(\tR\tshortCode\"_\n" +
	"\x0eGetURLResponse\x12!\n" +
	"\foriginal_url\x18\x01 \x01(\tR\voriginalUrl\x12\x14\n" +
	"\x05found\x18\x02 \x01(\bR\x05found\x12\x14\n" +
	"\x05error\x18\x03 \x01(\tR\x05error\"\x0f\n" +
	"\rHealthRequest\"*\n" +
	"\x0eHealthResponse\x12\x18\n" +
	"\ahealthy\x18\x01 \x01(\bR\ahealthy2\xea\x01\n" +
	"\n" +
	"URLService\x12M\n" +
	"\x0eCreateShortURL\x12\x1c.urlservice.CreateURLRequest\x1a\x1d.urlservice.CreateURLResponse\x12G\n" +
	"\x0eGetOriginalURL\x12\x19.urlservice.GetURLRequest\x1a\x1a.urlservice.GetURLResponse\x12D\n" +
	"\vHealthCheck\x12\x19.urlservice.HealthRequest\x1a\x1a.urlservice.HealthResponseB+Z)github.com/sammyqtran/url-shortener/protob\x06proto3"

var (
	file_proto_url_service_proto_rawDescOnce sync.Once
	file_proto_url_service_proto_rawDescData []byte
)

func file_proto_url_service_proto_rawDescGZIP() []byte {
	file_proto_url_service_proto_rawDescOnce.Do(func() {
		file_proto_url_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_url_service_proto_rawDesc), len(file_proto_url_service_proto_rawDesc)))
	})
	return file_proto_url_service_proto_rawDescData
}

var file_proto_url_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_url_service_proto_goTypes = []any{
	(*CreateURLRequest)(nil),  // 0: urlservice.CreateURLRequest
	(*CreateURLResponse)(nil), // 1: urlservice.CreateURLResponse
	(*GetURLRequest)(nil),     // 2: urlservice.GetURLRequest
	(*GetURLResponse)(nil),    // 3: urlservice.GetURLResponse
	(*HealthRequest)(nil),     // 4: urlservice.HealthRequest
	(*HealthResponse)(nil),    // 5: urlservice.HealthResponse
}
var file_proto_url_service_proto_depIdxs = []int32{
	0, // 0: urlservice.URLService.CreateShortURL:input_type -> urlservice.CreateURLRequest
	2, // 1: urlservice.URLService.GetOriginalURL:input_type -> urlservice.GetURLRequest
	4, // 2: urlservice.URLService.HealthCheck:input_type -> urlservice.HealthRequest
	1, // 3: urlservice.URLService.CreateShortURL:output_type -> urlservice.CreateURLResponse
	3, // 4: urlservice.URLService.GetOriginalURL:output_type -> urlservice.GetURLResponse
	5, // 5: urlservice.URLService.HealthCheck:output_type -> urlservice.HealthResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_url_service_proto_init() }
func file_proto_url_service_proto_init() {
	if File_proto_url_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_url_service_proto_rawDesc), len(file_proto_url_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_url_service_proto_goTypes,
		DependencyIndexes: file_proto_url_service_proto_depIdxs,
		MessageInfos:      file_proto_url_service_proto_msgTypes,
	}.Build()
	File_proto_url_service_proto = out.File
	file_proto_url_service_proto_goTypes = nil
	file_proto_url_service_proto_depIdxs = nil
}
