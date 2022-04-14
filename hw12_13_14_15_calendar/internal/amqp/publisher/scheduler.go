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
	checkNotificationFreqMinutes  int
	checkExpiredEventsFreqMinutes int
}

var (
	configFile   = flag.String("config", "../../../configs/scheduler_config.toml", "Path to configuration file")
	exchangeName = flag.String("exchange", "calendar-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchangeType", "direct", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("routingKey", "notification-key", "AMQP routing key")
	bindingKey   = flag.String("bindingKey", "notification-key", "AMQP binding key")
	queueName    = flag.String("queueName", "event-notification-queue", "AMQP queue name")
)

func init() {
	flag.Parse()
}

func New(c schedulerConfig.Config) *Scheduler {
	return &Scheduler{
		storage: &sqlstorage.Storage{
			Driver: c.Storage.Driver,
			Source: c.Storage.Source,
		},
		amqpClient: &amqpClient.AMQPManager{
			AmqpURI: c.AMQP.Source,
		},
		checkNotificationFreqMinutes:  c.Scheduler.CheckNotificationFreqMinutes,
		checkExpiredEventsFreqMinutes: c.Scheduler.CheckExpiredEventsFreqMinutes,
	}
}

func (scheduler *Scheduler) setupAMQP() {
	err := scheduler.amqpClient.InitConnectionAndChannel()
	failOnError(err, "Failed to initialize to AMQP client")
	log.Printf("AMQP Connection and Channel initialized")

	scheduler.amqpClient.DeclareExchange(*exchangeName, *exchangeType)
	failOnError(err, "Failed to declare the Exchange")
	log.Printf("Declared Exchange :%s", *exchangeName)

	scheduler.amqpClient.DeclareQueue(*queueName)
	failOnError(err, "Failed to declare the Queue")
	log.Printf("Declared Queue :%s", *queueName)

	scheduler.amqpClient.BindQueue(*exchangeName, *queueName, *bindingKey)
	failOnError(err, "Failed to bind to the Queue")
	log.Printf("Queue %s bound to %s with bindingKey %s", *queueName, *exchangeName, *bindingKey)
}

func (scheduler *Scheduler) GetEventNotifications(ctx context.Context) (eventModels []domainModels.Event, err error) {
	now := time.Now()
	events, err := scheduler.storage.DbQueries.GetNotifyEvents(ctx, now)
	if err != nil {
		return eventModels, err
	}
	return toViewModels(events), err
}

func (scheduler *Scheduler) UpdateNotificationStatus(ctx context.Context, eventID int64) error {
	_, err := scheduler.storage.DbQueries.UpdateEventNotificationStatus(
		ctx, sqlc.UpdateEventNotificationStatusParams{
			Notificationsended: sql.NullBool{Bool: true, Valid: true},
			ID:                 eventID,
		})
	if err != nil {
		return err
	}
	return nil
}

func (scheduler *Scheduler) DeleteExpiredEvents(ctx context.Context) (err error) {
	return scheduler.storage.DbQueries.DeleteExpiredEvents(ctx)
}

func main() {
	log.Println("Starting publisher...")

	config, err := schedulerConfig.NewConfig(*configFile)
	failOnError(err, "can't read config file")

	scheduler := New(config)
	defer scheduler.amqpClient.Shutdown()

	scheduler.setupAMQP()

	context, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if err := scheduler.storage.Connect(context); err != nil {
		failOnError(err, "can't connect to database")
		cancel()
	}

	go ProccesEventNotifications(context, scheduler)

	go DeleteExpiredEvents(context, scheduler)

	<-context.Done()

	log.Println("\nFinishing publishing...")
}

func ProccesEventNotifications(context context.Context, scheduler *Scheduler) {
	for {
		events, err := scheduler.GetEventNotifications(context)
		failOnError(err, "error during getting event notifications")

		if len(events) > 0 {
			log.Printf("Got some events, count: %d\n", len(events))

			for _, event := range events {
				err := scheduler.UpdateNotificationStatus(context, event.ID)
				failOnError(err, "error during udpating event notification status")

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
		time.Sleep(time.Duration(scheduler.checkNotificationFreqMinutes) * time.Minute)
	}
}

func DeleteExpiredEvents(context context.Context, scheduler *Scheduler) {
	for {
		log.Println("Check on expired events")
		err := scheduler.DeleteExpiredEvents(context)
		failOnError(err, "error during deleting expired events")

		time.Sleep(time.Duration(scheduler.checkExpiredEventsFreqMinutes) * time.Minute)
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
