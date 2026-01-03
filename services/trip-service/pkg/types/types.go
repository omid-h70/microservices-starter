package types

import (
	pb "ride-sharing/shared/proto/trip"
)

type OsrmApiResponse struct {
	Routes []struct {
		Distance float64 `josn:"distance"`
		Duration float64 `josn:"duration"`
		Geometry struct {
			Coordinate [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (o *OsrmApiResponse) ToProto() *pb.Route {

	//NOTE
	//why just getting first route ?
	route := o.Routes[0]
	geometry := route.Geometry.Coordinate
	coordinates := make([]*pb.Coordinate, len(geometry))

	for i, coordinate := range route.Geometry.Coordinate {
		coordinates[i] = &pb.Coordinate{
			Latitude:  coordinate[0],
			Longitude: coordinate[1],
		}
	}

	return &pb.Route{
		Geometry: []*pb.Geometry{
			// Initialize the first Item struct
			{
				Coordinates: coordinates,
			},
		},
		Distance: route.Distance,
		Duration: route.Duration,
	}
}
