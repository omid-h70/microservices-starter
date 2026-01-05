package grpc

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	Service domain.TripService

	pb.UnimplementedTripServiceServer
	grpcServer *grpc.Server
}

// NOTE :: its dependency injection dude !
func NewGRPCHandler(grpcServer *grpc.Server, service domain.TripService) *GRPCHandler {

	handler := &GRPCHandler{
		grpcServer: grpcServer,
		Service:    service,
	}
	pb.RegisterTripServiceServer(grpcServer, handler)
	return handler
}

func (h *GRPCHandler) CreateTrip(ctx context.Context, pbReq *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	return nil, nil
}

func (h *GRPCHandler) PreviewTrip(ctx context.Context, pbReq *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {

	//ctx := r.Context()
	startLocation := &types.Coordinate{
		Latitude:  pbReq.GetStartLocation().Latitude,
		Longitude: pbReq.GetStartLocation().Longitude,
	}

	endLocation := &types.Coordinate{
		Latitude:  pbReq.GetEndLocation().Latitude,
		Longitude: pbReq.GetEndLocation().Longitude,
	}

	apiResp, err := h.Service.GetRoute(ctx, startLocation, endLocation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     apiResp.ToProto(),
		RideFares: []*pb.RideFare{},
	}, nil
}
