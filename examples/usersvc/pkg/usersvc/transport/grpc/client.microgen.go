// Code generated by microgen. DO NOT EDIT.

package grpc

import (
	pb "github.com/devimteam/microgen/examples/usersvc/pb"
	transport "github.com/devimteam/microgen/examples/usersvc/pkg/usersvc/transport"
	log "github.com/go-kit/kit/log"
	opentracing "github.com/go-kit/kit/tracing/opentracing"
	grpckit "github.com/go-kit/kit/transport/grpc"
	types "github.com/gogo/protobuf/types"
	empty "github.com/golang/protobuf/ptypes/empty"
	opentracinggo "github.com/opentracing/opentracing-go"
	grpc "google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, addr string, opts ...grpckit.ClientOption) transport.Endpoints {
	return transport.Endpoints{
		CreateComment_Endpoint: grpckit.NewClient(
			conn, addr, "CreateComment",
			_Encode_CreateComment_Request,
			_Decode_CreateComment_Response,
			new(types.StringValue),
			opts...,
		).Endpoint(),
		CreateUser_Endpoint: grpckit.NewClient(
			conn, addr, "CreateUser",
			_Encode_CreateUser_Request,
			_Decode_CreateUser_Response,
			new(types.StringValue),
			opts...,
		).Endpoint(),
		FindUsers_Endpoint: grpckit.NewClient(
			conn, addr, "FindUsers",
			_Encode_FindUsers_Request,
			_Decode_FindUsers_Response,
			new(pb.FindUsersResponse),
			opts...,
		).Endpoint(),
		GetComment_Endpoint: grpckit.NewClient(
			conn, addr, "GetComment",
			_Encode_GetComment_Request,
			_Decode_GetComment_Response,
			new(pb.GetCommentResponse),
			opts...,
		).Endpoint(),
		GetUserComments_Endpoint: grpckit.NewClient(
			conn, addr, "GetUserComments",
			_Encode_GetUserComments_Request,
			_Decode_GetUserComments_Response,
			new(pb.GetUserCommentsResponse),
			opts...,
		).Endpoint(),
		GetUser_Endpoint: grpckit.NewClient(
			conn, addr, "GetUser",
			_Encode_GetUser_Request,
			_Decode_GetUser_Response,
			new(pb.GetUserResponse),
			opts...,
		).Endpoint(),
		UpdateUser_Endpoint: grpckit.NewClient(
			conn, addr, "UpdateUser",
			_Encode_UpdateUser_Request,
			_Decode_UpdateUser_Response,
			new(empty.Empty),
			opts...,
		).Endpoint(),
	}
}

func ClientOptionsBuilder(opts []grpckit.ClientOption, fns ...func([]grpckit.ClientOption) []grpckit.ClientOption) []grpckit.ClientOption {
	for i := range fns {
		opts = fns[i](opts)
	}
	return opts
}

func TracingClientOptions(tracer opentracinggo.Tracer, logger log.Logger) func([]grpckit.ClientOption) []grpckit.ClientOption {
	return func(opts []grpckit.ClientOption) []grpckit.ClientOption {
		return append(opts, grpckit.ClientBefore(
			opentracing.ContextToGRPC(tracer, logger),
		))
	}
}
