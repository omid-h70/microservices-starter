package grpc_clients

import (
	"os"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type driverServiceClient struct {
	Client pb.DriverServiceClient
	conn   *grpc.ClientConn
}

var (
	driverServiceDefaultURL = "driver-service:9093"
)

func NewDriverServiceClient() (*driverServiceClient, error) {

	tripServiceUrl := os.Getenv("DRIVER_SERVICE_URL")
	if tripServiceUrl == "" {
		tripServiceUrl = tripServiceDefaultURL
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(tripServiceUrl, opts)
	if err != nil {
		return nil, err
	}

	pbTripServiceClient := pb.NewDriverServiceClient(conn)

	return &driverServiceClient{
		Client: pbTripServiceClient,
		conn:   conn,
	}, nil

}

func (t *driverServiceClient) Close() {
	t.conn.Close()
}
