package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/domain"
	ihttp "ride-sharing/services/trip-service/internal/infrastructure/http"
	"syscall"
	"time"
)

//--------------------------- HTTP Server

type HttpServer struct {
	mux *http.ServeMux
	svc domain.TripService
}

func NewHttpServer(svc domain.TripService) (HttpServer, error) {

	return HttpServer{
		mux: http.NewServeMux(),
		svc: svc,
	}, nil
}

func (s *HttpServer) SetupRoutes() {

	httpHandler := ihttp.NewHttpHandler(s.svc)

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("path %s", r.URL.Path)
		w.Write([]byte("Hello from Trip Service \n path => " + msg))
	})

	s.mux.HandleFunc("/trip/preview", httpHandler.HandleTripPreview)
}

func (s *HttpServer) RunServer(httpAddr string) error {

	var err error
	server := &http.Server{
		Addr:    httpAddr,
		Handler: s.mux,
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, os.Interrupt)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	//go func() {

	select {
	case err = <-serverErr:
		log.Printf("server got error %s", err)
	case <-shutdownChan:
		{
			log.Printf("got %v signal", shutdownChan)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Printf("could not stop the server %v", err)
				server.Close()
			}
		}
	}
	//}()

	return err
}
