package dto

import "github.com/Pelfox/gidock/internal/models"

type CreateProjectRequest struct {
	Name string `json:"name"`
}

type CreateProjectResponse struct {
	models.Project
}
