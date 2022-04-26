package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	senderConfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	amqpManager "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp"
	amqpModels "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/amqp/models"
	"github.com/streadway/amqp"
)

type AMQPManager interface {
	InitConnectionAndChannel() error
	Consume(consumerName, queueName string) (<-chan amqp.Delivery, error)
	Shutdown()
}

type Sender struct {
	AMQPManager
}

var replies <-chan amqp.Delivery

func New(c senderConfig.Config) *Sender {
	return &Sender{
		AMQPManager: &amqpManager.AMQPManager{
			AMQPURI: c.AMQP.Source,
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

func (s *Sender) ProcessReceivedMessages() {
	count := 1
	for r := range replies {
		log.Printf("Consuming reply number %d", count)
		v := amqpModels.Notification{}
		json.Unmarshal(r.Body, &v)

		fmt.Printf("IdEvent: %d,\nTitle: %s\nEventStart: %s\nIdUser: %d",
			v.IDEvent, v.EventTitle, v.EventStart.String(), v.IDUser)
		count++
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
