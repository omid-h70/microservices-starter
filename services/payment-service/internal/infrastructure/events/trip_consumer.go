package events

import (
	"context"
	"encoding/json"
	"log"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/contracts"

	"github.com/rabbitmq/amqp091-go"
)

type driverConsumer struct {
	rabbitMq *messaging.RabbitMQ
	service  *service.TripService
}

func NewTripConsumer(rabbitMq *messaging.RabbitMQ) *driverConsumer {
	return &driverConsumer{
		rabbitMq: rabbitMq,
	}
}

func (c *driverConsumer) Listen(ctx context.Context, queueName string) error {
	return c.rabbitMq.ConsumeMessages(ctx, queueName, func(ctx context.Context, msg amqp091.Delivery) error {

		var tripEvent contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("failed to unmarshal the message %v", err)
			return err
		}

		var payload messaging.PaymentTripResponseData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			log.Printf("failed to unmarshal the message payload %v", err)
			return err
		}

		switch msg.RoutingKey {
		case contracts.PaymentCmdCreateSession:
			if err := c.handleTripAccepted(ctx, payload); err != nil {
				log.Printf("failed to handleTripAccepted %v", err)
				return err
			}
		case contracts.DriverCmdTripDecline:
			log.Println("trip declined")
		}

		log.Printf("unknown driver event %+v", payload)
		return nil
	})
}

func (c *driverConsumer) handleTripAccepted(ctx context.Context, payload messaging.DriverTripResponseData) error {

	log.Printf("Handling trip accepted by driver: %s", payload.TripID)

	paymentSession, err := c.service.CreatePaymentSession(
		ctx,
		payload.TripID,
		payload.UserID,
		payload.DriverID,
		int64(payload.Amount),
		payload.Currency,
	)

	if err != nil {
		log.Printf("failed to create payment session %v", err)
		return err
	}

	log.Printf("payment session created %s", paymentSession.StripeSessionID)

	//publish payment session created
	paymentPayLoad := messaging.PaymentEventSessionCreatedData{
		TripID:    payload.TripID,
		SessionID: paymentSession.StripeSessionID,
		Amount:    float(paymentSession.Amount) / 100.0, // convert cents to dollars
	}

	payloadBytes, err := json.Marshal(paymentPayLoad)
	if err != nil {
		log.Printf("failed to marshall payment session %v", err)
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.PaymentEventSessionCreated,
		contracts.AmqpMessage{
			OwnerID: payload.UserID,
			Data:    payloadBytes,
		},
	); err != nil {
		log.Printf("failed to publish payment session created event %v", err)
		return err
	}

	log.Printf("Published payment session created event %s", payload.TripID)
	return nil
}
