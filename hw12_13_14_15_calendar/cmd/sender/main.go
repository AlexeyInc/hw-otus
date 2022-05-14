package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	senderConfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	consumer "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp/consumer"
)

var (
	configFile   = flag.String("config", "../../configs/calendar_config.toml", "Path to configuration file")
	queueName    = flag.String("queueName", "event-notification-queue", "AMQP queue name")
	consumerName = flag.String("consumer-name", "sender-consumer", "AMQP consumer name (should not be blank)")
)

func init() {
	flag.Parse()
}

func main() {
	log.Println("Start consuming the Queue...")

	config, err := senderConfig.NewConfig(*configFile)
	failOnError(err, "can't read config file")

	sender := consumer.New(config)
	defer sender.AMQPManager.Shutdown()

	sender.SetupAMQP(*consumerName, *queueName)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go sender.ProcessReceivedMessages()

	<-ctx.Done()

	log.Println("\nStopped consuming the Queue...")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
