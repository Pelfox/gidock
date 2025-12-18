package controllers

import (
	"net/http"

	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/services"
	"github.com/gin-gonic/gin"
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
