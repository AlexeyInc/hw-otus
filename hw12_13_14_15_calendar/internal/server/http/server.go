package internalhttp

import (
	"context"
	"log"
	"net/http"

	api "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/api/protoc"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Server struct {
	Host string
	Port string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func RunHTTPServer(context context.Context, config configs.Config, app api.EventServiceServer, logger Logger) (err error) {
	mux := runtime.NewServeMux()

	err = api.RegisterEventServiceHandlerServer(context, mux, app)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:    config.HTTPServer.Host + config.HTTPServer.Port,
		Handler: addLoggingMiddleware(logger, mux),
	}

	go func() {
		<-context.Done()

		if err := s.Shutdown(context); err != nil {
			log.Fatal("Failed to shutdown http server: ", err)
		}
	}()

	logger.Info("calendar HTTP server is running...")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("Failed to listen and serve: ", err)
		return err
	}
	return
}
