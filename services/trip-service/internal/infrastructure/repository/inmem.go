package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
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
	return nil, nil
}
