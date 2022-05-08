package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	schedulerConfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	publisher "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp/publisher"
)

var (
	configFile   = flag.String("config", "../../configs/scheduler_config.toml", "Path to configuration file")
	exchangeName = flag.String("exchange", "calendar-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchangeType", "direct", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("routingKey", "notification-key", "AMQP routing key")
	bindingKey   = flag.String("bindingKey", "notification-key", "AMQP binding key")
	queueName    = flag.String("queueName", "event-notification-queue", "AMQP queue name")
)

func init() {
	flag.Parse()
}

func main() {
	log.Println("Starting publisher...")

	config, err := schedulerConfig.NewConfig(*configFile)
	failOnError(err, "can't read config file")

	scheduler := publisher.New(config)
	defer scheduler.AMQPManager.Shutdown()

	scheduler.SetupAMQP(*exchangeName, *exchangeType, *queueName, *bindingKey)

	ctx, _ := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if err := scheduler.Storage.Connect(ctx); err != nil {
		log.Fatalf("%s: %s", "can't connect to database", err)
	}

	go scheduler.ProccesEventNotifications(ctx, *exchangeName, *routingKey)

	go scheduler.DeleteExpiredEvents(ctx)

	<-ctx.Done()

	log.Println("\nFinishing publishing...")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
