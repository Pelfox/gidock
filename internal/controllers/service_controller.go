package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/services"
	"github.com/Pelfox/gidock/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// TODO: make all endpoints return data via response DTOs
// TODO: handle errors correctly, returning appropriate status codes and messages
// TODO: add other endpoints (from Service)

type ServiceController struct {
	serviceService *services.ServiceService
}

func NewServiceController(serviceService *services.ServiceService) *ServiceController {
	return &ServiceController{serviceService: serviceService}
}

func (c *ServiceController) Create(ctx *gin.Context) {
	var request dto.CreateServiceRequest
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
		return
	}

	service, err := c.serviceService.Create(ctx.Request.Context(), request)
	if err != nil {
		log.Error().Err(err).Msg("failed to create service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service."})
		return
	}

	response := dto.CreateServiceResponse{Service: *service}
	ctx.JSON(http.StatusCreated, response)
}

func (c *ServiceController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided service ID is invalid."})
		return
	}

	service, err := c.serviceService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get service."})
		return
	}

	ctx.JSON(http.StatusOK, service)
}

func (c *ServiceController) ListAll(ctx *gin.Context) {
	servicesList, err := c.serviceService.ListAll(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to list services")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to list services."})
		return
	}
	ctx.JSON(http.StatusOK, servicesList)
}

func (c *ServiceController) Start(ctx *gin.Context) {
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

	service, err := c.serviceService.Start(ctx.Request.Context(), id, forcePull)
	if err != nil {
		log.Error().Err(err).Msg("failed to start service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start service."})
		return
	}

	ctx.JSON(http.StatusOK, service)
}

func (c *ServiceController) Stop(ctx *gin.Context) {
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

	if err = c.serviceService.Stop(ctx.Request.Context(), id, kill); err != nil {
		log.Error().Err(err).Msg("failed to stop service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to stop service."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OK."})
}

func (c *ServiceController) GetStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided service ID is invalid."})
		return
	}

	status, err := c.serviceService.GetStatus(ctx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get service status")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get service status."})
		return
	}

	ctx.JSON(http.StatusOK, status)
}

func (c *ServiceController) StreamLogs(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The provided service ID is invalid."})
		return
	}

	logsChannel, err := c.serviceService.StreamLogs(ctx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get service logs")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get service logs."})
		return
	}

	conn := pkg.NewSSEConn(ctx, 10*time.Second)
	conn.SetupHeaders()
	conn.StartHeartbeats()
	defer conn.Close()

	for {
		select {
		case line, ok := <-logsChannel:
			if !ok {
				return
			}
			if err := conn.SendEvent("log", line); err != nil {
				log.Error().Err(err).Msg("failed to send log event")
			}
		case <-ctx.Request.Context().Done():
			return
		}
	}
}
