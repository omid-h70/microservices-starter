package main

import (
	"fmt"
	"log"
	"net/http"

	myhttp "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting Trip Service")

	mux := http.NewServeMux()
	repo := repository.NewDefaultTripRepository()
	svc := service.NewDefaultTripService(&repo)
	h := myhttp.HttpHandler{Service: &svc}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("path %s", r.URL.Path)
		w.Write([]byte("Hello from API Gateway \n path => " + msg))
	})

	mux.HandleFunc("/trip/preview", h.HandleTripPreview)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("%v", err)
	}
}
