// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pkg/logproto/pattern.proto

package logproto

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	push "github.com/grafana/loki/pkg/push"
	github_com_prometheus_common_model "github.com/prometheus/common/model"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type QueryPatternsRequest struct {
	Selector string                                  `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
	From     github_com_prometheus_common_model.Time `protobuf:"varint,2,opt,name=from,proto3,customtype=github.com/prometheus/common/model.Time" json:"from"`
	Through  github_com_prometheus_common_model.Time `protobuf:"varint,3,opt,name=through,proto3,customtype=github.com/prometheus/common/model.Time" json:"through"`
}

func (m *QueryPatternsRequest) Reset()      { *m = QueryPatternsRequest{} }
func (*QueryPatternsRequest) ProtoMessage() {}
func (*QueryPatternsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_aaf4192acc66a4ea, []int{0}
}
func (m *QueryPatternsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPatternsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPatternsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPatternsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPatternsRequest.Merge(m, src)
}
func (m *QueryPatternsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryPatternsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPatternsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPatternsRequest proto.InternalMessageInfo

func (m *QueryPatternsRequest) GetSelector() string {
	if m != nil {
		return m.Selector
	}
	return ""
}

type QueryPatternsResponse struct {
	Series []*PatternSeries `protobuf:"bytes,1,rep,name=series,proto3" json:"series,omitempty"`
}

func (m *QueryPatternsResponse) Reset()      { *m = QueryPatternsResponse{} }
func (*QueryPatternsResponse) ProtoMessage() {}
func (*QueryPatternsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_aaf4192acc66a4ea, []int{1}
}
func (m *QueryPatternsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPatternsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPatternsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPatternsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPatternsResponse.Merge(m, src)
}
func (m *QueryPatternsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryPatternsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPatternsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPatternsResponse proto.InternalMessageInfo

func (m *QueryPatternsResponse) GetSeries() []*PatternSeries {
	if m != nil {
		return m.Series
	}
	return nil
}

type PatternSeries struct {
	Pattern string           `protobuf:"bytes,1,opt,name=pattern,proto3" json:"pattern,omitempty"`
	Samples []*PatternSample `protobuf:"bytes,2,rep,name=samples,proto3" json:"samples,omitempty"`
}

func (m *PatternSeries) Reset()      { *m = PatternSeries{} }
func (*PatternSeries) ProtoMessage() {}
func (*PatternSeries) Descriptor() ([]byte, []int) {
	return fileDescriptor_aaf4192acc66a4ea, []int{2}
}
func (m *PatternSeries) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PatternSeries) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PatternSeries.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PatternSeries) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PatternSeries.Merge(m, src)
}
func (m *PatternSeries) XXX_Size() int {
	return m.Size()
}
func (m *PatternSeries) XXX_DiscardUnknown() {
	xxx_messageInfo_PatternSeries.DiscardUnknown(m)
}

var xxx_messageInfo_PatternSeries proto.InternalMessageInfo

func (m *PatternSeries) GetPattern() string {
	if m != nil {
		return m.Pattern
	}
	return ""
}

func (m *PatternSeries) GetSamples() []*PatternSample {
	if m != nil {
		return m.Samples
	}
	return nil
}

type PatternSample struct {
	Timestamp github_com_prometheus_common_model.Time `protobuf:"varint,1,opt,name=timestamp,proto3,customtype=github.com/prometheus/common/model.Time" json:"timestamp"`
	Value     int64                                   `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *PatternSample) Reset()      { *m = PatternSample{} }
func (*PatternSample) ProtoMessage() {}
func (*PatternSample) Descriptor() ([]byte, []int) {
	return fileDescriptor_aaf4192acc66a4ea, []int{3}
}
func (m *PatternSample) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PatternSample) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PatternSample.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PatternSample) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PatternSample.Merge(m, src)
}
func (m *PatternSample) XXX_Size() int {
	return m.Size()
}
func (m *PatternSample) XXX_DiscardUnknown() {
	xxx_messageInfo_PatternSample.DiscardUnknown(m)
}

var xxx_messageInfo_PatternSample proto.InternalMessageInfo

func (m *PatternSample) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*QueryPatternsRequest)(nil), "logproto.QueryPatternsRequest")
	proto.RegisterType((*QueryPatternsResponse)(nil), "logproto.QueryPatternsResponse")
	proto.RegisterType((*PatternSeries)(nil), "logproto.PatternSeries")
	proto.RegisterType((*PatternSample)(nil), "logproto.PatternSample")
}

func init() { proto.RegisterFile("pkg/logproto/pattern.proto", fileDescriptor_aaf4192acc66a4ea) }

var fileDescriptor_aaf4192acc66a4ea = []byte{
	// 443 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0xbf, 0x6e, 0xd4, 0x30,
	0x1c, 0xb6, 0x7b, 0x6d, 0xd3, 0x1a, 0xb1, 0x98, 0x2b, 0x44, 0x41, 0xf2, 0x9d, 0x22, 0x24, 0x6e,
	0x8a, 0xa1, 0x0c, 0xec, 0x65, 0x01, 0x09, 0xa4, 0x12, 0x98, 0x2a, 0x96, 0xf4, 0xf8, 0x35, 0x89,
	0x1a, 0xc7, 0xc6, 0x7f, 0x90, 0xd8, 0x78, 0x84, 0x3e, 0x06, 0xcf, 0xc1, 0xd4, 0xf1, 0xc6, 0x8a,
	0xa1, 0xe2, 0x72, 0x0b, 0x63, 0x1f, 0x01, 0xd5, 0x49, 0x7a, 0xd7, 0xaa, 0x5d, 0xba, 0x24, 0xfe,
	0xfe, 0xf8, 0x73, 0xf2, 0xfd, 0x12, 0x12, 0xa9, 0xe3, 0x9c, 0x57, 0x32, 0x57, 0x5a, 0x5a, 0xc9,
	0x55, 0x66, 0x2d, 0xe8, 0x3a, 0xf1, 0x88, 0x6e, 0xf5, 0x7c, 0xf4, 0xf4, 0x9a, 0xab, 0x5f, 0xb4,
	0xb6, 0xe8, 0xd1, 0xa5, 0xa8, 0x9c, 0x29, 0xfc, 0xa5, 0x23, 0x87, 0xb9, 0xcc, 0x65, 0x6b, 0xbf,
	0x5c, 0xb5, 0x6c, 0xfc, 0x1b, 0x93, 0xe1, 0x47, 0x07, 0xfa, 0xc7, 0x7e, 0x7b, 0x90, 0x49, 0xe1,
	0x9b, 0x03, 0x63, 0x69, 0x44, 0xb6, 0x0c, 0x54, 0x30, 0xb5, 0x52, 0x87, 0x78, 0x8c, 0x27, 0xdb,
	0xe9, 0x15, 0xa6, 0x6f, 0xc8, 0xfa, 0x91, 0x96, 0x22, 0x5c, 0x1b, 0xe3, 0xc9, 0x60, 0x8f, 0x9f,
	0x9e, 0x8f, 0xd0, 0x9f, 0xf3, 0xd1, 0xf3, 0xbc, 0xb4, 0x85, 0x3b, 0x4c, 0xa6, 0x52, 0x70, 0xa5,
	0xa5, 0x00, 0x5b, 0x80, 0x33, 0x7c, 0x2a, 0x85, 0x90, 0x35, 0x17, 0xf2, 0x2b, 0x54, 0xc9, 0xe7,
	0x52, 0x40, 0xea, 0x37, 0xd3, 0x77, 0x24, 0xb0, 0x85, 0x96, 0x2e, 0x2f, 0xc2, 0xc1, 0xfd, 0x72,
	0xfa, 0xfd, 0xf1, 0x5b, 0xb2, 0x73, 0xe3, 0x1d, 0x8c, 0x92, 0xb5, 0x01, 0xca, 0xc9, 0xa6, 0x01,
	0x5d, 0x82, 0x09, 0xf1, 0x78, 0x30, 0x79, 0xb0, 0xfb, 0x24, 0xb9, 0x6a, 0xaa, 0xf3, 0x7e, 0xf2,
	0x72, 0xda, 0xd9, 0xe2, 0x2f, 0xe4, 0xe1, 0x35, 0x81, 0x86, 0x24, 0xe8, 0x46, 0xd0, 0xb5, 0xd0,
	0x43, 0xfa, 0x92, 0x04, 0x26, 0x13, 0xaa, 0x02, 0x13, 0xae, 0xdd, 0x15, 0xee, 0xf5, 0xb4, 0xf7,
	0xc5, 0x76, 0x99, 0xee, 0x19, 0xfa, 0x81, 0x6c, 0xdb, 0x52, 0x80, 0xb1, 0x99, 0x50, 0x3e, 0xff,
	0x1e, 0x2d, 0x2c, 0x13, 0xe8, 0x90, 0x6c, 0x7c, 0xcf, 0x2a, 0x07, 0xed, 0x60, 0xd2, 0x16, 0xec,
	0x9e, 0x60, 0x12, 0x74, 0xc7, 0xd2, 0xd7, 0x64, 0x7d, 0xdf, 0x99, 0x82, 0xee, 0xac, 0x3c, 0xab,
	0x33, 0x45, 0x37, 0xf4, 0xe8, 0xf1, 0x4d, 0xba, 0xed, 0x31, 0x46, 0xf4, 0x3d, 0xd9, 0xf0, 0x15,
	0x53, 0xb6, 0xb4, 0xdc, 0xf6, 0xdd, 0x44, 0xa3, 0x3b, 0xf5, 0x3e, 0xeb, 0x05, 0xde, 0x3b, 0x98,
	0xcd, 0x19, 0x3a, 0x9b, 0x33, 0x74, 0x31, 0x67, 0xf8, 0x67, 0xc3, 0xf0, 0xaf, 0x86, 0xe1, 0xd3,
	0x86, 0xe1, 0x59, 0xc3, 0xf0, 0xdf, 0x86, 0xe1, 0x7f, 0x0d, 0x43, 0x17, 0x0d, 0xc3, 0x27, 0x0b,
	0x86, 0x66, 0x0b, 0x86, 0xce, 0x16, 0x0c, 0x1d, 0x3c, 0x5b, 0xa9, 0x24, 0xd7, 0xd9, 0x51, 0x56,
	0x67, 0xbc, 0x92, 0xc7, 0x25, 0x5f, 0xfd, 0x17, 0x0e, 0x37, 0xfd, 0xed, 0xd5, 0xff, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xa5, 0x4d, 0x43, 0x9d, 0x48, 0x03, 0x00, 0x00,
}

func (this *QueryPatternsRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*QueryPatternsRequest)
	if !ok {
		that2, ok := that.(QueryPatternsRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Selector != that1.Selector {
		return false
	}
	if !this.From.Equal(that1.From) {
		return false
	}
	if !this.Through.Equal(that1.Through) {
		return false
	}
	return true
}
func (this *QueryPatternsResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*QueryPatternsResponse)
	if !ok {
		that2, ok := that.(QueryPatternsResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Series) != len(that1.Series) {
		return false
	}
	for i := range this.Series {
		if !this.Series[i].Equal(that1.Series[i]) {
			return false
		}
	}
	return true
}
func (this *PatternSeries) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PatternSeries)
	if !ok {
		that2, ok := that.(PatternSeries)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Pattern != that1.Pattern {
		return false
	}
	if len(this.Samples) != len(that1.Samples) {
		return false
	}
	for i := range this.Samples {
		if !this.Samples[i].Equal(that1.Samples[i]) {
			return false
		}
	}
	return true
}
func (this *PatternSample) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PatternSample)
	if !ok {
		that2, ok := that.(PatternSample)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	if this.Value != that1.Value {
		return false
	}
	return true
}
func (this *QueryPatternsRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&logproto.QueryPatternsRequest{")
	s = append(s, "Selector: "+fmt.Sprintf("%#v", this.Selector)+",\n")
	s = append(s, "From: "+fmt.Sprintf("%#v", this.From)+",\n")
	s = append(s, "Through: "+fmt.Sprintf("%#v", this.Through)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *QueryPatternsResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&logproto.QueryPatternsResponse{")
	if this.Series != nil {
		s = append(s, "Series: "+fmt.Sprintf("%#v", this.Series)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *PatternSeries) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&logproto.PatternSeries{")
	s = append(s, "Pattern: "+fmt.Sprintf("%#v", this.Pattern)+",\n")
	if this.Samples != nil {
		s = append(s, "Samples: "+fmt.Sprintf("%#v", this.Samples)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *PatternSample) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&logproto.PatternSample{")
	s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	s = append(s, "Value: "+fmt.Sprintf("%#v", this.Value)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringPattern(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PatternClient is the client API for Pattern service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PatternClient interface {
	Push(ctx context.Context, in *push.PushRequest, opts ...grpc.CallOption) (*push.PushResponse, error)
	Query(ctx context.Context, in *QueryPatternsRequest, opts ...grpc.CallOption) (Pattern_QueryClient, error)
}

type patternClient struct {
	cc *grpc.ClientConn
}

func NewPatternClient(cc *grpc.ClientConn) PatternClient {
	return &patternClient{cc}
}

func (c *patternClient) Push(ctx context.Context, in *push.PushRequest, opts ...grpc.CallOption) (*push.PushResponse, error) {
	out := new(push.PushResponse)
	err := c.cc.Invoke(ctx, "/logproto.Pattern/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *patternClient) Query(ctx context.Context, in *QueryPatternsRequest, opts ...grpc.CallOption) (Pattern_QueryClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Pattern_serviceDesc.Streams[0], "/logproto.Pattern/Query", opts...)
	if err != nil {
		return nil, err
	}
	x := &patternQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Pattern_QueryClient interface {
	Recv() (*QueryPatternsResponse, error)
	grpc.ClientStream
}

type patternQueryClient struct {
	grpc.ClientStream
}

func (x *patternQueryClient) Recv() (*QueryPatternsResponse, error) {
	m := new(QueryPatternsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PatternServer is the server API for Pattern service.
type PatternServer interface {
	Push(context.Context, *push.PushRequest) (*push.PushResponse, error)
	Query(*QueryPatternsRequest, Pattern_QueryServer) error
}

// UnimplementedPatternServer can be embedded to have forward compatible implementations.
type UnimplementedPatternServer struct {
}

func (*UnimplementedPatternServer) Push(ctx context.Context, req *push.PushRequest) (*push.PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (*UnimplementedPatternServer) Query(req *QueryPatternsRequest, srv Pattern_QueryServer) error {
	return status.Errorf(codes.Unimplemented, "method Query not implemented")
}

func RegisterPatternServer(s *grpc.Server, srv PatternServer) {
	s.RegisterService(&_Pattern_serviceDesc, srv)
}

func _Pattern_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(push.PushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PatternServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logproto.Pattern/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PatternServer).Push(ctx, req.(*push.PushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pattern_Query_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryPatternsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PatternServer).Query(m, &patternQueryServer{stream})
}

type Pattern_QueryServer interface {
	Send(*QueryPatternsResponse) error
	grpc.ServerStream
}

type patternQueryServer struct {
	grpc.ServerStream
}

func (x *patternQueryServer) Send(m *QueryPatternsResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _Pattern_serviceDesc = grpc.ServiceDesc{
	ServiceName: "logproto.Pattern",
	HandlerType: (*PatternServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _Pattern_Push_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Query",
			Handler:       _Pattern_Query_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/logproto/pattern.proto",
}

func (m *QueryPatternsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPatternsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPatternsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Through != 0 {
		i = encodeVarintPattern(dAtA, i, uint64(m.Through))
		i--
		dAtA[i] = 0x18
	}
	if m.From != 0 {
		i = encodeVarintPattern(dAtA, i, uint64(m.From))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Selector) > 0 {
		i -= len(m.Selector)
		copy(dAtA[i:], m.Selector)
		i = encodeVarintPattern(dAtA, i, uint64(len(m.Selector)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryPatternsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPatternsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPatternsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Series) > 0 {
		for iNdEx := len(m.Series) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Series[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPattern(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *PatternSeries) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PatternSeries) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PatternSeries) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Samples) > 0 {
		for iNdEx := len(m.Samples) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Samples[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPattern(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Pattern) > 0 {
		i -= len(m.Pattern)
		copy(dAtA[i:], m.Pattern)
		i = encodeVarintPattern(dAtA, i, uint64(len(m.Pattern)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PatternSample) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PatternSample) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PatternSample) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != 0 {
		i = encodeVarintPattern(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x10
	}
	if m.Timestamp != 0 {
		i = encodeVarintPattern(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPattern(dAtA []byte, offset int, v uint64) int {
	offset -= sovPattern(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryPatternsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Selector)
	if l > 0 {
		n += 1 + l + sovPattern(uint64(l))
	}
	if m.From != 0 {
		n += 1 + sovPattern(uint64(m.From))
	}
	if m.Through != 0 {
		n += 1 + sovPattern(uint64(m.Through))
	}
	return n
}

func (m *QueryPatternsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Series) > 0 {
		for _, e := range m.Series {
			l = e.Size()
			n += 1 + l + sovPattern(uint64(l))
		}
	}
	return n
}

func (m *PatternSeries) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Pattern)
	if l > 0 {
		n += 1 + l + sovPattern(uint64(l))
	}
	if len(m.Samples) > 0 {
		for _, e := range m.Samples {
			l = e.Size()
			n += 1 + l + sovPattern(uint64(l))
		}
	}
	return n
}

func (m *PatternSample) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Timestamp != 0 {
		n += 1 + sovPattern(uint64(m.Timestamp))
	}
	if m.Value != 0 {
		n += 1 + sovPattern(uint64(m.Value))
	}
	return n
}

func sovPattern(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPattern(x uint64) (n int) {
	return sovPattern(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *QueryPatternsRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&QueryPatternsRequest{`,
		`Selector:` + fmt.Sprintf("%v", this.Selector) + `,`,
		`From:` + fmt.Sprintf("%v", this.From) + `,`,
		`Through:` + fmt.Sprintf("%v", this.Through) + `,`,
		`}`,
	}, "")
	return s
}
func (this *QueryPatternsResponse) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForSeries := "[]*PatternSeries{"
	for _, f := range this.Series {
		repeatedStringForSeries += strings.Replace(f.String(), "PatternSeries", "PatternSeries", 1) + ","
	}
	repeatedStringForSeries += "}"
	s := strings.Join([]string{`&QueryPatternsResponse{`,
		`Series:` + repeatedStringForSeries + `,`,
		`}`,
	}, "")
	return s
}
func (this *PatternSeries) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForSamples := "[]*PatternSample{"
	for _, f := range this.Samples {
		repeatedStringForSamples += strings.Replace(f.String(), "PatternSample", "PatternSample", 1) + ","
	}
	repeatedStringForSamples += "}"
	s := strings.Join([]string{`&PatternSeries{`,
		`Pattern:` + fmt.Sprintf("%v", this.Pattern) + `,`,
		`Samples:` + repeatedStringForSamples + `,`,
		`}`,
	}, "")
	return s
}
func (this *PatternSample) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PatternSample{`,
		`Timestamp:` + fmt.Sprintf("%v", this.Timestamp) + `,`,
		`Value:` + fmt.Sprintf("%v", this.Value) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringPattern(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *QueryPatternsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPattern
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryPatternsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPatternsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Selector", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPattern
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPattern
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Selector = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			m.From = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.From |= github_com_prometheus_common_model.Time(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Through", wireType)
			}
			m.Through = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Through |= github_com_prometheus_common_model.Time(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPattern(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPattern
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPattern
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryPatternsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPattern
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryPatternsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPatternsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Series", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPattern
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPattern
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Series = append(m.Series, &PatternSeries{})
			if err := m.Series[len(m.Series)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPattern(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPattern
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPattern
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PatternSeries) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPattern
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PatternSeries: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PatternSeries: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pattern", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPattern
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPattern
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pattern = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Samples", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPattern
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPattern
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Samples = append(m.Samples, &PatternSample{})
			if err := m.Samples[len(m.Samples)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPattern(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPattern
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPattern
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PatternSample) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPattern
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PatternSample: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PatternSample: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= github_com_prometheus_common_model.Time(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			m.Value = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Value |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPattern(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPattern
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPattern
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPattern(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPattern
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPattern
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPattern
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthPattern
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowPattern
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipPattern(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthPattern
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthPattern = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPattern   = fmt.Errorf("proto: integer overflow")
)
