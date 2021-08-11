package grpc_reqAuth

import (
	"context"
	"errors"
	"runtime"

	"github.com/hypebid/go-kit/grpc/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor(log *logrus.Logger, opts Options) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get transaction Id from ctx and build method logger
		tId := ctx.Value(middleware.Grpc_ReqId_Marker)
		pc, _, _, _ := runtime.Caller(0)
		logger := log.WithFields(logrus.Fields{"transaction-id": tId, "method": runtime.FuncForPC(pc).Name()})

		logger.Info("starting req auth...")
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("no metadata")
		}
		logger.Info("metadata: ", md)

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		return nil
	}
}
