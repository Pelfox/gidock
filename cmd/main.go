package main

import (
	"context"

	"github.com/Pelfox/gidock/internal"
	"github.com/Pelfox/gidock/internal/controllers"
	"github.com/Pelfox/gidock/internal/repositories"
	"github.com/Pelfox/gidock/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	config, err := internal.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	dbPool, err := internal.CreatePool(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	dockerClient, err := client.New(client.FromEnv)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Docker client")
	}

	projectRepository := repositories.NewProjectRepository(dbPool)
	projectService := services.NewProjectService(projectRepository)
	projectController := controllers.NewProjectController(projectService)

	serviceRepository := repositories.NewServiceRepository(dbPool)
	serviceService := services.NewServiceService(serviceRepository, dockerClient)
	serviceController := controllers.NewServiceController(serviceService)

	router := gin.New()

	projectGroup := router.Group("/projects")
	projectGroup.POST("/", projectController.CreateProject)

	serviceGroup := router.Group("/services")
	serviceGroup.GET("/", serviceController.ListServices)
	serviceGroup.POST("/", serviceController.CreateService)
	serviceGroup.POST("/:id/start", serviceController.StartService)
	serviceGroup.POST("/:id/stop", serviceController.StopService)

	if err := router.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
