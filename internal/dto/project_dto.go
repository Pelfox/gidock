package dto

import "github.com/Pelfox/gidock/internal/models"

// CreateProjectRequest is the request payload for creating a new project.
type CreateProjectRequest struct {
	// Name is the human-readable name for the new project.
	Name string `json:"name"`
}

// CreateProjectResponse is the response payload after successfully creating a project.
type CreateProjectResponse struct {
	models.Project
}

// UpdateProjectRequest is the request payload for updating an existing project.
type UpdateProjectRequest struct {
	// Name is the new human-readable name for the project.
	Name *string `json:"name,omitempty"`
}

// UpdateProjectResponse is the response payload after successfully updating a project.
type UpdateProjectResponse struct {
	models.Project
}
