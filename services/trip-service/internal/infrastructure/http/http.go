package http

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/types"
	"ride-sharing/shared/util"
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
	log.Printf("Got %v", reqBody)

	//TODO add more validation
	if reqBody.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	osrmData, err := h.service.GetRoute(r.Context(), &reqBody.Destination, &reqBody.Pickup)
	if err != nil {
		http.Error(w, "GetRoute Api failed ", http.StatusBadRequest)
		return
	}

	apiResp := contracts.APIResponse{
		Data: osrmData,
	}
	util.WriteJSON(w, http.StatusCreated, apiResp)
}
