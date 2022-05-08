package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	calendarconfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	app "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	sqlstorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile, logFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/calendar_config.toml", "Path to configuration file")
	flag.StringVar(&logFile, "log", "../../log/logs.log", "Path to log file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := calendarconfig.NewConfig(configFile)
	if err != nil {
		log.Println("can't read config file: " + err.Error())
		return
	}

	zapLogg := logger.New(logFile, config.Logger.Level)
	defer zapLogg.ZapLogger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	storage := sqlstorage.New(config)
	if err := storage.Connect(ctx); err != nil {
		zapLogg.Info("connection to database failed: " + err.Error())
		cancel()
		return
	}
	zapLogg.Info("Successfully connected to database...")
	defer storage.Close(ctx)

	calendar := app.New(zapLogg, storage)

	go internalgrpc.RunGRPCServer(ctx, config, calendar, zapLogg)

	go internalhttp.RunHTTPServer(ctx, config, calendar, zapLogg)

	log.Println("Calendar service started")

	<-ctx.Done()

	zapLogg.Info("\nAll servers are stopped...")
}
