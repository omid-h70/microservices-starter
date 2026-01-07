package messaging

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("faild to connect to rabbitmq %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("faild to create rabbitmq channel %w", err)
	}

	rmq := &RabbitMQ{
		conn:    conn,
		Channel: channel,
	}

	if err := rmq.setupExchangesAndQueues(); err != nil {

		rmq.Close()
		return nil, fmt.Errorf("faild to setup exchanges or queues %w", err)
	}

	return rmq, nil
}

func (amap *RabbitMQ) setupExchangesAndQueues() error {
	return nil
}

func (amqp *RabbitMQ) Close() {
	if amqp.Channel != nil {
		amqp.Channel.Close()
	}

	if amqp.conn != nil {
		amqp.conn.Close()
	}
}
