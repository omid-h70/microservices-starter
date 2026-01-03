package grpc_clients

import (
	"os"
	pb "ride-sharing/shared/proto/trip"

	"google.golang.org/grpc"
)

type tripServiceClient struct {
	Client pb.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient() (*tripServiceClient, error) {

	tripServiceUrl := os.Getenv("TRIP_SERVICE_URL")
	if tripServiceUrl == "" {
		tripServiceUrl = "trip-service:9093"
	}

	conn, err := grpc.NewClient(tripServiceUrl)
	if err != nil {
		return nil, err
	}

	pbTripServiceClient := pb.NewTripServiceClient(conn)

	return &tripServiceClient{
		Client: pbTripServiceClient,
		conn:   conn,
	}, nil

}

func (t *tripServiceClient) Close() {
	t.conn.Close()
}
