package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	pb "ride-sharing/shared/proto/trip"
	"syscall"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedTripServiceServer
	server *grpc.Server
}

func NewGRPCServer() (GRPCServer, error) {
	return GRPCServer{
		server: grpc.NewServer(),
	}, nil
}

func (g *GRPCServer) SetupRoutes() {
	pb.RegisterTripServiceServer(g.server, g)
}

func (g *GRPCServer) RunServer(grpcAddr string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errorChan := make(chan error, 1)

	go func() {
		shutdownChan := make(chan os.Signal, 1)
		signal.Notify(shutdownChan, syscall.SIGTERM, os.Interrupt)
		<-shutdownChan
		cancel()
	}()

	go func() {

		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		if err := g.server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		errorChan <- err
	}()

	var err error
	select {
	case <-ctx.Done():
	case err = <-errorChan:
		log.Println("grpc server is about to end")
		g.server.GracefulStop()
	}

	return err
}
