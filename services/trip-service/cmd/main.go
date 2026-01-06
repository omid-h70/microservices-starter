package main

import (
	"log"

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
	log.Printf("Starting Trip Service At http %s GRPC %s", httpAddr, grpcAddr)

	repo := repository.NewInMemRepository()
	svc := service.NewDefaultTripService(repo)

	errorChan := make(chan error)
	// httpSever, _ := httpserver.NewHttpServer(&svc)
	// httpSever.SetupRoutes()

	// go func() {
	// 	errorChan <- httpSever.RunServer(httpAddr)
	// }()

	gRPCServer, _ := grpcserver.NewGRPCServer(&svc)
	gRPCServer.SetupRoutes()

	go func() {
		errorChan <- gRPCServer.RunServer(grpcAddr)
	}()

	err := <-errorChan
	if err != nil {
		log.Fatalf("sth went wrong ::: %v", err)
	}
	log.Println("Done !")
}
