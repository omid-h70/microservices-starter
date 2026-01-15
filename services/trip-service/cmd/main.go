package main

import (
	"log"

	httpserver "ride-sharing/services/trip-service/cmd/api"
	grpcserver "ride-sharing/services/trip-service/cmd/gapi"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
)

var (
	httpAddr    = env.GetString("HTTP_ADDR", ":8081")
	grpcAddr    = env.GetString("GRPC_ADDR", ":9091")
	rabbitMQURI = env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672")
)

func main() {
	log.Printf("Starting Trip Service At http %s GRPC %s", httpAddr, grpcAddr)

	repo := repository.NewInMemRepository()
	svc := service.NewTripService(repo)

	rabbitmq, err := messaging.NewRabbitMQ(rabbitMQURI)
	if err != nil {
		log.Fatalf("rabbit is down %v", err)
	}
	defer rabbitmq.Close()
	log.Printf("Rabbitmq started on %s ", rabbitMQURI)

	//publisher := events.NewTripEventPubisher(rabbitmq)
	driverConsumer := events.NewDriverConsumer(rabbitmq, svc)
	go driverConsumer.Listen()

	errorChan := make(chan error)
	httpSever, _ := httpserver.NewHttpServer(&svc)
	httpSever.SetupRoutes()

	go func() {
		errorChan <- httpSever.RunServer(httpAddr)
	}()

	gRPCServer, _ := grpcserver.NewGRPCServer(&svc)
	gRPCServer.SetupRoutes()

	go func() {
		errorChan <- gRPCServer.RunServer(grpcAddr)
	}()

	err = <-errorChan
	if err != nil {
		log.Fatalf("sth went wrong ::: %v", err)
	}
	log.Println("Done !")
}
