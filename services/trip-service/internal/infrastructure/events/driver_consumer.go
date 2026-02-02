package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	pbd "ride-sharing/shared/proto/driver"

	"github.com/rabbitmq/amqp091-go"
)

type driverConsumer struct {
	rabbitMq *messaging.RabbitMQ
	service  domain.TripService
}

func NewDriverConsumer(rabbitMq *messaging.RabbitMQ, svc domain.TripService) *driverConsumer {
	return &driverConsumer{
		rabbitMq: rabbitMq,
		service:  svc,
	}
}

func (c *driverConsumer) Listen(ctx context.Context, queueName string) error {
	return c.rabbitMq.ConsumeMessages(ctx, queueName, func(ctx context.Context, msg amqp091.Delivery) error {

		var tripEvent contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("failed to unmarshal the message %v", err)
			return err
		}

		var payload messaging.DriverTripResponseData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			log.Printf("failed to unmarshal the message payload %v", err)
			return err
		}

		switch msg.RoutingKey {
		case contracts.DriverCmdTripAccept:
			if err := c.handleTripAccepted(ctx, payload.TripID, payload.Driver); err != nil {
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

func (c *driverConsumer) handleTripAccepted(ctx context.Context, tripID string, driver *pbd.Driver) error {
	//1. validate if trip exists
	tripModel, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("1.get trip by id - failed %v", err)
	}
	//2. update the trip
	if err := c.service.UpdateTrip(ctx, tripID, "accepted", driver); err != nil {
		log.Printf("failed to update the trip %v", err)
		return fmt.Errorf("update trip failed %v", err)
	}

	//FIXME you can return it in update rather than another db fetch
	//get it again to have in updated last form
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("2.get trip by id - failed %v", err)
	}
	//3. Driver has been assigned -> publish the event
	marshalledTrip, err := json.Marshal(&tripModel)
	if err != nil {
		return fmt.Errorf("marshalling trip data failed %v", err)
	}

	err = c.rabbitMq.PublishMessage(ctx, contracts.TripEventDriverAssigned, contracts.AmqpMessage{
		OwnerID: tripModel.UserID,
		Data:    marshalledTrip,
	})
	if err != nil {
		return fmt.Errorf("publish message failed %v", err)
	}
	//TODO :: Notify payment service to do the payment

	marshalledPayLoad, err := json.Marshal(messaging.PaymentTripResponseData{
		TripID:   tripID,
		UserID:   trip.UserID,
		DriverID: driver.Id,
		Amount:   trip.RideFare.TotalPricesInCents,
		Currency: "USD",
	})
	if err != nil {
		return fmt.Errorf("failed to marshall payment trip response message %v", err)
	}

	err = c.rabbitMq.PublishMessage(ctx, contracts.TripEventDriverAssigned, contracts.AmqpMessage{
		OwnerID: tripModel.UserID,
		Data:    marshalledPayLoad,
	})
	if err != nil {
		return fmt.Errorf("publish message failed %v", err)
	}

	return nil
}

func (c *driverConsumer) handleTripDeclined(ctx context.Context, payload messaging.DriverTripResponseData) error {
	//when a driver declines, we should try to find another driver, it means we have to republish the message
	trip, err := c.service.GetTripByID(ctx, payload.TripID)
	if err != nil {
		return err
	}

	newPayload := messaging.TripEventData{
		Trip: trip.ToProto(),
	}

	marshalledPayLoad, err := json.Marshal(newPayload)
	if err != nil {
		return err
	}

	if err := c.rabbitMq.PublishMessage(ctx, contracts.TripEventDriverNotInterested, contracts.AmqpMessage{
		OwnerID: payload.RiderID,
		Data:    marshalledPayLoad,
	}); err != nil {
		return err
	}

	return nil
}
