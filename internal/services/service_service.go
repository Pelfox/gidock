package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/internal/repositories"
	"github.com/Pelfox/gidock/internal/repositories/commands"
	"github.com/Pelfox/gidock/pkg"
	"github.com/containerd/errdefs"
	"github.com/google/uuid"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

type ServiceService struct {
	serviceRepository *repositories.ServiceRepository
	dockerClient      *client.Client
}

func NewServiceService(
	serviceRepository *repositories.ServiceRepository,
	dockerClient *client.Client,
) *ServiceService {
	return &ServiceService{
		serviceRepository: serviceRepository,
		dockerClient:      dockerClient,
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

func (s *ServiceService) pullServiceImage(ctx context.Context, service *models.Service) error {
	_, err := s.dockerClient.ImageInspect(ctx, service.Image)

	// image already exists - no need to pull
	if err == nil {
		return nil
	}
	// if error isn't "not found", then return it
	if !errdefs.IsNotFound(err) {
		return err
	}

	log.Info().Str("service_id", service.ID.String()).
		Str("image", service.Image).
		Msg("pulling image for service")

	pullOptions := client.ImagePullOptions{} // TODO: support for auth
	pullResult, err := s.dockerClient.ImagePull(ctx, service.Image, pullOptions)
	if err != nil {
		return err
	}

	defer pullResult.Close()
	if err := pullResult.Wait(ctx); err != nil {
		return fmt.Errorf("pull image failed: %w", err)
	}

	return nil
}

func (s *ServiceService) createServiceContainer(ctx context.Context, service *models.Service) (*string, error) {
	environment := make([]string, 0)
	for key, value := range service.Environment {
		environment = append(environment, fmt.Sprintf("%s=%s", key, value))
	}

	mounts := make([]mount.Mount, len(service.Mounts))
	for _, serviceMount := range service.Mounts {
		mounts = append(mounts, mount.Mount{
			Type:     mount.TypeVolume,
			Source:   serviceMount.Source,
			Target:   serviceMount.Target,
			ReadOnly: serviceMount.ReadOnly,
		})
	}

	createOptions := client.ContainerCreateOptions{
		// TODO: we should probably support more fields from this struct
		// TODO: scan `dependencies` from `service` and attach to them
		Config: &container.Config{
			Env:             environment,
			NetworkDisabled: !service.NetworkAccess,
			Labels: map[string]string{
				"gidock.service":    "true",
				"gidock.service_id": service.ID.String(),
				"gidock.project_id": service.ProjectID.String(),
			},
		},
		HostConfig: &container.HostConfig{
			Mounts: mounts,
		},
		// TODO: specify `Name`, when `models.Service` model will be updated to support it
		Image: service.Image,
	}
	createResult, err := s.dockerClient.ContainerCreate(ctx, createOptions)
	if err != nil {
		return nil, err
	}

	return &createResult.ID, nil
}

func (s *ServiceService) startServiceContainer(ctx context.Context, containerID string, service *models.Service) (*string, error) {
	startOptions := client.ContainerStartOptions{}
	_, err := s.dockerClient.ContainerStart(ctx, containerID, startOptions)

	// if container isn't found, create a new one and start it
	if errdefs.IsNotFound(err) {
		containerID, err := s.createServiceContainer(ctx, service)
		if err != nil {
			return nil, err
		}

		_, err = s.dockerClient.ContainerStart(ctx, *containerID, startOptions)
		if err != nil {
			return nil, err
		}
		return containerID, nil
	}

	return &containerID, err
}

func (s *ServiceService) StartService(ctx context.Context, id uuid.UUID, forcePull bool) (*models.Service, error) {
	service, err := s.serviceRepository.GetServiceByID(id)
	if err != nil {
		return nil, err
	}

	startContainerID := service.ContainerID

	// this service is starting for the first time, pulling its image
	if service.ContainerID == nil || forcePull {
		if err := s.pullServiceImage(ctx, service); err != nil {
			return nil, err
		}
		containerID, err := s.createServiceContainer(ctx, service)
		if err != nil {
			return nil, err
		}
		startContainerID = containerID
	}

	containerID, err := s.startServiceContainer(ctx, *startContainerID, service)
	if err != nil {
		return nil, err
	}

	updatedService, err := s.serviceRepository.UpdateServiceByID(service.ID, commands.UpdateServiceCommand{ContainerID: containerID})
	if err != nil {
		return nil, err
	}

	return updatedService, nil
}

func (s *ServiceService) stopServiceContainer(ctx context.Context, containerID string, kill bool) error {
	signal := "SIGTERM"
	if kill {
		signal = "SIGKILL"
	}

	stopOptions := client.ContainerStopOptions{
		Signal: signal,
	}
	_, err := s.dockerClient.ContainerStop(ctx, containerID, stopOptions)

	if errdefs.IsNotFound(err) {
		return nil
	}
	return err
}

func (s *ServiceService) StopService(ctx context.Context, id uuid.UUID, kill bool) error {
	service, err := s.serviceRepository.GetServiceByID(id)
	if err != nil {
		return err
	}

	if service.ContainerID == nil {
		return errors.New("service has no container attached to it")
	}

	return s.stopServiceContainer(ctx, *service.ContainerID, kill)
}

func (s *ServiceService) GetService(id uuid.UUID) (*models.Service, error) {
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
		return nil, errors.New("service has no container attached to it")
	}

	inspectResult, err := s.dockerClient.ContainerInspect(ctx, *service.ContainerID, client.ContainerInspectOptions{})
	if err != nil {
		return nil, err
	}

	return &dto.ServiceStatusResponse{
		State:     inspectResult.Container.State.Status,
		StartedAt: inspectResult.Container.State.StartedAt,
		ExitCode:  inspectResult.Container.State.ExitCode,
	}, nil
}

func (s *ServiceService) GetServiceLogs(ctx context.Context, id uuid.UUID) (<-chan pkg.LogEntry, error) {
	service, err := s.serviceRepository.GetServiceByID(id)
	if err != nil {
		return nil, err
	}

	if service.ContainerID == nil {
		return nil, errors.New("service has no container attached to it")
	}

	logsOptions := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "all",
		Details:    true,
	}
	logsResult, err := s.dockerClient.ContainerLogs(ctx, *service.ContainerID, logsOptions)
	if err != nil {
		return nil, err
	}

	channel := make(chan pkg.LogEntry)
	writer := pkg.NewLogsWriter(channel)

	go func() {
		defer logsResult.Close()
		defer close(channel)

		if _, err = stdcopy.StdCopy(writer, writer, logsResult); err != nil && !errors.Is(err, context.Canceled) {
			log.Error().Err(err).
				Str("service_id", id.String()).
				Str("container_id", *service.ContainerID).
				Msg("error copying logs")
		}
		if err := writer.FlushRemaining(); err != nil && !errors.Is(err, context.Canceled) {
			log.Error().Err(err).
				Str("service_id", id.String()).
				Str("container_id", *service.ContainerID).
				Msg("error flushing remaining logs")
		}
	}()

	return channel, nil
}
