package middleware

import (
	"context"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type TransactionIdMarker string

func UnaryServerInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		id := uuid.NewString()

		logger.Printf("TransactionId created: %v", id)
		newCtx = context.WithValue(ctx, TransactionIdMarker("transaction_id_ctx_marker"), id)

		return handler(newCtx, req)
	}
}

func StreamServerInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		id := uuid.NewString()

		logger.Println("TransactionId created: %v", id)
		newCtx = context.WithValue(stream.Context(), TransactionIdMarker("transaction_id_ctx_marker"), id)

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
