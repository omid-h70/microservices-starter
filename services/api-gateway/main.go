package main

import (
	"fmt"
	"log"
	"net/http"

	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway at " + httpAddr)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("got %s", r.URL.Path)
		w.Write([]byte("Hello from API Gateway \n path => " + msg))
	})

	mux.HandleFunc("POST /trip/preview", handleTripPreview)
	mux.HandleFunc("/ws/drivers", handleDriverWebSocket)
	mux.HandleFunc("/ws/riders", handleRidersWebSocket)

	handler := corsMiddleware(
		loggingMiddleware(
			mux,
		),
	)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("%v", err)
	}
}
