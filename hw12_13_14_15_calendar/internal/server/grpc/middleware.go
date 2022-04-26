package internalgrpc

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func addLoggingMiddleware(logger Logger) grpc.UnaryServerInterceptor {
	return grpc.UnaryServerInterceptor(func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		p, _ := peer.FromContext(ctx)
		mD, _ := metadata.FromIncomingContext(ctx)

		clientID := p.Addr.String()
		method := info.FullMethod
		latency := time.Since(start)
		userAgent := mD["user-agent"]

		h, err := handler(ctx, req)

		statucCode := status.Code(err).String()
		contentLength, err := getGRPCResponseSize(h)
		if err != nil {
			log.Fatal("encode gRPC Response error:", err)
			return nil, err
		}

		msg := fmt.Sprintf("gRPC request has been made...\nClientIP:%s;\nMethod:%s;\nStatusCode:%s;"+
			"\nContentLength:%d;\nLatency:%s;\nUser_agent:%s",
			clientID, method, statucCode, contentLength, latency, userAgent)

		logger.Info(msg)

		return h, err
	})
}

func getGRPCResponseSize(val interface{}) (int, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(val); err != nil {
		return 0, err
	}
	return binary.Size(buff.Bytes()), nil
}
