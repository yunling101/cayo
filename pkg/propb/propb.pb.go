// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.19.4
// source: propb.proto

package propb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string            `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Nid   int32             `protobuf:"varint,2,opt,name=nid,proto3" json:"nid,omitempty"`
	Param map[string]string `protobuf:"bytes,3,rep,name=Param,proto3" json:"Param,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_propb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_propb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_propb_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Request) GetNid() int32 {
	if x != nil {
		return x.Nid
	}
	return 0
}

func (x *Request) GetParam() map[string]string {
	if x != nil {
		return x.Param
	}
	return nil
}

type Reply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code   int32             `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg    string            `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Result map[string]string `protobuf:"bytes,3,rep,name=Result,proto3" json:"Result,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Reply) Reset() {
	*x = Reply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_propb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reply) ProtoMessage() {}

func (x *Reply) ProtoReflect() protoreflect.Message {
	mi := &file_propb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reply.ProtoReflect.Descriptor instead.
func (*Reply) Descriptor() ([]byte, []int) {
	return file_propb_proto_rawDescGZIP(), []int{1}
}

func (x *Reply) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Reply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Reply) GetResult() map[string]string {
	if x != nil {
		return x.Result
	}
	return nil
}

var File_propb_proto protoreflect.FileDescriptor

var file_propb_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x70, 0x62, 0x22, 0x98, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x6e, 0x69, 0x64, 0x12, 0x2f, 0x0a, 0x05, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x1a, 0x38, 0x0a, 0x0a, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x9a, 0x01, 0x0a, 0x05, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12,
	0x30, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x70, 0x72, 0x6f, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x1a, 0x39, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x93, 0x01, 0x0a,
	0x07, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x12, 0x2c, 0x0a, 0x0a, 0x4f, 0x62, 0x74, 0x61,
	0x69, 0x6e, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x70, 0x62, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x70, 0x62, 0x2e, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x0b, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70,
	0x74, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x70, 0x62, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x70, 0x62, 0x2e, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x2b, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65,
	0x61, 0x74, 0x12, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x00, 0x42, 0x03, 0x5a, 0x01, 0x2e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_propb_proto_rawDescOnce sync.Once
	file_propb_proto_rawDescData = file_propb_proto_rawDesc
)

func file_propb_proto_rawDescGZIP() []byte {
	file_propb_proto_rawDescOnce.Do(func() {
		file_propb_proto_rawDescData = protoimpl.X.CompressGZIP(file_propb_proto_rawDescData)
	})
	return file_propb_proto_rawDescData
}

var file_propb_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_propb_proto_goTypes = []interface{}{
	(*Request)(nil), // 0: propb.Request
	(*Reply)(nil),   // 1: propb.Reply
	nil,             // 2: propb.Request.ParamEntry
	nil,             // 3: propb.Reply.ResultEntry
}
var file_propb_proto_depIdxs = []int32{
	2, // 0: propb.Request.Param:type_name -> propb.Request.ParamEntry
	3, // 1: propb.Reply.Result:type_name -> propb.Reply.ResultEntry
	0, // 2: propb.Monitor.ObtainTask:input_type -> propb.Request
	0, // 3: propb.Monitor.ReceiptTask:input_type -> propb.Request
	0, // 4: propb.Monitor.Heartbeat:input_type -> propb.Request
	1, // 5: propb.Monitor.ObtainTask:output_type -> propb.Reply
	1, // 6: propb.Monitor.ReceiptTask:output_type -> propb.Reply
	1, // 7: propb.Monitor.Heartbeat:output_type -> propb.Reply
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_propb_proto_init() }
func file_propb_proto_init() {
	if File_propb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_propb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_propb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reply); i {
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
			RawDescriptor: file_propb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_propb_proto_goTypes,
		DependencyIndexes: file_propb_proto_depIdxs,
		MessageInfos:      file_propb_proto_msgTypes,
	}.Build()
	File_propb_proto = out.File
	file_propb_proto_rawDesc = nil
	file_propb_proto_goTypes = nil
	file_propb_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MonitorClient is the client API for Monitor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MonitorClient interface {
	ObtainTask(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error)
	ReceiptTask(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error)
	Heartbeat(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error)
}

type monitorClient struct {
	cc grpc.ClientConnInterface
}

func NewMonitorClient(cc grpc.ClientConnInterface) MonitorClient {
	return &monitorClient{cc}
}

func (c *monitorClient) ObtainTask(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, "/propb.Monitor/ObtainTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitorClient) ReceiptTask(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, "/propb.Monitor/ReceiptTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitorClient) Heartbeat(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, "/propb.Monitor/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MonitorServer is the server API for Monitor service.
type MonitorServer interface {
	ObtainTask(context.Context, *Request) (*Reply, error)
	ReceiptTask(context.Context, *Request) (*Reply, error)
	Heartbeat(context.Context, *Request) (*Reply, error)
}

// UnimplementedMonitorServer can be embedded to have forward compatible implementations.
type UnimplementedMonitorServer struct {
}

func (*UnimplementedMonitorServer) ObtainTask(context.Context, *Request) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ObtainTask not implemented")
}
func (*UnimplementedMonitorServer) ReceiptTask(context.Context, *Request) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiptTask not implemented")
}
func (*UnimplementedMonitorServer) Heartbeat(context.Context, *Request) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}

func RegisterMonitorServer(s *grpc.Server, srv MonitorServer) {
	s.RegisterService(&_Monitor_serviceDesc, srv)
}

func _Monitor_ObtainTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitorServer).ObtainTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/propb.Monitor/ObtainTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitorServer).ObtainTask(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Monitor_ReceiptTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitorServer).ReceiptTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/propb.Monitor/ReceiptTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitorServer).ReceiptTask(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Monitor_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitorServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/propb.Monitor/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitorServer).Heartbeat(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Monitor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "propb.Monitor",
	HandlerType: (*MonitorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ObtainTask",
			Handler:    _Monitor_ObtainTask_Handler,
		},
		{
			MethodName: "ReceiptTask",
			Handler:    _Monitor_ReceiptTask_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _Monitor_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "propb.proto",
}