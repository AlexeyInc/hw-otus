package internalgrpc

import (
	"context"
	"log"
	"net"

	api "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/api/protoc"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	"google.golang.org/grpc"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type GRPCServer struct {
	App api.EventServiceServer
	api.UnimplementedEventServiceServer
}

func RunGRPCServer(context context.Context, config configs.Config, app api.EventServiceServer, logger Logger) (err error) {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(addLoggingMiddleware(logger)),
	)

	api.RegisterEventServiceServer(gRPCServer, app)

	l, err := net.Listen(config.GRPCServer.Network, config.GRPCServer.Host+config.GRPCServer.Port)
	if err != nil {
		log.Fatal("can't run listener: ", err)
	}

	logger.Info("calendar gRPC server is running...")
	if err = gRPCServer.Serve(l); err != nil {
		log.Fatal("can't run server: ", err)
	}
	return
}
