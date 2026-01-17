package stripe

import (
	"context"
	"fmt"
	"ride-sharing/services/payment-service/internal/domain"
	"ride-sharing/services/payment-service/pkg/types"
	//github.com/stripe/strip-go/v81
	//github.com/stripe/strip-go/v81/checkout/session
)

type stripeClient struct {
	config *types.PaymentConfig
}

func NewStripeClient(config *types.PaymentConfig) domain.PaymentProcessor {
	stripe.key = config.StripeSecretKey

	return &stripeClient{
		config: config,
	}
}

func (s *stripeClient) CreatePaymentSession(ctx context.Context, amount int64, currency string, metadata map[string]string) (string, error) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(s.config.SuccessURL),
		CancelURL:  stripe.String(s.config.CancelURL),
		MetaData:   metadata,
		LineItems: *[]strip.CheckoutSessionLineItemParams{
			{
				PriceData: &strip.CheckoutSessionLineItemPriceDataParams{
					Currency: strip.String(""),
					ProductData: &strip.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Ride Payment"),
					},
					UnitAmount: amount,
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(stripe.CheckoutSessionModePayment),
	}
	result, err := sessions.New(param)
	if err != nil {
		return "", fmt.Errorf("failed to create payment session %w", err)
	}
	return result.ID, nil
}

func (s *stripeClient) GetSessionsStatus(ctx context.Context, sessionID string) (types.PaymentStatus, error) {
	return "", nil
}
