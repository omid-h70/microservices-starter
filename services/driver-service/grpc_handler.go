package main

import (
	"context"
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

// NOTE :: its dependency injection dude !
func NewGRPCHandler(grpcServer *grpc.Server, service *Service) *grpcHandler {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterDriverServiceServer(grpcServer, handler)
	return handler
}

func (h *grpcHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResposne, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterDriver not implemented")
}
func (h *grpcHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResposne, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterDriver not implemented")
}
