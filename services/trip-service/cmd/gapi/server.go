package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/domain"
	hgrpc "ride-sharing/services/trip-service/internal/infrastructure/grpc"
	pb "ride-sharing/shared/proto/trip"
	"syscall"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	handler *hgrpc.GRPCHandler
	server  *grpc.Server
}

func NewGRPCServer(service domain.TripService) (GRPCServer, error) {

	grpcServer := grpc.NewServer()

	return GRPCServer{
		server:  grpcServer,
		handler: hgrpc.NewGRPCHandler(grpcServer, service),
	}, nil
}

func (grpcServer *GRPCServer) SetupRoutes() {
	pb.RegisterTripServiceServer(grpcServer.server, grpcServer.handler)
}

func (grpcServer *GRPCServer) printDebugInfo() {
	for s, info := range grpcServer.server.GetServiceInfo() {
		log.Println("registered service:", s)
		for _, m := range info.Methods {
			log.Println("  method:", m.Name)
		}
	}
}

func (grpcServer *GRPCServer) RunServer(grpcAddr string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errorChan := make(chan error, 1)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer.printDebugInfo()

	go func() {
		shutdownChan := make(chan os.Signal, 1)
		signal.Notify(shutdownChan, syscall.SIGTERM, os.Interrupt)
		<-shutdownChan
		cancel()
	}()

	go func() {
		if err := grpcServer.server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		errorChan <- err
	}()

	select {
	case <-ctx.Done():
	case err = <-errorChan:
		log.Println("grpc server is about to end")
		grpcServer.server.GracefulStop()
	}

	return err
}
