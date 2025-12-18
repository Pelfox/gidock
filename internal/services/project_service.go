package services

import (
	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/internal/repositories"
)

type ProjectService struct {
	projectRepository *repositories.ProjectRepository
}

func NewProjectService(projectRepository *repositories.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepository: projectRepository}
}

func (s *ProjectService) Create(request dto.CreateProjectRequest) (*models.Project, error) {
	project, err := s.projectRepository.CreateProject(request.Name)
	if err != nil {
		return nil, err
	}
	return project, nil
}
