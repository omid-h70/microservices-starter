package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/shared/contracts"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TripExchangeName = "trip"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

var (
	queueMap map[string][]string = map[string][]string{
		FindAvailableDriversQueue: {
			contracts.TripEventCreated,
			contracts.TripEventDriverNotInterested,
		},
		DriverCmdTripRequestQueue: {
			contracts.DriverCmdTripRequest,
		},
		DriverTripResponseQueue: {
			contracts.DriverCmdTripAccept,
			contracts.DriverCmdTripDecline,
		},
		NotifyDriverNoDriverFoundQueue: {
			contracts.TripEventNoDriversFound,
		},
		NotifyDriverAssingedQueue: {
			contracts.TripEventDriverAssigned,
		},
	}
)

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

	if err := rmq.setupExchangesAndQueues(TripExchangeName); err != nil {

		rmq.Close()
		return nil, fmt.Errorf("faild to setup exchanges or queues %w", err)
	}

	return rmq, nil
}

type MessageHandler func(context.Context, amqp.Delivery) error

func (r *RabbitMQ) ConsumeMessages(ctx context.Context, quequeName string, handler MessageHandler) error {

	//set prefetch count to 1 for fair dispatch
	//this tells rabbitmq to not to give more than one message to a service  at a time
	//the worker will only get the next message after it has acknowledeged the prev one

	err := r.Channel.Qos(
		1,     //prefetch count : limit to 1 unacknowledged per consumer
		0,     //prefetch size : no specific limit on message size
		false, //global: apply prefetch count to eache consumer individually
	)

	if err != nil {
		return fmt.Errorf("set Qos failed %v", err)
	}

	deliveryChan, err := r.Channel.Consume(
		quequeName,
		"",           //consumer
		true,         //auto-ack
		false,        //exclusive
		false,        //no local
		false,        //no wait
		amqp.Table{}, //
	)

	if err != nil {
		return fmt.Errorf("consume messages failed %v", err)
	}

	go func() {
		for msg := range deliveryChan {
			err := handler(ctx, msg)
			if err != nil {
				log.Printf("ERROR ::: failed to process msg %v, Message Body %s", err, msg.Body)

				//TODO:: consider a dead-letter or other retry mechanism
				if err := msg.Nack(false, true); err != nil {
					log.Printf("ERROR ::: failed to nack ! %v", err)
				}
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Printf("ERROR ::: failed to ack ! %v Message Body %s", err, msg.Body)
			}
		}
	}()
	return nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshaling queue messaged failed ,err %w", err)
	}

	err = r.Channel.PublishWithContext(
		ctx,
		"",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         jsonData,
			DeliveryMode: amqp.Persistent,
		},
	)
	return err
}
func (r *RabbitMQ) setupExchangesAndQueues(exchangeName string) error {

	err := r.Channel.ExchangeDeclare(
		exchangeName,
		"topic", //direct, topic, fanout, ....
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declaring exchange faield %s, err %v", exchangeName, err)
	}

	for qName, qEvents := range queueMap {

		err = r.declareAndBindQueue(
			exchangeName,
			qName,
			qEvents,
		)
		if err != nil {
			return fmt.Errorf("failed to create queue %s - %w", qName, err)
		}
	}

	return nil
}

func (r *RabbitMQ) declareAndBindQueue(exchangeName, queueName string, routes []string) error {

	q, err := r.Channel.QueueDeclare(
		queueName,
		true, //durable
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue %w", err)
	}

	for _, route := range routes {

		err = r.Channel.QueueBind(
			q.Name,
			route,
			exchangeName,
			false,
			nil,
		)

		if err != nil {
			return fmt.Errorf("failed to declare queue %w", err)
		}
	}
	return nil
}

func (rabbitMQ *RabbitMQ) Close() {
	if rabbitMQ.Channel != nil {
		rabbitMQ.Channel.Close()
	}

	if rabbitMQ.conn != nil {
		rabbitMQ.conn.Close()
	}
}
