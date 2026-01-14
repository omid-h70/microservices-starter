package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/services/api-gateway/types"
	"ride-sharing/services/api-gateway/utils"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
)

type HttpApi struct {
	mux      *http.ServeMux
	handler  http.Handler
	rabbitmq *messaging.RabbitMQ
}

func NewHttpApiServer(rabbit *messaging.RabbitMQ) *HttpApi {

	mux := http.NewServeMux()
	return &HttpApi{
		rabbitmq: rabbit,
		mux:      mux,
	}
}

func (api *HttpApi) AddRoutes() {
	api.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("got %s", r.URL.Path)
		w.Write([]byte("Hello from API Gateway \n path => " + msg))
	})

	api.mux.HandleFunc("POST /trip/preview", handleTripPreview)
	api.mux.HandleFunc("/ws/drivers", api.handleDriverWebSocket)
	api.mux.HandleFunc("/ws/riders", api.handleRidersWebSocket)
}

func (api *HttpApi) RunServer(httpAddr string) error {

	server := &http.Server{
		Addr:    httpAddr,
		Handler: api.handler,
	}
	return server.ListenAndServe()
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {

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
