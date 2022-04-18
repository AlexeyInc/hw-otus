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

var _notificationSendedStatus int32 = 2

type AMQPClient interface {
	InitConnectionAndChannel() error
	Consume(consumerName, queueName string) (<-chan amqp.Delivery, error)
	Shutdown()
}

type Sender struct {
	storage    *sqlstorage.Storage
	amqpClient AMQPClient
}

var replies <-chan amqp.Delivery

var (
	configFile   = flag.String("config", "../../../configs/calendar_config.toml", "Path to configuration file")
	queueName    = flag.String("queueName", "event-notification-queue", "AMQP queue name")
	consumerName = flag.String("consumer-name", "sender-consumer", "AMQP consumer name (should not be blank)")
)

func init() {
	flag.Parse()
}

func New(c senderConfig.Config) *Sender {
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
	failOnError(err, "Failed to initialize to AMQP client")

	log.Printf("AMQP Connection and Channel initialized")

	replies, err = s.amqpClient.Consume(*consumerName, *queueName)
	failOnError(err, "Failed to consume message")
}

func main() {
	log.Println("Start consuming the Queue...")

	config, err := senderConfig.NewConfig(*configFile)
	failOnError(err, "can't read config file")

	sender := New(config)
	defer sender.amqpClient.Shutdown()

	sender.setupAMQP()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if err := sender.storage.Connect(ctx); err != nil {
		failOnError(err, "can't connect to database")
		cancel()
	}
	go ProcessReceivedMessages(ctx, sender)

	<-ctx.Done()

	log.Println("\nStopped consuming the Queue...")
}

func ProcessReceivedMessages(context context.Context, sender *Sender) {
	count := 1

	for r := range replies {
		log.Printf("Consuming reply number %d", count)
		v := amqpModels.Notification{}
		json.Unmarshal(r.Body, &v)

		sender.UpdateNotificationStatus(context, v.IdEvent)

		fmt.Printf("\nIdEvent: %d,\nTitle: %s\nEventStart: %s\nIdUser: %d",
			v.IdEvent, v.EventTitle, v.EventStart.String(), v.IdUser)
		count++
	}
}

func (s *Sender) UpdateNotificationStatus(ctx context.Context, eventID int64) error {
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
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
