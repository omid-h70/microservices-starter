package types

import (
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

type PreviewTripRequest struct {
	UserID string           `json:"userID"`
	Pickup types.Coordinate `json:"pickup"`
	Dest   types.Coordinate `json:"destination"`
}

func (r PreviewTripRequest) ToProto() *pb.PreviewTripRequest {

	return &pb.PreviewTripRequest{
		UserId: r.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  r.Pickup.Latitude,
			Longitude: r.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  r.Dest.Latitude,
			Longitude: r.Dest.Longitude,
		},
	}
}
