package grpc_reqAuth

import (
	"context"
	"runtime"

	"github.com/hypebid/go-kit/grpc/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(log *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get transaction Id from ctx and build method logger
		tId := ctx.Value(middleware.Grpc_ReqId_Marker)
		pc, _, _, _ := runtime.Caller(0)
		logger := log.WithFields(logrus.Fields{"transaction-id": tId, "method": runtime.FuncForPC(pc).Name()})

		logger.Info("starting req auth...")

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		return nil
	}
}
