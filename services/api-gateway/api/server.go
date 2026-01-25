package api

import (
	"fmt"
	"net/http"
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
	api.mux.HandleFunc("/webhook/stripe", api.handleStripeWebhook)

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
