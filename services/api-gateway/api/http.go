package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/services/api-gateway/types"
	"ride-sharing/services/api-gateway/utils"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/env"

	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	//github.com/stripe/stripe-go/v81/webhook
)

func (api *HttpApi) handleTripPreview(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	//DEPRECATED - but let it stay here for now
	if r.Method != http.MethodPost {
		http.Error(w, "bad REST verb request", http.StatusBadRequest)
		return
	}

	var tripRequestArgs types.PreviewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&tripRequestArgs); err != nil {
		http.Error(w, "failed to parse json data", http.StatusBadRequest)
		return
	}

	log.Printf("got after decode %v", tripRequestArgs)

	//TODO add more validation
	if tripRequestArgs.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	//TODO - make connections better and more powerful !
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Printf("failed to stablish grpc connection %v", err)
		http.Error(w, "internal server.Error", http.StatusInternalServerError)
		return
	}
	defer tripService.Close()

	grpcResp, err := tripService.Client.PreviewTrip(r.Context(), tripRequestArgs.ToProto())
	if err != nil {
		log.Printf("failed to preview trip %v", err)
		http.Error(w, "internal server Error", http.StatusInternalServerError)
		return
	}

	resp := contracts.APIResponse{
		Data: grpcResp,
	}
	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (api *HttpApi) handleStripeWebhook(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read ", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	webhookKey := env.GetString("STRIPE_WEBHOOK_KEY", "")
	if webhookKey == "" {
		log.Printf("webhookkey required")
		return
	}

	event, err := webhook.ConstructEventWithOptions(
		body,
		r.Header.Get("Strip-Signature"),
		webhookKey,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)

	if err != nil {
		log.Printf("error verifying webhook signature %v", err)
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	log.Printf("received stripe event: %v", event)

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession

		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("error pasing webhook event data %v", err)
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}
	}

	payload := messaging.PaymentStatusUpdateData{
		TripID:   session.Metadata["trip_id"],
		UserID:   session.Metadata["user_id"],
		DriverID: session.Metadata["driver_id"],
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error pasing webhook event data %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	message := contracts.AmqpMessage{
		OwnerID: session.MetaData["user_id"],
		Data:    payloadBytes,
	}

	if err := rb.PublishMessage(
		r.Context(),
		contracts.PaymentEventSuccess,
		message,
	); err != nil {
		log.Printf("error publishing payment event %v", err)
		http.Error(w, "failed to publish payment event", http.StatusBadRequest)
		return
	}
}
