package domain

import (
	pb "ride-sharing/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                 primitive.ObjectID
	UserID             string
	PackageSlug        string
	TotalPricesInCents float64
	Route              *types.OsrmApiResponse
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                 r.ID.Hex(),
		UserId:             r.UserID,
		PackageSlug:        r.PackageSlug,
		TotalPricesInCents: r.TotalPricesInCents,
	}
}

// ToRideFaresProto - code has changed here rather to original
func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	pbRideFares := make([]*pb.RideFare, len(fares))
	for i, fare := range fares {
		pbRideFares[i] = fare.ToProto()
	}
	return pbRideFares
}
