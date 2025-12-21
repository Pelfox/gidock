package controllers

import (
	"net/http"
	"strconv"

	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// TODO: make all endpoints return data via response DTOs
// TODO: handle errors correctly, returning appropriate status codes and messages
// TODO: add other endpoints (from Service)

type ProjectController struct {
	projectService *services.ProjectService
}

func NewProjectController(projectService *services.ProjectService) *ProjectController {
	return &ProjectController{projectService: projectService}
}

func (c *ProjectController) Create(ctx *gin.Context) {
	var request dto.CreateProjectRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
		return
	}

	project, err := c.projectService.Create(ctx.Request.Context(), request)
	if err != nil {
		log.Error().Err(err).Msg("failed to create project")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create project."})
		return
	}

	response := dto.CreateProjectResponse{Project: *project}
	ctx.JSON(http.StatusCreated, response)
}

func (c *ProjectController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided project ID is invalid."})
		return
	}

	includeServices, err := strconv.ParseBool(ctx.DefaultQuery("includeServices", "false"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided `includeServices` flag is invalid."})
		return
	}

	service, err := c.projectService.GetByID(ctx.Request.Context(), id, includeServices)
	if err != nil {
		log.Error().Err(err).Msg("failed to get project")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get project."})
		return
	}

	ctx.JSON(http.StatusOK, service)
}

func (c *ProjectController) ListAll(ctx *gin.Context) {
	projects, err := c.projectService.ListAll(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to list projects")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to list projects."})
		return
	}
	ctx.JSON(http.StatusOK, projects)
}
