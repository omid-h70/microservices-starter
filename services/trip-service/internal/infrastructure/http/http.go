package http

import (
	"encoding/json"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
)

type HttpHandler struct {
	service domain.TripService
}

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func NewHttpHandler(service domain.TripService) *HttpHandler {

	handler := &HttpHandler{
		service: service,
	}
	return handler
}

func (h *HttpHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "bad REST verb request", http.StatusBadRequest)
		return
	}

	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse json data", http.StatusBadRequest)
		return
	}

	//TODO add more validation
	if reqBody.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	h.service.GetRoute(ctx, &reqBody.Destination, &reqBody.Pickup)
}
