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

func (c *ProjectController) ListProjects(ctx *gin.Context) {
	projects, err := c.projectService.ListProjects()
	if err != nil {
		log.Error().Err(err).Msg("failed to list projects")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to list projects."})
		return
	}
	ctx.JSON(http.StatusOK, projects)
}

func (c *ProjectController) GetProjectByID(ctx *gin.Context) {
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

	service, err := c.projectService.GetProjectByID(id, includeServices)
	if err != nil {
		log.Error().Err(err).Msg("failed to get project")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get project."})
		return
	}

	ctx.JSON(http.StatusOK, service)
}
