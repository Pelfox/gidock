package models

import (
	"time"

	"github.com/google/uuid"
)

// Project is the core domain model of the system. A project represents a
// logical boundary for deployment and management. It groups multiple services
// that collectively implement a complete application.
type Project struct {
	// ID is the unique identifier of the project (UUID).
	ID uuid.UUID `json:"id" db:"id"`
	// Name is the human-readable name of the project.
	Name string `json:"name" db:"name"`
	// CreatedAt is the timestamp when the project was created.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt is the timestamp of the last update to the project.
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ProjectWithServices is an extended model that includes a project and its
// associated services.
type ProjectWithServices struct {
	Project

	// Services contains all services belonging to this project. It is `nil`
	// when services are not loaded.
	Services *[]Service `json:"services,omitempty" db:"services"`
}
