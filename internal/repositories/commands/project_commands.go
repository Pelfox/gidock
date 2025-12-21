package commands

import "github.com/google/uuid"

// CreateProjectCommand represents the data required to create a new project.
type CreateProjectCommand struct {
	// Name is the human-readable name for the new project.
	Name string
}

// GetProjectCommand represents the data required to retrieve a project.
type GetProjectCommand struct {
	// ID is the unique identifier of the project to be retrieved.
	ID uuid.UUID
	// IncludeServices indicates whether to include associated services.
	IncludeServices bool
}

// UpdateProjectCommand represents a partial update request for a project.
type UpdateProjectCommand struct {
	// ID is the unique identifier of the project to be updated.
	ID uuid.UUID
	// Name is the new human-readable name for the project.
	Name *string
}

// DeleteProjectCommand represents the data required to delete a project.
type DeleteProjectCommand struct {
	// ID is the unique identifier of the project to be deleted.
	ID uuid.UUID
}
