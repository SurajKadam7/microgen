// Code generated by microgen 0.9.0. DO NOT EDIT.

package transport

import (
	"context"
	endpoint "github.com/go-kit/kit/endpoint"
	opentracing "github.com/go-kit/kit/tracing/opentracing"
	opentracinggo "github.com/opentracing/opentracing-go"
	addsvc "github.com/recolabs/microgen/examples/addsvc/addsvc"
)

func Endpoints(svc addsvc.Service) EndpointsSet {
	return EndpointsSet{
		ConcatEndpoint: ConcatEndpoint(svc),
		SumEndpoint:    SumEndpoint(svc),
	}
}

// TraceServerEndpoints is used for tracing endpoints on server side.
func TraceServerEndpoints(endpoints EndpointsSet, tracer opentracinggo.Tracer) EndpointsSet {
	return EndpointsSet{
		ConcatEndpoint: opentracing.TraceServer(tracer, "Concat")(endpoints.ConcatEndpoint),
		SumEndpoint:    opentracing.TraceServer(tracer, "Sum")(endpoints.SumEndpoint),
	}
}

func SumEndpoint(svc addsvc.Service) endpoint.Endpoint {
	return func(arg0 context.Context, request interface{}) (interface{}, error) {
		req := request.(*SumRequest)
		res0, res1 := svc.Sum(arg0, req.A, req.B)
		return &SumResponse{Result: res0}, res1
	}
}

func ConcatEndpoint(svc addsvc.Service) endpoint.Endpoint {
	return func(arg0 context.Context, request interface{}) (interface{}, error) {
		req := request.(*ConcatRequest)
		res0, res1 := svc.Concat(arg0, req.A, req.B)
		return &ConcatResponse{Result: res0}, res1
	}
}
