package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"syscall"

	"google.golang.org/grpc"
)

var (
	httpAddr    = env.GetString("HTTP_ADDR", ":8081")
	grpcAddr    = env.GetString("GRPC_ADDR", ":9091")
	rabbitMQURI = env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672")
)

func main() {
	log.Printf("Driver Service !")

	grpcServer := grpc.NewServer()
	service := NewServicce()

	NewGRPCHandler(grpcServer, service)
	printDebugInfo(grpcServer)

	rabbitmq, err := messaging.NewRabbitMQ(rabbitMQURI)
	if err != nil {
		log.Fatalf("rabbit is down %v", err)
	}
	defer rabbitmq.Close()
	log.Printf("Rabbitmq started on %s ", rabbitMQURI)

	consumer := NewTripConsumer(rabbitmq)
	go func() {
		err := consumer.Listen(context.Background(), "hello")
		if err != nil {
			log.Fatalf("failed to listen for messages")
		}
	}()
}

func printDebugInfo(grpcServer *grpc.Server) {
	for s, info := range grpcServer.GetServiceInfo() {
		log.Println("registered service:", s)
		for _, m := range info.Methods {
			log.Println("  method:", m.Name)
		}
	}
}

func runGRPCServer(server *grpc.Server, grpcAddr string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errorChan := make(chan error, 1)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//printDebugInfo()

	go func() {
		shutdownChan := make(chan os.Signal, 1)
		signal.Notify(shutdownChan, syscall.SIGTERM, os.Interrupt)
		<-shutdownChan
		cancel()
	}()

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		errorChan <- err
	}()

	select {
	case <-ctx.Done():
	case err = <-errorChan:
		log.Println("grpc server is about to end")
		server.GracefulStop()
	}

	return err
}
