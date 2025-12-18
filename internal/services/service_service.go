package services

import (
	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/internal/repositories"
)

type ServiceService struct {
	serviceRepository *repositories.ServiceRepository
}

func NewServiceService(serviceRepository *repositories.ServiceRepository) *ServiceService {
	return &ServiceService{serviceRepository: serviceRepository}
}

func (s *ServiceService) CreateService(request dto.CreateServiceRequest) (*models.Service, error) {
	service, err := s.serviceRepository.CreateService(
		request.ProjectID,
		request.Name,
		request.Image,
		request.Environment,
		request.Mounts,
		request.Dependencies,
		request.NetworkAccess,
	)
	if err != nil {
		return nil, err
	}
	return service, nil
}
