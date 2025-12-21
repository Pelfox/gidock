package commands

import (
	"github.com/Pelfox/gidock/internal/models"
	"github.com/google/uuid"
)

// CreateServiceCommand represents the data required to create a new service.
type CreateServiceCommand struct {
	// ProjectID is the unique identifier of the project to which the service belongs.
	ProjectID uuid.UUID
	// Name is the human-readable name for the new service.
	Name string
	// Image is the container image used for the service.
	Image string
	// Environment contains environment variables for the service.
	Environment map[string]string
	// Mounts contains volume mounts for the service.
	Mounts []models.ServiceMount
	// Dependencies contains service dependencies.
	Dependencies []models.ServiceDependency
	// NetworkAccess indicates whether the service has network access.
	NetworkAccess bool
}

// GetServiceCommand represents the data required to retrieve a service.
type GetServiceCommand struct {
	// ID is the unique identifier of the service to be retrieved.
	ID uuid.UUID
}

// UpdateServiceCommand represents a partial update request for a service.
type UpdateServiceCommand struct {
	// ID is the unique identifier of the service to be updated.
	ID uuid.UUID
	// ContainerID is the new container ID for the service.
	ContainerID *string
}

// DeleteServiceCommand represents the data required to delete a service.
type DeleteServiceCommand struct {
	// ID is the unique identifier of the service to be deleted.
	ID uuid.UUID
}
