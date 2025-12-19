package dto

import (
	"github.com/Pelfox/gidock/internal/models"
	"github.com/google/uuid"
)

type CreateServiceRequest struct {
	ProjectID uuid.UUID `json:"project_id"`

	Name  string `json:"name"`
	Image string `json:"image"`

	Environment  map[string]string          `json:"environment"`
	Mounts       []models.ServiceMount      `json:"mounts"`
	Dependencies []models.ServiceDependency `json:"dependencies"`

	NetworkAccess bool `json:"network_access"`
}

type CreateServiceResponse struct {
	models.Service
}

// UpdateServiceFields is an internal struct for updating service's fields. It
// shouldn't be used outside (i.e. HTTP requests, etc.)
type UpdateServiceFields struct {
	ContainerID *string
}
