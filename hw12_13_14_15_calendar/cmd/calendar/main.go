package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	app "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/services"

	//memorystorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile, logFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.toml", "Path to configuration file")
	flag.StringVar(&logFile, "log", "../../internal/logger/requests.log", "Path to log file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := configs.NewConfig(configFile)
	if err != nil {
		log.Fatalln("can't read config file: " + err.Error())
		os.Exit(1)
	}

	zapLogg := logger.New(logFile, config.Logger.Level)
	defer zapLogg.ZapLogger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := sqlstorage.New(config)
	if err := storage.Connect(ctx); err != nil {
		log.Fatalln("connection to database failed: " + err.Error())
		os.Exit(1)
	}
	defer storage.Close(ctx)

	calendar := app.New(zapLogg, storage)

	// Run gRPC Server...

	go internalgrpc.RunGRPCServer(ctx, config, calendar, zapLogg)

	// Run HTTP Server...

	go internalhttp.RunHTTPServer(ctx, config, calendar, zapLogg)

	<-ctx.Done()

	println("\nAll servers are stopped...")
	cancel()
	os.Exit(0)
}
