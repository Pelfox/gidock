package models

import (
	"time"

	"github.com/google/uuid"
)

type ServiceMount struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	ReadOnly bool   `json:"read_only"`
}

type ServiceCondition string

const (
	ServiceConditionHealthy ServiceCondition = "healthy"
	ServiceConditionReady   ServiceCondition = "ready"
)

type ServiceDependency struct {
	ServiceID uuid.UUID        `json:"service_id"`
	Condition ServiceCondition `json:"condition"`
}

type Service struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`

	Name  string `json:"name"`
	Image string `json:"image"`

	Environment  map[string]string   `json:"environment"`
	Mounts       []ServiceMount      `json:"mounts"`
	Dependencies []ServiceDependency `json:"dependencies"`

	NetworkAccess bool       `json:"network_access"`
	ContainerID   *uuid.UUID `json:"container_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
