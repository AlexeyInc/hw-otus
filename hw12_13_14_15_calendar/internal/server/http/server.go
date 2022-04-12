package internalhttp

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
)

type Server struct {
	Host string
	Port string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func myHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Hello world!")
}

func NewServer(logger Logger, config configs.Config, app Application) *Server {
	http.Handle("/", loggingMiddleware(logger, http.HandlerFunc(myHandler)))

	return &Server{
		Host: config.Server.Host,
		Port: config.Server.Port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	http.ListenAndServe(s.Host+s.Port, nil)

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	os.Exit(1)
	return nil
}

// TODO
