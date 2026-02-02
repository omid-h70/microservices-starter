package grpc

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	service   domain.TripService
	publisher *events.TripEventPublisher
	pb.UnimplementedTripServiceServer
}

var _ pb.TripServiceServer = (*GRPCHandler)(nil)

// NOTE :: its dependency injection dude !
func NewGRPCHandler(grpcServer *grpc.Server, service domain.TripService) *GRPCHandler {

	handler := &GRPCHandler{
		//GrpcServer: grpcServer,
		service: service,
	}
	return handler
}

func (h *GRPCHandler) CreateTrip(ctx context.Context, pbReq *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {

	fareID := pbReq.GetRideFareId()
	userID := pbReq.GetUserId()

	rideFare, err := h.service.GetAndValidateFare(ctx, fareID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate fare %v", err)
	}

	trip, err := h.service.CreateTrip(ctx, rideFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trip %v", err)
	}

	if err := h.publisher.PublishTripCreated(ctx, trip); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish trip created event %v", err)
	}

	return &pb.CreateTripResponse{
		TripId: trip.ID.Hex(),
	}, nil
}

func (h *GRPCHandler) PreviewTrip(ctx context.Context, pbReq *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {

	startLocation := &types.Coordinate{
		Latitude:  pbReq.GetStartLocation().Latitude,
		Longitude: pbReq.GetStartLocation().Longitude,
	}

	endLocation := &types.Coordinate{
		Latitude:  pbReq.GetEndLocation().Latitude,
		Longitude: pbReq.GetEndLocation().Longitude,
	}

	apiResp, err := h.service.GetRoute(ctx, startLocation, endLocation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route %v", err)
	}

	//1. Estimate RideFares prices based on the route (ex: distance)
	estimatedFares := h.service.EstimatePackagesPriceWithRoute(apiResp)
	//2. Store the ride fares for the create trip to fetch and validate
	rideFares, err := h.service.GenerateTripFares(ctx, estimatedFares, pbReq.GetUserId(), apiResp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate ride fares :%v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     apiResp.ToProto(),
		RideFares: domain.ToRideFaresProto(rideFares),
	}, nil
}
