package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8999")
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

func main() {
	log.Println("Starting API Gateway at " + httpAddr)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("Got %s", r.URL.Path)
		w.Write([]byte("Hello from API Gateway \n path => " + msg))
	})

	mux.HandleFunc("POST /trip/preview", handleTripPreview)
	mux.HandleFunc("/ws/driver", handleDriverWebSocket)
	mux.HandleFunc("/ws/rider", handleRidersWebSocket)

	handler := loggingMiddleware(mux)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("%v", err)
	}
}
