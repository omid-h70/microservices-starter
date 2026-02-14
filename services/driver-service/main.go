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

	// Create root context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutdown signal received")
		cancel()
	}()
	//-------Create root context Done !

	service := NewServicce()

	rabbitmq, err := messaging.NewRabbitMQ(rabbitMQURI)
	if err != nil {
		log.Fatalf("rabbit is down %v", err)
	}
	defer rabbitmq.Close()
	log.Printf("Rabbitmq started on %s ", rabbitMQURI)

	errChan := make(chan error, 1)
	consumer := NewTripConsumer(rabbitmq)
	go func() {
		err = consumer.Listen(ctx, messaging.DriverCmdTripRequestQueue)
		if err != nil {
			errChan <- err
		}
	}()

	grpcServer := grpc.NewServer()
	NewGRPCHandler(grpcServer, service)
	go func() {
		err = runGRPCServer(grpcServer, grpcAddr)
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("we're done here")
	case err := <-errChan:
		log.Printf("sth went wrong %v", err)
		cancel()
	}
}

func runGRPCServer(server *grpc.Server, grpcAddr string) error {

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return err
	}

	if err := server.Serve(lis); err != nil {
		log.Printf("failed to grpc serve: %v", err)
		return err
	}

	log.Printf("grpc server started on %s", grpcAddr)
	return nil
}
