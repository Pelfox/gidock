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
