// Code generated by microgen 0.9.0. DO NOT EDIT.

package transportgrpc

import (
	log "github.com/go-kit/kit/log"
	opentracing "github.com/go-kit/kit/tracing/opentracing"
	grpckit "github.com/go-kit/kit/transport/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	opentracinggo "github.com/opentracing/opentracing-go"
	pb "github.com/recolabs/microgen/examples/protobuf"
	transport "github.com/recolabs/microgen/examples/usersvc/pkg/transport"
	grpc "google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, addr string, opts ...grpckit.ClientOption) transport.EndpointsSet {
	return transport.EndpointsSet{
		CreateCommentEndpoint: grpckit.NewClient(
			conn, addr, "CreateComment",
			_Encode_CreateComment_Request,
			_Decode_CreateComment_Response,
			pb.CreateCommentResponse{},
			opts...,
		).Endpoint(),
		CreateUserEndpoint: grpckit.NewClient(
			conn, addr, "CreateUser",
			_Encode_CreateUser_Request,
			_Decode_CreateUser_Response,
			pb.CreateUserResponse{},
			opts...,
		).Endpoint(),
		FindUsersEndpoint: grpckit.NewClient(
			conn, addr, "FindUsers",
			_Encode_FindUsers_Request,
			_Decode_FindUsers_Response,
			pb.FindUsersResponse{},
			opts...,
		).Endpoint(),
		GetCommentEndpoint: grpckit.NewClient(
			conn, addr, "GetComment",
			_Encode_GetComment_Request,
			_Decode_GetComment_Response,
			pb.GetCommentResponse{},
			opts...,
		).Endpoint(),
		GetUserCommentsEndpoint: grpckit.NewClient(
			conn, addr, "GetUserComments",
			_Encode_GetUserComments_Request,
			_Decode_GetUserComments_Response,
			pb.GetUserCommentsResponse{},
			opts...,
		).Endpoint(),
		GetUserEndpoint: grpckit.NewClient(
			conn, addr, "GetUser",
			_Encode_GetUser_Request,
			_Decode_GetUser_Response,
			pb.GetUserResponse{},
			opts...,
		).Endpoint(),
		UpdateUserEndpoint: grpckit.NewClient(
			conn, addr, "UpdateUser",
			_Encode_UpdateUser_Request,
			_Decode_UpdateUser_Response,
			empty.Empty{},
			opts...,
		).Endpoint(),
	}
}

func TracingGRPCClientOptions(tracer opentracinggo.Tracer, logger log.Logger) func([]grpckit.ClientOption) []grpckit.ClientOption {
	return func(opts []grpckit.ClientOption) []grpckit.ClientOption {
		return append(opts, grpckit.ClientBefore(
			opentracing.ContextToGRPC(tracer, logger),
		))
	}
}
