package services

import (
	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/internal/repositories"
	"github.com/google/uuid"
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

func (s *ProjectService) ListProjects() ([]models.Project, error) {
	projects, err := s.projectRepository.ListProjects()
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *ProjectService) GetProjectByID(id uuid.UUID, includeServices bool) (*models.ProjectWithServices, error) {
	service, err := s.projectRepository.GetProjectById(id, includeServices)
	if err != nil {
		return nil, err
	}
	return service, err
}
