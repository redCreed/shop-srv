package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"runtime/debug"
)

func CrashInterceptor() grpc.UnaryServerInterceptor  {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	 	defer handleCrash(func(i interface{}) {
	 		err = toPanicError(i)
		})
		return handler(ctx, req)
	}
}

func handleCrash(handle func(interface{}))  {
	if r := recover(); r != nil {
		handle(r)
	}
}

func toPanicError(r interface{}) error  {
	//todo  日志记录
	log.Printf("%+v %s", r, debug.Stack())
	return status.Errorf(codes.Internal, "panic: %v ", r)
}