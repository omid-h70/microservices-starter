package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DefaultTripService struct {
	repo domain.TripRepository
}

func NewDefaultTripService(repo domain.TripRepository) DefaultTripService {
	return DefaultTripService{
		repo: repo,
	}
}

func (d *DefaultTripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {

	return &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}, nil
}

func (d *DefaultTripService) GetRoute(ctx context.Context, pickup, dest *types.Coordinate) (osrm *tripTypes.OsrmApiResponse, err error) {

	//original smaple
	//url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/13.388860,52.517037;13.397634,52.529407;13.428555,52.523219?overview=full&geometries=geojson")

	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude,
		pickup.Latitude,
		dest.Longitude,
		dest.Latitude,
	)

	var rsp *http.Response
	rsp, err = http.Get(url)

	defer func() {
		err = rsp.Body.Close()
	}()

	if err != nil {
		//TODO handle errors from external apis
		log.Printf("remote url request failed %v", err)
		return nil, err
	}

	var body []byte
	body, err = io.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("reading body failed %v", err)
		return nil, err
	}

	var osrmResp tripTypes.OsrmApiResponse
	err = json.Unmarshal(body, &osrmResp)
	if err != nil {
		log.Printf("json parsing failed %v", err)
		return nil, err
	}

	osrm = &osrmResp
	return
}

func estimateFareRoute(f *domain.RideFareModel, route *tripTypes.OsrmApiResponse) *domain.RideFareModel {
	pricingCnf := tripTypes.DefaultPricingConfig()
	carPackagePrice := f.TotalPricesInCents

	distannceInKM := route.Routes[0].Distance
	durationInMin := route.Routes[0].Duration

	distanceFare := distannceInKM * pricingCnf.PriceUnitOfDistance
	timeFare := durationInMin * pricingCnf.PricingPerMinute

	return &domain.RideFareModel{
		TotalPricesInCents: distanceFare + timeFare + carPackagePrice,
		PackageSlug:        f.PackageSlug,
	}
}

func getBaseFare() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:        "suv",
			TotalPricesInCents: 200,
		},
		{
			PackageSlug:        "sedan",
			TotalPricesInCents: 350,
		},
		{
			PackageSlug:        "van",
			TotalPricesInCents: 400,
		},
		{
			PackageSlug:        "luxury",
			TotalPricesInCents: 1000,
		},
	}
}

func (d *DefaultTripService) EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) []*domain.RideFareModel {

	baseFares := getBaseFare()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))

	for i, f := range baseFares {
		estimatedFares[i] = estimateFareRoute(f, route)
	}

	return nil
}

func (d *DefaultTripService) GenerateTripFares(ctx context.Context, fares []*domain.RideFareModel, userID string) ([]*domain.RideFareModel, error) {

	rideFares := make([]*domain.RideFareModel, len(fares))

	var err error
	for i, f := range fares {
		rideFares[i] = &domain.RideFareModel{
			ID:                 primitive.NewObjectID(),
			UserID:             userID,
			TotalPricesInCents: f.TotalPricesInCents,
			PackageSlug:        f.PackageSlug,
		}

		err = d.repo.SaveRideFare(ctx, rideFares[i])
		if err != nil {
			return nil, fmt.Errorf("failed to save trip fare %w", err)
		}
	}

	return rideFares, err
}
func (d *DefaultTripService) GetAndValidateFare(ctx context.Context, fareID, userID string) (*domain.RideFareModel, error) {
	return nil, nil
}
