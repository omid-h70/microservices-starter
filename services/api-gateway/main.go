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
	log.Println("Starting API Gateway")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("path %s", r.URL.Path)
		w.Write([]byte("Hello from API Gateway \n path => " + msg))
	})

	mux.HandleFunc("POST /trip/preview", handleTripPreview)
	mux.HandleFunc("/ws/driver", handleDriverWebSocket)
	mux.HandleFunc("/ws/rider", handleRidersWebSocket)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("%v", err)
	}
}
