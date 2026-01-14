package domain

import (
	"context"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	tripTypes "ride-sharing/services/trip-service/pkg/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver //FIXME
}

func (t *TripModel) ToProto() *pb.Trip {
	return &pb.Trip{
		Id:           t.ID.Hex(),
		UserId:       t.UserID,
		Status:       t.Status,
		SelectedFare: t.RideFare.ToProto(),
		Route:        t.RideFare.Route.toProto(),
		Driver:       t.Driver,
	}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, f *RideFareModel) error

	GetRideFareByID(ctx context.Context, id string) (*RideFareModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, dest *types.Coordinate) (*tripTypes.OsrmApiResponse, error)

	//EstimatePackagesPriceWithRoute taks the route and estimates the price
	EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) []*RideFareModel
	//GenerateTripFares it returns multiple packages for user to decide based on it
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userID string, route *tripTypes.OsrmApiResponse) ([]*RideFareModel, error)
	GetAndValidateFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
}
