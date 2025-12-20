package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/pkg"
	"github.com/containerd/errdefs"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

type DockerService struct {
	client *client.Client
}

func NewDockerService(dockerClient *client.Client) *DockerService {
	return &DockerService{client: dockerClient}
}

func (s *DockerService) PullServiceImage(ctx context.Context, service *models.Service) error {
	_, err := s.client.ImageInspect(ctx, service.Image)

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

	// TODO: support for auth
	pullResult, err := s.client.ImagePull(ctx, service.Image, client.ImagePullOptions{})
	if err != nil {
		return err
	}

	defer pullResult.Close()
	if err := pullResult.Wait(ctx); err != nil {
		return fmt.Errorf("pull image failed: %w", err)
	}

	return nil
}

func (s *DockerService) CreateServiceContainer(ctx context.Context, service *models.Service) (*string, error) {
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
	createResult, err := s.client.ContainerCreate(ctx, createOptions)
	if err != nil {
		return nil, err
	}

	return &createResult.ID, nil
}

func (s *DockerService) StartServiceContainer(
	ctx context.Context,
	containerID string,
	service *models.Service,
) (*string, error) {
	startOptions := client.ContainerStartOptions{}
	_, err := s.client.ContainerStart(ctx, containerID, startOptions)

	// if container isn't found, create a new one and start it
	if errdefs.IsNotFound(err) {
		containerID, err := s.CreateServiceContainer(ctx, service)
		if err != nil {
			return nil, err
		}

		_, err = s.client.ContainerStart(ctx, *containerID, startOptions)
		if err != nil {
			return nil, err
		}
		return containerID, nil
	}

	return &containerID, err
}

func (s *DockerService) StopContainer(ctx context.Context, containerID string, kill bool) error {
	signal := "SIGTERM"
	if kill {
		signal = "SIGKILL"
	}

	_, err := s.client.ContainerStop(ctx, containerID, client.ContainerStopOptions{
		Signal: signal,
	})

	if errdefs.IsNotFound(err) {
		return nil
	}
	return err
}

func (s *DockerService) GetContainerStatus(ctx context.Context, containerID string) (*dto.ServiceStatusResponse, error) {
	inspectResult, err := s.client.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		return nil, err
	}
	return &dto.ServiceStatusResponse{
		State:      inspectResult.Container.State.Status,
		ExitCode:   inspectResult.Container.State.ExitCode,
		StartedAt:  inspectResult.Container.State.StartedAt,
		FinishedAt: inspectResult.Container.State.FinishedAt,
	}, nil
}

func (s *DockerService) GetContainerLogs(ctx context.Context, containerID string) (<-chan pkg.LogEntry, error) {
	logsOptions := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "all",
		Details:    true,
	}
	logsResult, err := s.client.ContainerLogs(ctx, containerID, logsOptions)
	if err != nil {
		return nil, err
	}

	channel := make(chan pkg.LogEntry)
	writer := pkg.NewLogsWriter(channel)

	go func() {
		defer logsResult.Close()
		defer close(channel)

		if _, err = stdcopy.StdCopy(writer, writer, logsResult); err != nil && !errors.Is(err, context.Canceled) {
			log.Error().Err(err).Str("container_id", containerID).
				Msg("error copying logs")
		}
		if err := writer.FlushRemaining(); err != nil && !errors.Is(err, context.Canceled) {
			log.Error().Err(err).Str("container_id", containerID).
				Msg("error flushing remaining logs")
		}
	}()

	return channel, nil
}
