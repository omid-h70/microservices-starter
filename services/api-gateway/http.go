package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/proto/trip"
)

func _handleTripPreview(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonData)

	resp, err := http.Post("http://trip-service:8083/preview", "appplication/json", reader)
	if err != nil {
		http.Error(w, "internal server.Error", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	//TODO call external trip service

	apiResp := contracts.APIResponse{
		Data: "ok",
	}
	writeJSON(w, http.StatusCreated, apiResp)
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	//DEPRECATED
	if r.Method != http.MethodPost {
		http.Error(w, "bad REST verb request", http.StatusBadRequest)
		return
	}

	var tripRequestArgs previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&tripRequestArgs); err != nil {
		http.Error(w, "failed to parse json data", http.StatusBadRequest)
		return
	}

	//TODO add more validation
	if tripRequestArgs.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	//TODO - make connections better and more powerful !
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		http.Error(w, "internal server.Error", http.StatusInternalServerError)
		return
	}
	defer tripService.Close()

	grpcReqArgs := trip.PreviewTripRequest{
		UserId: tripRequestArgs.UserID,
		StartLocation: &trip.Coordinate{
			Latitude:  tripRequestArgs.Pickup.Latitude,
			Longitude: tripRequestArgs.Pickup.Longitude,
		},
		EndLocation: &trip.Coordinate{
			Latitude:  tripRequestArgs.Dest.Latitude,
			Longitude: tripRequestArgs.Dest.Longitude,
		},
	}
	grpcResp, err := tripService.Client.PreviewTrip(r.Context(), &grpcReqArgs)
	if err != nil {
		http.Error(w, "internal server.Error", http.StatusInternalServerError)
		return
	}

	//jsonData, err := json.Marshal(reqBody)

	//TODO call external trip service

	resp := contracts.APIResponse{
		Data: grpcResp,
	}
	writeJSON(w, http.StatusCreated, resp)
}
