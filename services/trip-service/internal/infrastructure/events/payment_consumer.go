package events

import (
	"context"
	"encoding/json"
	"log"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type paymentConsumer struct {
	rabbitMq *messaging.RabbitMQ
	service  *service.TripService
}

func NewPaymentConsumer(rabbitMq *messaging.RabbitMQ) *paymentConsumer {
	return &paymentConsumer{
		rabbitMq: rabbitMq,
	}
}

func (c *paymentConsumer) Listen(ctx context.Context, queueName string) error {
	//NotifyPaymentSuccessQueue
	return c.rabbitMq.ConsumeMessages(ctx, queueName, func(ctx context.Context, msg amqp091.Delivery) error {

		var tripEvent contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("failed to unmarshal the message %v", err)
			return err
		}

		var payload messaging.PaymentStatusUpdateData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			log.Printf("failed to unmarshal the message payload %v", err)
			return err
		}

		log.Println("trip has been compelted and paid !")

		return c.service.UpdateTrip(
			ctx,
			payload.TripID,
			"payed",
			nil,
		)
	})
}
