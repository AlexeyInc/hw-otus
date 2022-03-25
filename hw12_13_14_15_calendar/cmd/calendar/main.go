package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
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
	defer zapLogg.SugarLogger.Sync()

	storage := memorystorage.New(config)

	calendar := app.New(zapLogg, storage)

	server := internalhttp.NewServer(zapLogg, config, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			zapLogg.Error("failed to stop http server: " + err.Error())
		}
	}()

	zapLogg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		zapLogg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
