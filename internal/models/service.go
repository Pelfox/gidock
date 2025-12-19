package models

import (
	"time"

	"github.com/google/uuid"
)

type ServiceMount struct {
	// TODO: specify whether it is a volume or a host mount
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
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`

	Name  string `json:"name" db:"name"`
	Image string `json:"image" db:"image"`

	Environment  map[string]string   `json:"environment" db:"environment"`
	Mounts       []ServiceMount      `json:"mounts" db:"mounts"`
	Dependencies []ServiceDependency `json:"dependencies" db:"dependencies"`

	NetworkAccess bool    `json:"network_access" db:"network_access"`
	ContainerID   *string `json:"container_id" db:"container_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
