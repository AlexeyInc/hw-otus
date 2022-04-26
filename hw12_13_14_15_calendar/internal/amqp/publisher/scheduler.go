package publisher

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	schedulerConfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	amqpManager "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp"
	amqpModels "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp/models"
	sqlstorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
	sqlc "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql/sqlc"
	domainModels "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"
)

type AMQPManager interface {
	InitConnectionAndChannel() error
	Publish(payload []byte, exchangeName, routingKey string) error
	DeclareExchange(exchangeName, exchangeKind string) error
	DeclareQueue(queueName string) error
	BindQueue(exchangeName, queueName, bindKey string) error
	Shutdown()
}

type Scheduler struct {
	AMQPManager
	*sqlstorage.Storage
	checkNotificationFreqMinutes  int
	checkExpiredEventsFreqMinutes int
}

func New(c schedulerConfig.Config) *Scheduler {
	return &Scheduler{
		Storage: &sqlstorage.Storage{
			Driver: c.Storage.Driver,
			Source: c.Storage.Source,
		},
		AMQPManager: &amqpManager.AMQPManager{
			AMQPURI: c.AMQP.Source,
		},
		checkNotificationFreqMinutes:  c.Scheduler.CheckNotificationFreqMinutes,
		checkExpiredEventsFreqMinutes: c.Scheduler.CheckExpiredEventsFreqMinutes,
	}
}

func (s *Scheduler) SetupAMQP(exchangeName, exchangeType, queueName, bindingKey string) {
	err := s.AMQPManager.InitConnectionAndChannel()
	failOnError(err, "Failed to initialize to AMQP client")
	log.Printf("AMQP Connection and Channel initialized")

	s.AMQPManager.DeclareExchange(exchangeName, exchangeType)
	failOnError(err, "Failed to declare the Exchange")
	log.Printf("Declared Exchange :%s", exchangeName)

	s.AMQPManager.DeclareQueue(queueName)
	failOnError(err, "Failed to declare the Queue")
	log.Printf("Declared Queue :%s", queueName)

	s.AMQPManager.BindQueue(exchangeName, queueName, bindingKey)
	failOnError(err, "Failed to bind to the Queue")
	log.Printf("Queue %s bound to %s with bindingKey %s", queueName, exchangeName, bindingKey)
}

func (s *Scheduler) GetEventNotifications(ctx context.Context) (eventModels []domainModels.Event, err error) {
	now := time.Now()
	events, err := s.Storage.DBQueries.GetNotifyEvents(ctx, now)
	if err != nil {
		return eventModels, err
	}
	return toViewModels(events), err
}

func (s *Scheduler) ProccesEventNotifications(context context.Context, exchangeName, routingKey string) {
	for {
		events, err := s.GetEventNotifications(context)
		failOnError(err, "error during getting event notifications")

		if len(events) > 0 {
			log.Printf("Got some events, count: %d\n", len(events))

			for _, event := range events {
				err := s.updateNotificationStatus(context, event.ID)
				failOnError(err, "error during udpating event notification status")

				notification := amqpModels.Notification{
					IDEvent:    event.ID,
					EventTitle: event.Title,
					EventStart: event.StartEvent,
					IDUser:     event.IDUser,
				}

				payload, err := json.Marshal(notification)
				failOnError(err, "Failed to marshal JSON")

				err = s.AMQPManager.Publish(payload, exchangeName, routingKey)
				failOnError(err, "Failed to public message")
			}
		}
		time.Sleep(time.Duration(s.checkNotificationFreqMinutes) * time.Minute)
	}
}

func (s *Scheduler) DeleteExpiredEvents(context context.Context) {
	for {
		log.Println("Check on expired events")
		err := s.Storage.DBQueries.DeleteExpiredEvents(context)
		failOnError(err, "error during deleting expired events")

		time.Sleep(time.Duration(s.checkExpiredEventsFreqMinutes) * time.Minute)
	}
}

func (s *Scheduler) updateNotificationStatus(ctx context.Context, eventID int64) error {
	_, err := s.Storage.DBQueries.UpdateEventNotificationStatus(
		ctx, sqlc.UpdateEventNotificationStatusParams{
			Notificationsended: sql.NullBool{Bool: true, Valid: true},
			ID:                 eventID,
		})
	if err != nil {
		return err
	}
	return nil
}

func toViewModel(ev sqlc.Event) domainModels.Event {
	return domainModels.Event{
		ID:           ev.ID,
		Title:        ev.Title,
		StartEvent:   ev.StartEvent,
		EndEvent:     ev.EndEvent,
		Description:  ev.Description.String,
		IDUser:       ev.IDUser,
		Notification: ev.Notification.Time,
	}
}

func toViewModels(events []sqlc.Event) []domainModels.Event {
	result := make([]domainModels.Event, len(events))
	for i := 0; i < len(events); i++ {
		result[i] = toViewModel(events[i])
	}
	return result
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
