package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	senderConfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	amqpManager "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp"
	amqpModels "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp/models"
	sqlstorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
	sqlc "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql/sqlc"
	"github.com/streadway/amqp"
)

const _notificationSendedStatus int32 = 2

type AMQPManager interface {
	InitConnectionAndChannel() error
	Consume(consumerName, queueName string) (<-chan amqp.Delivery, error)
	Shutdown()
}

type Sender struct {
	AMQPManager
	*sqlstorage.Storage
}

var replies <-chan amqp.Delivery

func New(c senderConfig.Config) *Sender {
	return &Sender{
		AMQPManager: &amqpManager.AMQPManager{
			AMQPURI: c.AMQP.Source,
		},
		Storage: &sqlstorage.Storage{
			Driver: c.Storage.Driver,
			Source: c.Storage.Source,
		},
	}
}

func (s *Sender) SetupAMQP(consumerName, queueName string) {
	err := s.AMQPManager.InitConnectionAndChannel()
	failOnError(err, "Failed to initialize to AMQP client")

	log.Printf("AMQP Connection and Channel initialized")

	replies, err = s.AMQPManager.Consume(consumerName, queueName)
	failOnError(err, "Failed to consume message")
}

func (s *Sender) ProcessReceivedMessages(ctx context.Context) {
	count := 1
	for r := range replies {
		log.Printf("Consuming event number %d", count)
		v := amqpModels.Notification{}
		err := json.Unmarshal(r.Body, &v)
		failOnError(err, "can't unmarshal response body")

		err = s.updateNotificationStatus(context.Background(), v.IDEvent)
		failOnError(err, "can't update notification status")

		fmt.Printf("\nIdEvent: %d,\nTitle: %s\nEventStart: %s\nIdUser: %d",
			v.IDEvent, v.EventTitle, v.EventStart.String(), v.IDUser)
		count++
	}
}

func (s *Sender) updateNotificationStatus(ctx context.Context, eventID int64) error {
	log.Println("here 1", eventID)
	log.Println("1", s)
	log.Println("2", s.Storage)
	log.Println("3", s.Storage.DBQueries)

	_, err := s.Storage.DBQueries.UpdateEventNotificationStatus(
		ctx, sqlc.UpdateEventNotificationStatusParams{
			Notificationstatus: sql.NullInt32{Int32: _notificationSendedStatus, Valid: true},
			ID:                 eventID,
		})
	log.Printf("here 3")

	if err != nil {
		log.Printf("here 2")

		return err
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
