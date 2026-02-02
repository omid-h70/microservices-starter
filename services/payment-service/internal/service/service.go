package service

import (
	"context"
	"fmt"
	"ride-sharing/services/payment-service/internal/domain"
	"ride-sharing/services/payment-service/pkg/types"
	"time"

	"github.com/google/uuid"
)

type paymentService struct {
	processor domain.PaymentProcessor
}

func NewPaymentService(processor domain.PaymentProcessor) domain.Service {
	return &paymentService{
		processor: processor,
	}
}

func (s *paymentService) CreatePaymentSession(
	ctx context.Context,
	tripID string,
	userID string,
	driverID string,
	amount int64,
	currency string,
) (*types.PaymentIntent, error) {
	metadata := map[string]string{
		"tripID":   tripID,
		"userID":   userID,
		"driverID": driverID,
	}

	sessionID, err := s.processor.CreatePaymentSession(
		ctx,
		amount,
		currency,
		metadata,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment session %v", err)
	}

	paymentIntent := &types.PaymentIntent{
		ID:             uuid.New().String(),
		TripID:         tripID,
		UserID:         userID,
		DriverID:       driverID,
		Amount:         amount,
		Currency:       currency,
		StripSessionID: sessionID,
		CreatedAt:      time.Now(),
	}
	return paymentIntent, nil
}
