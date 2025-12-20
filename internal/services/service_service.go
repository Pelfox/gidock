package services

import (
	"context"

	"github.com/Pelfox/gidock/internal"
	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/internal/repositories"
	"github.com/Pelfox/gidock/internal/repositories/commands"
	"github.com/Pelfox/gidock/pkg"
	"github.com/google/uuid"
)

type ServiceService struct {
	serviceRepository *repositories.ServiceRepository
	dockerService     *DockerService
}

func NewServiceService(
	serviceRepository *repositories.ServiceRepository,
	dockerService *DockerService,
) *ServiceService {
	return &ServiceService{
		serviceRepository: serviceRepository,
		dockerService:     dockerService,
	}
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

func (s *ServiceService) StartService(ctx context.Context, id uuid.UUID, forcePull bool) (*models.Service, error) {
	service, err := s.serviceRepository.GetServiceByID(id)
	if err != nil {
		return nil, err
	}

	startContainerID := service.ContainerID

	// this service is starting for the first time, pulling its image
	if service.ContainerID == nil || forcePull {
		if err := s.dockerService.PullServiceImage(ctx, service); err != nil {
			return nil, err
		}
		containerID, err := s.dockerService.CreateServiceContainer(ctx, service)
		if err != nil {
			return nil, err
		}
		startContainerID = containerID
	}

	containerID, err := s.dockerService.StartServiceContainer(ctx, *startContainerID, service)
	if err != nil {
		return nil, err
	}

	updatedService, err := s.serviceRepository.UpdateServiceByID(service.ID, commands.UpdateServiceCommand{ContainerID: containerID})
	if err != nil {
		return nil, err
	}

	return updatedService, nil
}

func (s *ServiceService) StopService(ctx context.Context, id uuid.UUID, kill bool) error {
	service, err := s.serviceRepository.GetServiceByID(id)
	if err != nil {
		return err
	}
	if service.ContainerID == nil {
		return internal.ErrNoContainer
	}
	return s.dockerService.StopContainer(ctx, *service.ContainerID, kill)
}

func (s *ServiceService) GetServiceByID(id uuid.UUID) (*models.Service, error) {
	service, err := s.serviceRepository.GetServiceByID(id)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceService) ListServices() ([]models.Service, error) {
	services, err := s.serviceRepository.ListServices()
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (s *ServiceService) GetServiceStatus(ctx context.Context, id uuid.UUID) (*dto.ServiceStatusResponse, error) {
	service, err := s.serviceRepository.GetServiceByID(id)
	if err != nil {
		return nil, err
	}
	if service.ContainerID == nil {
		return nil, internal.ErrNoContainer
	}
	return s.dockerService.GetContainerStatus(ctx, *service.ContainerID)
}

func (s *ServiceService) GetServiceLogs(ctx context.Context, id uuid.UUID) (<-chan pkg.LogEntry, error) {
	service, err := s.serviceRepository.GetServiceByID(id)
	if err != nil {
		return nil, err
	}
	if service.ContainerID == nil {
		return nil, internal.ErrNoContainer
	}
	return s.dockerService.GetContainerLogs(ctx, *service.ContainerID)
}
