// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: data/v1/data.proto

package datav1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DataServiceClient is the client API for DataService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataServiceClient interface {
	GetTest(ctx context.Context, in *GetTestRequest, opts ...grpc.CallOption) (*GetTestResponse, error)
	QueryTests(ctx context.Context, in *QueryTestsRequest, opts ...grpc.CallOption) (*QueryTestsResponse, error)
	GetRun(ctx context.Context, in *GetRunRequest, opts ...grpc.CallOption) (*GetRunResponse, error)
	QueryRuns(ctx context.Context, in *QueryRunsRequest, opts ...grpc.CallOption) (*QueryRunsResponse, error)
	SummarizeRuns(ctx context.Context, in *SummarizeRunsRequest, opts ...grpc.CallOption) (*SummarizeRunsResponse, error)
	GetRunner(ctx context.Context, in *GetRunnerRequest, opts ...grpc.CallOption) (*GetRunnerResponse, error)
	QueryRunners(ctx context.Context, in *QueryRunnersRequest, opts ...grpc.CallOption) (*QueryRunnersResponse, error)
}

type dataServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataServiceClient(cc grpc.ClientConnInterface) DataServiceClient {
	return &dataServiceClient{cc}
}

func (c *dataServiceClient) GetTest(ctx context.Context, in *GetTestRequest, opts ...grpc.CallOption) (*GetTestResponse, error) {
	out := new(GetTestResponse)
	err := c.cc.Invoke(ctx, "/tstr.data.v1.DataService/GetTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataServiceClient) QueryTests(ctx context.Context, in *QueryTestsRequest, opts ...grpc.CallOption) (*QueryTestsResponse, error) {
	out := new(QueryTestsResponse)
	err := c.cc.Invoke(ctx, "/tstr.data.v1.DataService/QueryTests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataServiceClient) GetRun(ctx context.Context, in *GetRunRequest, opts ...grpc.CallOption) (*GetRunResponse, error) {
	out := new(GetRunResponse)
	err := c.cc.Invoke(ctx, "/tstr.data.v1.DataService/GetRun", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataServiceClient) QueryRuns(ctx context.Context, in *QueryRunsRequest, opts ...grpc.CallOption) (*QueryRunsResponse, error) {
	out := new(QueryRunsResponse)
	err := c.cc.Invoke(ctx, "/tstr.data.v1.DataService/QueryRuns", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataServiceClient) SummarizeRuns(ctx context.Context, in *SummarizeRunsRequest, opts ...grpc.CallOption) (*SummarizeRunsResponse, error) {
	out := new(SummarizeRunsResponse)
	err := c.cc.Invoke(ctx, "/tstr.data.v1.DataService/SummarizeRuns", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataServiceClient) GetRunner(ctx context.Context, in *GetRunnerRequest, opts ...grpc.CallOption) (*GetRunnerResponse, error) {
	out := new(GetRunnerResponse)
	err := c.cc.Invoke(ctx, "/tstr.data.v1.DataService/GetRunner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataServiceClient) QueryRunners(ctx context.Context, in *QueryRunnersRequest, opts ...grpc.CallOption) (*QueryRunnersResponse, error) {
	out := new(QueryRunnersResponse)
	err := c.cc.Invoke(ctx, "/tstr.data.v1.DataService/QueryRunners", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataServiceServer is the server API for DataService service.
// All implementations must embed UnimplementedDataServiceServer
// for forward compatibility
type DataServiceServer interface {
	GetTest(context.Context, *GetTestRequest) (*GetTestResponse, error)
	QueryTests(context.Context, *QueryTestsRequest) (*QueryTestsResponse, error)
	GetRun(context.Context, *GetRunRequest) (*GetRunResponse, error)
	QueryRuns(context.Context, *QueryRunsRequest) (*QueryRunsResponse, error)
	SummarizeRuns(context.Context, *SummarizeRunsRequest) (*SummarizeRunsResponse, error)
	GetRunner(context.Context, *GetRunnerRequest) (*GetRunnerResponse, error)
	QueryRunners(context.Context, *QueryRunnersRequest) (*QueryRunnersResponse, error)
	mustEmbedUnimplementedDataServiceServer()
}

// UnimplementedDataServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDataServiceServer struct {
}

func (UnimplementedDataServiceServer) GetTest(context.Context, *GetTestRequest) (*GetTestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTest not implemented")
}
func (UnimplementedDataServiceServer) QueryTests(context.Context, *QueryTestsRequest) (*QueryTestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryTests not implemented")
}
func (UnimplementedDataServiceServer) GetRun(context.Context, *GetRunRequest) (*GetRunResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRun not implemented")
}
func (UnimplementedDataServiceServer) QueryRuns(context.Context, *QueryRunsRequest) (*QueryRunsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRuns not implemented")
}
func (UnimplementedDataServiceServer) SummarizeRuns(context.Context, *SummarizeRunsRequest) (*SummarizeRunsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SummarizeRuns not implemented")
}
func (UnimplementedDataServiceServer) GetRunner(context.Context, *GetRunnerRequest) (*GetRunnerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRunner not implemented")
}
func (UnimplementedDataServiceServer) QueryRunners(context.Context, *QueryRunnersRequest) (*QueryRunnersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRunners not implemented")
}
func (UnimplementedDataServiceServer) mustEmbedUnimplementedDataServiceServer() {}

// UnsafeDataServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataServiceServer will
// result in compilation errors.
type UnsafeDataServiceServer interface {
	mustEmbedUnimplementedDataServiceServer()
}

func RegisterDataServiceServer(s grpc.ServiceRegistrar, srv DataServiceServer) {
	s.RegisterService(&DataService_ServiceDesc, srv)
}

func _DataService_GetTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).GetTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tstr.data.v1.DataService/GetTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).GetTest(ctx, req.(*GetTestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataService_QueryTests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).QueryTests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tstr.data.v1.DataService/QueryTests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).QueryTests(ctx, req.(*QueryTestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataService_GetRun_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRunRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).GetRun(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tstr.data.v1.DataService/GetRun",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).GetRun(ctx, req.(*GetRunRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataService_QueryRuns_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRunsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).QueryRuns(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tstr.data.v1.DataService/QueryRuns",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).QueryRuns(ctx, req.(*QueryRunsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataService_SummarizeRuns_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SummarizeRunsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).SummarizeRuns(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tstr.data.v1.DataService/SummarizeRuns",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).SummarizeRuns(ctx, req.(*SummarizeRunsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataService_GetRunner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRunnerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).GetRunner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tstr.data.v1.DataService/GetRunner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).GetRunner(ctx, req.(*GetRunnerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataService_QueryRunners_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRunnersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).QueryRunners(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tstr.data.v1.DataService/QueryRunners",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).QueryRunners(ctx, req.(*QueryRunnersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DataService_ServiceDesc is the grpc.ServiceDesc for DataService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tstr.data.v1.DataService",
	HandlerType: (*DataServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTest",
			Handler:    _DataService_GetTest_Handler,
		},
		{
			MethodName: "QueryTests",
			Handler:    _DataService_QueryTests_Handler,
		},
		{
			MethodName: "GetRun",
			Handler:    _DataService_GetRun_Handler,
		},
		{
			MethodName: "QueryRuns",
			Handler:    _DataService_QueryRuns_Handler,
		},
		{
			MethodName: "SummarizeRuns",
			Handler:    _DataService_SummarizeRuns_Handler,
		},
		{
			MethodName: "GetRunner",
			Handler:    _DataService_GetRunner_Handler,
		},
		{
			MethodName: "QueryRunners",
			Handler:    _DataService_QueryRunners_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "data/v1/data.proto",
}
