package dto

import (
	"github.com/Pelfox/gidock/internal/models"
	"github.com/google/uuid"
	"github.com/moby/moby/api/types/container"
)

// CreateServiceRequest is the request payload for creating a new service.
type CreateServiceRequest struct {
	// ProjectID identifies the parent project.
	ProjectID uuid.UUID `json:"project_id"`
	// Name is a name for the service.
	Name string `json:"name"`
	// Image is the Docker image and tag to deploy.
	Image string `json:"image"`
	// Environment is a map of environment variables passed to the container.
	Environment map[string]string `json:"environment"`
	// Mounts defines volume and bind mounts for the container.
	Mounts []models.ServiceMount `json:"mounts"`
	// Dependencies lists other services that must be running before this one starts.
	Dependencies []models.ServiceDependency `json:"dependencies"`
	// NetworkAccess indicates whether the service should be exposed externally.
	NetworkAccess bool `json:"network_access"`
}

// CreateServiceResponse is the response payload after successfully creating a service.
type CreateServiceResponse struct {
	models.Service
}

// ServiceStatusResponse provides runtime status information about a deployed
// service at the specific point of time.
type ServiceStatusResponse struct {
	// State is the current container state.
	State container.ContainerState `json:"state"`
	// StartedAt is the timestamp when the container started (if running).
	StartedAt string `json:"started_at"`

	// TODO: add `exited_at`

	// ExitCode is the exit code if the container has stopped.
	ExitCode int `json:"exit_code"`
}
