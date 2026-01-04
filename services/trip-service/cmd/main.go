package main

import (
	"log"

	httpserver "ride-sharing/services/trip-service/cmd/api"
	grpcserver "ride-sharing/services/trip-service/cmd/gapi"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
	grpcAddr = env.GetString("GRPC_ADDR", ":9091")
)

func main() {
	log.Println("Starting Trip Service")

	repo := repository.NewInMemRepository()
	svc := service.NewDefaultTripService(repo)

	errorChan := make(chan error)
	httpSever, _ := httpserver.NewHttpServer(&svc)
	httpSever.SetupRoutes()

	go func() {
		errorChan <- httpSever.RunServer(httpAddr)
	}()

	gRPCServer, _ := grpcserver.NewGRPCServer()
	gRPCServer.SetupRoutes()

	go func() {
		errorChan <- gRPCServer.RunServer(grpcAddr)
	}()

	err := <-errorChan
	log.Fatalf("sth went wrong ::: %v", err)
}
