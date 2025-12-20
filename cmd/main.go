package main

import (
	"context"

	"github.com/Pelfox/gidock/internal"
	"github.com/Pelfox/gidock/internal/controllers"
	"github.com/Pelfox/gidock/internal/repositories"
	"github.com/Pelfox/gidock/internal/services"
	"github.com/gin-contrib/cors"
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

	dockerService := services.NewDockerService(dockerClient)

	projectRepository := repositories.NewProjectRepository(dbPool)
	projectService := services.NewProjectService(projectRepository)
	projectController := controllers.NewProjectController(projectService)

	serviceRepository := repositories.NewServiceRepository(dbPool)
	serviceService := services.NewServiceService(serviceRepository, dockerService)
	serviceController := controllers.NewServiceController(serviceService)

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	projectGroup := router.Group("/projects")
	projectGroup.POST("/", projectController.CreateProject)
	projectGroup.GET("/", projectController.ListProjects)
	projectGroup.GET("/:id", projectController.GetProjectByID)

	serviceGroup := router.Group("/services")
	serviceGroup.GET("/", serviceController.ListServices)
	serviceGroup.POST("/", serviceController.CreateService)
	serviceGroup.GET("/:id", serviceController.GetServiceByID)
	serviceGroup.POST("/:id/start", serviceController.StartService)
	serviceGroup.POST("/:id/stop", serviceController.StopService)
	serviceGroup.GET("/:id/status", serviceController.GetServiceStatus)
	serviceGroup.GET("/:id/logs", serviceController.GetServiceLogs)

	if err := router.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
