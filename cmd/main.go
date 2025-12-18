package main

import (
	"context"

	s "github.com/Masterminds/squirrel"
	"github.com/Pelfox/gidock/internal"
	"github.com/Pelfox/gidock/internal/controllers"
	"github.com/Pelfox/gidock/internal/repositories"
	"github.com/Pelfox/gidock/internal/services"
	"github.com/gin-gonic/gin"
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

	psql := s.StatementBuilder.PlaceholderFormat(s.Dollar)

	projectRepository := repositories.NewProjectRepository(psql, dbPool)
	projectService := services.NewProjectService(projectRepository)
	projectController := controllers.NewProjectController(projectService)

	router := gin.New()

	projectGroup := router.Group("/projects")
	projectGroup.POST("/", projectController.CreateProject)

	if err := router.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
