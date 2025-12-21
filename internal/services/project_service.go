package services

import (
	"context"

	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/internal/repositories"
	"github.com/Pelfox/gidock/internal/repositories/commands"
	"github.com/google/uuid"
)

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

func (s *ProjectService) Update(
	ctx context.Context,
	id uuid.UUID,
	request dto.UpdateProjectRequest,
) (*models.Project, error) {
	return s.projectRepository.Update(ctx, commands.UpdateProjectCommand{
		ID:   id,
		Name: request.Name,
	})
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.projectRepository.Delete(ctx, commands.DeleteProjectCommand{ID: id})
}

func (s *ProjectService) ListAll(ctx context.Context) ([]models.Project, error) {
	projects, err := s.projectRepository.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
