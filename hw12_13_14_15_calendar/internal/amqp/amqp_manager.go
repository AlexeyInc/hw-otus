package amqpmanager

import (
	"time"

	"github.com/streadway/amqp"
)

type AMQPManager struct {
	AmqpURI string
	conn    *amqp.Connection
	ch      *amqp.Channel
}

func (amqpManager *AMQPManager) InitConnectionAndChannel() error {
	conn, err := amqp.Dial(amqpManager.AmqpURI)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	amqpManager.conn = conn
	amqpManager.ch = ch

	return nil
}

func (amqpManager *AMQPManager) Publish(payload []byte, exchangeName, routingKey string) error {
	err := amqpManager.ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "application/json",
			Body:         payload,
			Timestamp:    time.Now(),
		})
	if err != nil {
		return err
	}
	return nil
}

func (amqpManager *AMQPManager) Consume(consumerName, queueName string) (<-chan amqp.Delivery, error) {
	replies, err := amqpManager.ch.Consume(
		queueName,
		consumerName,
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return replies, nil
}

func (amqpManager *AMQPManager) DeclareExchange(exchangeName, exchangeKind string) error {
	err := amqpManager.ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (amqpManager *AMQPManager) DeclareQueue(queueName string) error {
	_, err := amqpManager.ch.QueueDeclare(
		queueName,
		true, // durable
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}
	return nil
}

func (amqpManager AMQPManager) BindQueue(exchangeName, queueName, bindKey string) error {
	err := amqpManager.ch.QueueBind(
		queueName,
		bindKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (amqpManager *AMQPManager) Shutdown() {
	amqpManager.ch.Close()
	amqpManager.conn.Close()
}
