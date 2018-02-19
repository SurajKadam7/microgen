// This file was automatically generated by "microgen 0.8.0alpha" utility.
// This file will never be overwritten.
package main

import (
	errors "errors"
	generated "github.com/devimteam/microgen/example/generated"
	middleware "github.com/devimteam/microgen/example/generated/middleware"
	grpc "github.com/devimteam/microgen/example/generated/transport/grpc"
	http "github.com/devimteam/microgen/example/generated/transport/http"
	protobuf "github.com/devimteam/microgen/example/protobuf"
	log "github.com/go-kit/kit/log"
	opentracinggo "github.com/opentracing/opentracing-go"
	grpc1 "google.golang.org/grpc"
	io "io"
	net "net"
	http1 "net/http"
	os "os"
	signal "os/signal"
	syscall "syscall"
)

func main() {
	logger := log.With(InitLogger(os.Stdout), "level", "info")
	errorLogger := log.With(InitLogger(os.Stderr), "level", "error")
	logger.Log("message", "Hello, I am alive")
	defer logger.Log("message", "goodbye, good luck")

	errorChan := make(chan error)
	go InterruptHandler(errorChan)

	service := generated.NewStringService()                      // Create new service.
	service = middleware.ServiceLogging(logger)(service)         // Setup service logging.
	service = middleware.ServiceErrorLogging(logger)(service)    // Setup error logging.
	service = middleware.ServiceRecovering(errorLogger)(service) // Setup service recovering.

	endpoints := generated.AllEndpoints(service,
		opentracinggo.NoopTracer{}, // TODO: Add tracer
	)

	grpcAddr := ":8081"
	// Start grpc server.
	go ServeGRPC(endpoints, errorChan, grpcAddr, log.With(logger, "transport", "GRPC"))

	httpAddr := ":8080"
	// Start http server.
	go ServeHTTP(endpoints, errorChan, httpAddr, log.With(logger, "transport", "HTTP"))

	logger.Log("error", <-errorChan)
}

// InitLogger initialize go-kit JSON logger with timestamp and caller.
func InitLogger(writer io.Writer) log.Logger {
	logger := log.NewJSONLogger(writer)
	logger = log.With(logger, "@timestamp", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}

// InterruptHandler handles first SIGINT and SIGTERM and sends messages to error channel.
func InterruptHandler(ch chan<- error) {
	interruptHandler := make(chan os.Signal, 1)
	signal.Notify(interruptHandler, syscall.SIGINT, syscall.SIGTERM)
	ch <- errors.New((<-interruptHandler).String())
}

// ServeGRPC starts new GRPC server on address and sends first error to channel.
func ServeGRPC(endpoints *generated.Endpoints, ch chan<- error, addr string, logger log.Logger) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		ch <- err
		return
	}
	// Here you can add middlewares for grpc server.
	server := grpc.NewGRPCServer(endpoints,
		logger,
		opentracinggo.NoopTracer{}, // TODO: Add tracer
	)
	grpcServer := grpc1.NewServer()
	protobuf.RegisterStringServiceServer(grpcServer, server)
	logger.Log("listen on", addr)
	ch <- grpcServer.Serve(listener)
}

// ServeHTTP starts new HTTP server on address and sends first error to channel.
func ServeHTTP(endpoints *generated.Endpoints, ch chan<- error, addr string, logger log.Logger) {
	handler := http.NewHTTPHandler(endpoints,
		logger,
		opentracinggo.NoopTracer{}, // TODO: Add tracer
	)
	httpServer := &http1.Server{
		Addr:    addr,
		Handler: handler,
	}
	logger.Log("listen on", addr)
	ch <- httpServer.ListenAndServe()
}
