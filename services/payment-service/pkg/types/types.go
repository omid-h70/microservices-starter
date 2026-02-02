package types

import (
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucces    PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type Payment struct {
	ID               string    `json:"id"`
	TripID           string    `json:"trip_id"`
	UserID           string    `json:"user_id"`
	DriverID         string    `json:"driver_id"`
	Amount           int64     `json:"amount"`
	Currency         string    `json:"currency"`
	Status           string    `json:"status"`
	StripeSessinonID string    `json:"strip_session_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PaymentIntent struct {
	ID             string    `json:"id"`
	TripID         string    `json:"trip_id"`
	UserID         string    `json:"user_id"`
	DriverID       string    `json:"driver_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	StripSessionID string    `json:"strip_session_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type PaymentConfig struct {
	StripeSecretKey      string `json:"strip_secret_key"`
	StripePublishableKey string `json:"strip_publishable_key"`
	StripeWebhookSecret  string `json:"strip_webhook_secret"`
	SuccessURL           string `json:"sucess_url"`
	CancelURL            string `json:"cancel_url"`
}
