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

type Server struct {
	gRPCServer *grpc.Server
	listener   net.Listener
}

func RunGRPCServer(context context.Context, config configs.Config, app api.EventServiceServer, logger Logger) {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(addLoggingMiddleware(logger)),
	)

	api.RegisterEventServiceServer(gRPCServer, app)

	l, err := net.Listen(config.GRPCServer.Network, config.GRPCServer.Host+config.GRPCServer.Port)
	if err != nil {
		log.Fatal("can't run listener: ", err)
	}

	go func() {
		<-context.Done()

		gRPCServer.GracefulStop()
	}()

	logger.Info("calendar gRPC server is running...")
	if err = gRPCServer.Serve(l); err != nil {
		log.Fatal("can't run server: ", err)
	}
}

func NewServer(context context.Context, config configs.Config, app api.EventServiceServer, logger Logger) *Server {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(addLoggingMiddleware(logger)),
	)

	api.RegisterEventServiceServer(gRPCServer, app)

	l, err := net.Listen(config.GRPCServer.Network, config.GRPCServer.Host+config.GRPCServer.Port)
	if err != nil {
		log.Fatal("can't run listener: ", err)
	}

	return &Server{
		gRPCServer: gRPCServer,
		listener:   l,
	}
}

func (s *Server) Start(l net.Listener) (err error) {
	//logger.Info("calendar gRPC server is running...")
	if err = s.gRPCServer.Serve(l); err != nil {
		log.Fatal("can't run server: ", err)
	}
	return err
}

func (s *Server) Stop() {
	s.gRPCServer.GracefulStop()
}
