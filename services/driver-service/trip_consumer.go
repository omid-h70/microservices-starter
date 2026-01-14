package main

import (
	"context"
	"encoding/json"
	"log"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type tripConsumer struct {
	rabbitMq *messaging.RabbitMQ
	service  *Service
}

func NewTripConsumer(rabbitMq *messaging.RabbitMQ) *tripConsumer {
	return &tripConsumer{
		rabbitMq: rabbitMq,
	}
}

func (c *tripConsumer) Listen(ctx context.Context, queueName string) error {
	return c.rabbitMq.ConsumeMessages(ctx, queueName, func(ctx context.Context, msg amqp091.Delivery) error {

		var tripEvent contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("failed to unmarshal the message %v", err)
			return err
		}

		var payload messaging.TripEventData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			log.Printf("failed to unmarshal the message payload %v", err)
			return err
		}

		switch msg.RoutingKey {
		case contracts.TripEventCreated, contracts.TripEventDriverNotInterested:
			return c.handleFindAndNotifyDriver(ctx, payload)
		}

		log.Printf("unknown trip event %+v", payload)
		return nil
	})
}

func (c *tripConsumer) handleFindAndNotifyDriver(ctx context.Context, payload messaging.TripEventData) error {
	suitableIDs := c.service.FindAavailableDrivers(payload.Trip.SelectedFare.PackageSlug)

	if len(suitableIDs) == 0 {
		//Failed Event
		if err := c.rabbitMq.PublishMessage(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: payload.Trip.Id,
		}); err != nil {
			log.Printf("failed to publish message to exchange %v", err)
			return err
		}
	}

	suitableDriverID := suitableIDs[0]
	log.Printf("found suitable driver %v", suitableDriverID)

	marshalledEvent, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	//Notify the driver for potential trip
	if err := c.rabbitMq.PublishMessage(ctx, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
		OwnerID: suitableDriverID,
		Data:    marshalledEvent,
	}); err != nil {
		log.Printf("failed to publish message to exchange %v", err)
		return err
	}

	return nil
}
