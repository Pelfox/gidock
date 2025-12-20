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

type ServiceController struct {
	serviceService *services.ServiceService
}

func NewServiceController(serviceService *services.ServiceService) *ServiceController {
	return &ServiceController{serviceService: serviceService}
}

func (c *ServiceController) CreateService(ctx *gin.Context) {
	var request dto.CreateServiceRequest
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
		return
	}

	service, err := c.serviceService.CreateService(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to create service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service."})
		return
	}

	response := dto.CreateServiceResponse{Service: *service}
	ctx.JSON(http.StatusCreated, response)
}

func (c *ServiceController) StartService(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided service ID is invalid."})
		return
	}

	forcePull, err := strconv.ParseBool(ctx.DefaultQuery("forcePull", "false"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided `forcePull` flag is invalid."})
		return
	}

	service, err := c.serviceService.StartService(ctx.Request.Context(), id, forcePull)
	if err != nil {
		log.Error().Err(err).Msg("failed to start service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start service."})
		return
	}

	ctx.JSON(http.StatusOK, service)
}

func (c *ServiceController) StopService(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided service ID is invalid."})
		return
	}

	kill, err := strconv.ParseBool(ctx.DefaultQuery("kill", "false"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided `kill` flag is invalid."})
		return
	}

	if err = c.serviceService.StopService(ctx.Request.Context(), id, kill); err != nil {
		log.Error().Err(err).Msg("failed to stop service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to stop service."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OK."})
}

func (c *ServiceController) ListServices(ctx *gin.Context) {
	servicesList, err := c.serviceService.ListServices()
	if err != nil {
		log.Error().Err(err).Msg("failed to list services")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to list services."})
		return
	}
	ctx.JSON(http.StatusOK, servicesList)
}

func (c *ServiceController) GetServiceStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided service ID is invalid."})
		return
	}

	status, err := c.serviceService.GetServiceStatus(ctx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get service status")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get service status."})
		return
	}

	ctx.JSON(http.StatusOK, status)
}
