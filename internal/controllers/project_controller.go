package controllers

import (
	"net/http"

	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ProjectController struct {
	projectService *services.ProjectService
}

func NewProjectController(projectService *services.ProjectService) *ProjectController {
	return &ProjectController{projectService: projectService}
}

func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var request dto.CreateProjectRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
		return
	}

	project, err := c.projectService.Create(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to create project")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create project."})
		return
	}

	response := dto.CreateProjectResponse{Project: *project}
	ctx.JSON(http.StatusCreated, response)
}
