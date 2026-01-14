package main

import (
	pb "ride-sharing/shared/proto/driver"
	"ride-sharing/shared/util"
	"sync"
)

type driverInMap struct {
	Driver *pb.Driver
}
type Service struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

func NewServicce() *Service {
	return &Service{
		drivers: make([]*driverInMap, 0),
	}
}

func (s *Service) FindAavailableDrivers(packageType string) []string {
	var matchingDrivers []string

	for _, driver := range s.drivers {
		if driver.Driver.PackageSlug == packageType {
			matchingDrivers = append(matchingDrivers, driver.Driver.Id)
		}
	}
	return matchingDrivers
}

func (s *Service) RegisterDriver(driverID string, pacakgeSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := math.IntN(len(PredefinedRoutes))
	randomRoute := PredefinedRoutes[randomIndex]

	randomPlate := GenerateRandomPlate()
	randomAvatar := util.GetRandomAvatar(randomIndex)

}
