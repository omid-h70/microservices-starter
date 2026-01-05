package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
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

	//DEPRECATED - but let it stay here for now
	if r.Method != http.MethodPost {
		http.Error(w, "bad REST verb request", http.StatusBadRequest)
		return
	}

	var tripRequestArgs previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&tripRequestArgs); err != nil {
		http.Error(w, "failed to parse json data", http.StatusBadRequest)
		return
	}

	log.Printf("Got %v", tripRequestArgs)

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

	grpcResp, err := tripService.Client.PreviewTrip(r.Context(), tripRequestArgs.toProto())
	if err != nil {
		log.Printf("failed to preview trip %v", err)
		http.Error(w, "internal server Error", http.StatusInternalServerError)
		return
	}

	resp := contracts.APIResponse{
		Data: grpcResp,
	}
	writeJSON(w, http.StatusCreated, resp)
}
