package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
)

type DefaultTripRepository struct {
}

func NewDefaultTripRepository() DefaultTripRepository {
	return DefaultTripRepository{}
}

func (d *DefaultTripRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	return &domain.TripModel{}, nil
}
