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

// TODO: add other methods (from Repository)

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

func (s *ServiceService) Create(
	ctx context.Context,
	request dto.CreateServiceRequest,
) (*models.Service, error) {
	return s.serviceRepository.Create(ctx, commands.CreateServiceCommand{
		ProjectID:     request.ProjectID,
		Name:          request.Name,
		Image:         request.Image,
		Environment:   request.Environment,
		Mounts:        request.Mounts,
		Dependencies:  request.Dependencies,
		NetworkAccess: request.NetworkAccess,
	})
}

func (s *ServiceService) GetByID(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	return s.serviceRepository.Get(ctx, commands.GetServiceCommand{ID: id})
}

func (s *ServiceService) ListAll(ctx context.Context) ([]models.Service, error) {
	return s.serviceRepository.ListAll(ctx)
}

func (s *ServiceService) Start(ctx context.Context, id uuid.UUID, forcePull bool) (*models.Service, error) {
	service, err := s.serviceRepository.Get(ctx, commands.GetServiceCommand{ID: id})
	if err != nil {
		return nil, err
	}

	// TODO: implement transaction boundary
	var containerID *string

	// create a new container if this is the first start or if forcePull is enabled
	if service.ContainerID == nil || forcePull {
		if err := s.dockerService.PullServiceImage(ctx, service); err != nil {
			return nil, err
		}
		containerID, err = s.dockerService.CreateServiceContainer(ctx, service)
		if err != nil {
			return nil, err
		}
	} else {
		containerID = service.ContainerID
	}

	startedContainerID, err := s.dockerService.StartServiceContainer(ctx, *containerID, service)
	if err != nil {
		return nil, err
	}

	updatedService, err := s.serviceRepository.Update(
		ctx,
		commands.UpdateServiceCommand{
			ID:          id,
			ContainerID: startedContainerID,
		},
	)
	if err != nil {
		return nil, err
	}

	return updatedService, nil
}

func (s *ServiceService) Stop(ctx context.Context, id uuid.UUID, kill bool) error {
	service, err := s.serviceRepository.Get(ctx, commands.GetServiceCommand{ID: id})
	if err != nil {
		return err
	}
	if service.ContainerID == nil {
		return internal.ErrNoContainer
	}
	// TODO: implement transaction boundary
	return s.dockerService.StopContainer(ctx, *service.ContainerID, kill)
}

func (s *ServiceService) GetStatus(ctx context.Context, id uuid.UUID) (*dto.ServiceStatusResponse, error) {
	service, err := s.serviceRepository.Get(ctx, commands.GetServiceCommand{ID: id})
	if err != nil {
		return nil, err
	}
	if service.ContainerID == nil {
		return nil, internal.ErrNoContainer
	}
	return s.dockerService.GetContainerStatus(ctx, *service.ContainerID)
}

func (s *ServiceService) StreamLogs(ctx context.Context, id uuid.UUID) (<-chan pkg.LogEntry, error) {
	service, err := s.serviceRepository.Get(ctx, commands.GetServiceCommand{ID: id})
	if err != nil {
		return nil, err
	}
	if service.ContainerID == nil {
		return nil, internal.ErrNoContainer
	}
	return s.dockerService.GetContainerLogs(ctx, *service.ContainerID)
}
