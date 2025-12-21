package services

import (
	"context"

	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/internal/repositories"
	"github.com/Pelfox/gidock/internal/repositories/commands"
	"github.com/google/uuid"
)

// TODO: add other methods (from Repository)

type ProjectService struct {
	projectRepository *repositories.ProjectRepository
}

func NewProjectService(projectRepository *repositories.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepository: projectRepository}
}

func (s *ProjectService) Create(
	ctx context.Context,
	request dto.CreateProjectRequest,
) (*models.Project, error) {
	return s.projectRepository.Create(ctx, commands.CreateProjectCommand{
		Name: request.Name,
	})
}

func (s *ProjectService) GetByID(
	ctx context.Context,
	id uuid.UUID,
	includeServices bool,
) (*models.ProjectWithServices, error) {
	return s.projectRepository.Get(ctx, commands.GetProjectCommand{
		ID:              id,
		IncludeServices: includeServices,
	})
}

func (s *ProjectService) ListAll(ctx context.Context) ([]models.Project, error) {
	projects, err := s.projectRepository.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
