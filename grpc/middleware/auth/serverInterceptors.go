package grpc_reqAuth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
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

		errorResp := fmt.Errorf("%v | %v", tId, "auth issue")

		logger.Info("starting req auth...")
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Info("no metadata")
			return nil, errorResp
		}

		// create hash
		var hmac_message string
		for _, v := range opts.MetadataKeyList {
			if len(md.Get(v)) == 0 {
				logger.Info("rpc request does not contain this metadata: ", v)
				return nil, errorResp
			}
			hmac_message = fmt.Sprintf("%v%v", hmac_message, md.Get(v)[0])
		}
		logger.Info("hmac message created")
		mac := hmac.New(sha256.New, []byte(opts.HashSecret))
		mac.Write([]byte(hmac_message))
		expectedMAC := mac.Sum(nil)

		// verify hash matches
		logger.Info("verify hash matches")
		if md.Get("hypebid-nohash")[0] == "false" &&
			len(md.Get(opts.MetadataHashKey)) != 0 &&
			hmac.Equal([]byte(md.Get(opts.MetadataHashKey)[0]), expectedMAC) {
			logger.Info("hmac hash does match")
			return handler(ctx, req)
		}
		logger.Info("hmac hash does not match")
		return nil, errorResp
	}
}

func StreamServerInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		return nil
	}
}
