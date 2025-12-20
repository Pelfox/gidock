package models

import (
	"time"

	"github.com/google/uuid"
)

// ServiceMount defines a single volume or bind mount for a service container.
type ServiceMount struct {
	// TODO: specify whether it is a volume or a host mount

	// Source is the host path, volume name, or named volume source.
	Source string `json:"source"`
	// Target is the path inside the container where the source is mounted.
	Target string `json:"target"`
	// ReadOnly indicates whether the mount should be read-only inside the container.
	ReadOnly bool `json:"read_only"`
}

// ServiceCondition represents the condition a dependency must satisfy.
type ServiceCondition string

const (
	// ServiceConditionHealthy means the dependency service must be healthy
	// (e.g. passing health checks).
	ServiceConditionHealthy ServiceCondition = "healthy"
	// ServiceConditionReady means the dependency service must be running and
	// ready (e.g., started successfully, but not necessarily passing health
	// checks).
	ServiceConditionReady ServiceCondition = "ready"
)

// ServiceDependency expresses that one service depends on another service
// being in a specific state before it can start.
type ServiceDependency struct {
	// ServiceID is the unique ID of the service this one depends on.
	ServiceID uuid.UUID `json:"service_id"`
	// Condition specifies the required state of the dependency service.
	Condition ServiceCondition `json:"condition"`
}

// Service represents a deployable service (internally a container) within a Project.
type Service struct {
	// ID is the unique identifier of the service.
	ID uuid.UUID `json:"id" db:"id"`
	// ProjectID references the parent project this service belongs to.
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	// Name is the human-readable name of the service.
	Name string `json:"name" db:"name"`
	// Image is the Docker image and tag to use.
	Image string `json:"image" db:"image"`
	// Environment contains key-value pairs of environment variables passed to
	// the container.
	Environment map[string]string `json:"environment" db:"environment"`
	// Mounts defines the volume and bind mounts for the container.
	Mounts []ServiceMount `json:"mounts" db:"mounts"`
	// Dependencies lists other services that must be running before this
	// service can start.
	Dependencies []ServiceDependency `json:"dependencies" db:"dependencies"`
	// NetworkAccess determines whether the service should be exposed to the
	// external network.
	NetworkAccess bool `json:"network_access" db:"network_access"`
	// ContainerID is the runtime identifier of the container (set after
	// deployment).
	ContainerID *string `json:"container_id" db:"container_id"`
	// CreatedAt is the timestamp when the service was created.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt is the timestamp of the last update.
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
