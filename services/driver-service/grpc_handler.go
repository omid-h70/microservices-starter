package main

import (
	"context"
	"log"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pb.UnimplementedDriverServiceServer
	service *Service
}

var _ pb.DriverServiceServer = (*grpcHandler)(nil)

func printDebugInfo(grpcServer *grpc.Server) {
	for s, info := range grpcServer.GetServiceInfo() {
		log.Println("registered service:", s)
		for _, m := range info.Methods {
			log.Println("  method:", m.Name)
		}
	}
}

// NOTE :: its dependency injection dude !
func NewGRPCHandler(grpcServer *grpc.Server, service *Service) *grpcHandler {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterDriverServiceServer(grpcServer, handler)
	printDebugInfo(grpcServer)
	return handler
}

func (h *grpcHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResposne, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterDriver not implemented")
}
func (h *grpcHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResposne, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterDriver not implemented")
}
