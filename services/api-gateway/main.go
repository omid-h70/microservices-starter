package main

import (
	"log"

	"ride-sharing/services/api-gateway/api"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
)

var (
	httpAddr    = env.GetString("HTTP_ADDR", ":8081")
	rabbitMQURI = env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672")
)

func main() {
	log.Println("Starting API Gateway at " + httpAddr)

	rabbitmq, err := messaging.NewRabbitMQ(rabbitMQURI)
	if err != nil {
		log.Fatalf("rabbit is down %v", err)
	}
	defer rabbitmq.Close()
	log.Printf("Rabbitmq started on %s ", rabbitMQURI)

	httpServer := api.NewHttpApiServer(rabbitmq)
	httpServer.AddRoutes()
	httpServer.AddMiddleWares()
	if err := httpServer.RunServer(httpAddr); err != nil {
		log.Fatalf("http server failed %v", err)
	}
}
