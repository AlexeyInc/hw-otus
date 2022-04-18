package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	schedulerConfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	amqpClient "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp"
	amqpModels "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp/models"
	sqlstorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
	sqlc "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql/sqlc"
	domainModels "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"
)

const _notificationInQueueStatus int32 = 1

var (
	configFile   = flag.String("config", "../../../configs/scheduler_config.toml", "Path to configuration file")
	exchangeName = flag.String("exchange", "calendar-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchangeType", "direct", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("routingKey", "notification-key", "AMQP routing key")
	bindingKey   = flag.String("bindingKey", "notification-key", "AMQP binding key")
	queueName    = flag.String("queueName", "event-notification-queue", "AMQP queue name")
)

type AMQPClient interface {
	InitConnectionAndChannel() error
	Publish(payload []byte, exchangeName, routingKey string) error
	DeclareExchange(exchangeName, exchangeKind string) error
	DeclareQueue(queueName string) error
	BindQueue(exchangeName, queueName, bindKey string) error
	Shutdown()
}

type Scheduler struct {
	storage                       *sqlstorage.Storage
	amqpClient                    AMQPClient
	checkNotificationFreqSeconds  int
	checkExpiredEventsFreqSeconds int
}

func newScheduler(c schedulerConfig.Config) *Scheduler {
	return &Scheduler{
		storage: &sqlstorage.Storage{
			Driver: c.Storage.Driver,
			Source: c.Storage.Source,
		},
		amqpClient: &amqpClient.AMQPManager{
			AmqpURI: c.AMQP.Source,
		},
		checkNotificationFreqSeconds:  c.Scheduler.CheckNotificationFreqSeconds,
		checkExpiredEventsFreqSeconds: c.Scheduler.CheckExpiredEventsFreqSeconds,
	}
}

func (scheduler *Scheduler) setupAMQP() {
	err := scheduler.amqpClient.InitConnectionAndChannel()
	failOnError(err, "can't read config file")
	log.Printf("AMQP Connection and Channel initialized")

	err = scheduler.amqpClient.DeclareExchange(*exchangeName, *exchangeType)
	failOnError(err, "failed to declare the Exchange")
	log.Printf("Declared Exchange :%s", *exchangeName)

	err = scheduler.amqpClient.DeclareQueue(*queueName)
	failOnError(err, "failed to declare the Queue")
	log.Printf("Declared Queue :%s", *queueName)

	err = scheduler.amqpClient.BindQueue(*exchangeName, *queueName, *bindingKey)
	failOnError(err, "failed to bind to the Queue")
	log.Printf("Queue %s bound to %s with bindingKey %s", *queueName, *exchangeName, *bindingKey)
}

func (scheduler *Scheduler) getEventNotifications(ctx context.Context) (eventModels []domainModels.Event, err error) {
	now := time.Now()
	events, err := scheduler.storage.DbQueries.GetNotifyEvents(ctx, now)
	if err != nil {
		return eventModels, err
	}
	return toViewModels(events), err
}

func (scheduler *Scheduler) updateNotificationStatus(ctx context.Context, eventID int64) error {
	_, err := scheduler.storage.DbQueries.UpdateEventNotificationStatus(
		ctx, sqlc.UpdateEventNotificationStatusParams{
			Notificationstatus: sql.NullInt32{Int32: _notificationInQueueStatus, Valid: true},
			ID:                 eventID,
		})
	if err != nil {
		return err
	}
	return nil
}

func (scheduler *Scheduler) deleteExpiredEvents(ctx context.Context) (err error) {
	return scheduler.storage.DbQueries.DeleteExpiredEvents(ctx)
}

func main() {
	flag.Parse()

	log.Println("Starting publisher...")

	config, err := schedulerConfig.NewConfig(*configFile)
	failOnError(err, "can't read config file")

	scheduler := newScheduler(config)
	defer scheduler.amqpClient.Shutdown()

	scheduler.setupAMQP()

	context, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if err := scheduler.storage.Connect(context); err != nil {
		failOnError(err, "can't connect to database")
		cancel()
	}

	go proccesEventNotifications(context, scheduler)

	go deleteExpiredEvents(context, scheduler)

	<-context.Done()

	log.Println("\nFinishing publishing...")
}

func proccesEventNotifications(context context.Context, scheduler *Scheduler) {
	for {
		log.Printf("Check on event notifications")

		events, err := scheduler.getEventNotifications(context)
		failOnError(err, "error during getting event notifications")

		if len(events) > 0 {
			log.Printf("Got some events, count: %d\n", len(events))

			for _, event := range events {
				err := scheduler.updateNotificationStatus(context, event.ID)
				failOnError(err, "failed to udpate event notification status")

				notification := amqpModels.Notification{
					IdEvent:    event.ID,
					EventTitle: event.Title,
					EventStart: event.StartEvent,
					IdUser:     event.IDUser,
				}

				payload, err := json.Marshal(notification)
				failOnError(err, "Failed to marshal JSON")

				err = scheduler.amqpClient.Publish(payload, *exchangeName, *routingKey)
				failOnError(err, "Failed to public message")
			}
		}
		time.Sleep(time.Duration(scheduler.checkNotificationFreqSeconds) * time.Second)
	}
}

func deleteExpiredEvents(context context.Context, scheduler *Scheduler) {
	for {
		log.Println("Check on expired events")
		err := scheduler.deleteExpiredEvents(context)
		failOnError(err, "error during deleting expired events")

		time.Sleep(time.Duration(scheduler.checkExpiredEventsFreqSeconds) * time.Second)
	}
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
