package internalhttp

import (
	"context"
	"log"
	"net/http"

	api "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/api/protoc"
	calendarconfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Server struct {
	httpServer *http.Server
}

func RunHTTPServer(context context.Context, config calendarconfig.Config, app api.EventServiceServer, logger Logger) {
	mux := runtime.NewServeMux()

	err := api.RegisterEventServiceHandlerServer(context, mux, app)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:    config.HTTPServer.Host + config.HTTPServer.Port,
		Handler: addLoggingMiddleware(logger, mux),
	}

	go func() {
		logger.Info("calendar HTTP server is running...")
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("Failed to listen and serve: ", err)
		}
	}()

	<-context.Done()

	if err := s.Close(); err != nil {
		log.Fatal("failed to close http server: ", err)
	}
}

func NewServer(context context.Context,
	config calendarconfig.Config, app api.EventServiceServer, logger Logger) *Server {
	mux := runtime.NewServeMux()

	err := api.RegisterEventServiceHandlerServer(context, mux, app)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:    config.HTTPServer.Host + config.HTTPServer.Port,
		Handler: addLoggingMiddleware(logger, mux),
	}
	return &Server{
		httpServer: s,
	}
}

func (server *Server) Start(ctx context.Context) (err error) {
	if err = server.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("Failed to listen and serve: ", err)
		return
	}
	return
}

func (server *Server) Stop(ctx context.Context) (err error) {
	if err = server.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutdown http server: ", err)
	}
	return
}
