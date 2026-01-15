package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	pbd "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
)

var _ domain.TripRepository = (*InMemRepository)(nil)

type InMemRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (inmem *InMemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	inmem.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (inmem *InMemRepository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	inmem.rideFares[f.ID.Hex()] = f
	return nil
}

func (inmem *InMemRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	trip, ok := inmem.rideFares[id]
	if !ok {
		return nil, fmt.Errorf("ridefare not found with id %s", id)
	}
	return trip, nil
}

func (inmem *InMemRepository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	trip, ok := inmem.trips[id]
	if !ok {
		return nil, fmt.Errorf("trip not found with id %s", id)
	}
	return trip, nil
}

func (inmem *InMemRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error {
	trip, ok := inmem.trips[tripID]
	if !ok {
		return fmt.Errorf("trip not found with id %s", tripID)
	}

	trip.Status = status

	if driver != nil {
		trip.Driver = &pb.TripDriver{
			Id:         driver.Id,
			Name:       driver.Name,
			CarPlate:   driver.CarPlate,
			ProfilePic: driver.ProfilePic,
		}
	}
	return nil
}
