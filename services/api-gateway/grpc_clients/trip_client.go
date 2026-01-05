package grpc_clients

import (
	"os"
	pb "ride-sharing/shared/proto/trip"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type tripServiceClient struct {
	Client pb.TripServiceClient
	conn   *grpc.ClientConn
}

var (
	tripServiceDefaultURL = "trip-service:9093"
)

func NewTripServiceClient() (*tripServiceClient, error) {

	tripServiceUrl := os.Getenv("TRIP_SERVICE_URL")
	if tripServiceUrl == "" {
		tripServiceUrl = tripServiceDefaultURL
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(tripServiceUrl, opts)
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
