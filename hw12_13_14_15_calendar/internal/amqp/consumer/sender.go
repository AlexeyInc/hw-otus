package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	senderConfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	amqpClient "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp"
	amqpModels "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp/models"
	sqlstorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
	sqlc "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql/sqlc"
	"github.com/streadway/amqp"
)

const _notificationSendedStatus int32 = 2

var replies <-chan amqp.Delivery

var (
	configFile   = flag.String("config", "../../../configs/calendar_config.toml", "Path to configuration file")
	queueName    = flag.String("queueName", "event-notification-queue", "AMQP queue name")
	consumerName = flag.String("consumer-name", "sender-consumer", "AMQP consumer name (should not be blank)")
)

type AMQPClient interface {
	InitConnectionAndChannel() error
	Consume(consumerName, queueName string) (<-chan amqp.Delivery, error)
	Shutdown()
}

type Sender struct {
	storage    *sqlstorage.Storage
	amqpClient AMQPClient
}

func newSender(c senderConfig.Config) *Sender {
	return &Sender{
		amqpClient: &amqpClient.AMQPManager{
			AmqpURI: c.AMQP.Source,
		},
		storage: &sqlstorage.Storage{
			Driver: c.Storage.Driver,
			Source: c.Storage.Source,
		},
	}
}

func (s *Sender) setupAMQP() {
	err := s.amqpClient.InitConnectionAndChannel()
	failOnError(err, "failed to initialize AMQP client")

	log.Printf("AMQP Connection and Channel initialized")

	replies, err = s.amqpClient.Consume(*consumerName, *queueName)
	failOnError(err, "failed to consume from specified queue")
}

func main() {
	flag.Parse()

	log.Println("Start consuming the Queue...")

	config, err := senderConfig.NewConfig(*configFile)
	failOnError(err, "can't read config file")

	sender := newSender(config)
	defer sender.amqpClient.Shutdown()

	sender.setupAMQP()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if err := sender.storage.Connect(ctx); err != nil {
		failOnError(err, "can't connect to database")
		cancel()
	}
	go processReceivedMessages(ctx, sender)

	<-ctx.Done()

	log.Println("\nStopped consuming the Queue...")
}

func processReceivedMessages(context context.Context, sender *Sender) {
	count := 1

	for r := range replies {
		log.Printf("Consuming reply number %d", count)

		v := amqpModels.Notification{}
		err := json.Unmarshal(r.Body, &v)
		failOnError(err, "can't unmarshal response body")

		err = sender.updateNotificationStatus(context, v.IdEvent)
		failOnError(err, "can't update notification status")

		fmt.Printf("\nIdEvent: %d,\nTitle: %s\nEventStart: %s\nIdUser: %d",
			v.IdEvent, v.EventTitle, v.EventStart.String(), v.IdUser)
		count++
	}
}

func (s *Sender) updateNotificationStatus(ctx context.Context, eventID int64) error {
	_, err := s.storage.DbQueries.UpdateEventNotificationStatus(
		ctx, sqlc.UpdateEventNotificationStatusParams{
			Notificationstatus: sql.NullInt32{Int32: _notificationSendedStatus, Valid: true},
			ID:                 eventID,
		})
	if err != nil {
		return err
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
