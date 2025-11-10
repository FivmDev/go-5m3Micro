package server_interceptors

import (
	"context"
	"go-5m3Micro/pkg/log"
	"google.golang.org/grpc"
	"runtime/debug"
)

func UnaryCrashInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v \n \n %s", r, debug.Stack())
	})
	return handler(ctx, req)
}

func StreamCrashInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v \n \n %s", r, debug.Stack())
	})
	return handler(srv, ss)
}

func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}
