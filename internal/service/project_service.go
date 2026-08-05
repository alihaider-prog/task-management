package service

import (
	"errors"
	"task-management/internal/models"
	"task-management/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetProjects(workspaceID int64) ([]models.Project, error) {
	return s.repo.AllProjects(workspaceID)
}

func (s *ProjectService) GetProjectByID(id int64) (*models.Project, error) {
	if id <= 0 {
		return nil, errors.New("invalid project id")
	}
	return s.repo.ProjectByID(id)
}

func (s *ProjectService) CreateProject(project *models.Project) (*int64, error) {
	if project.Name == "" {
		return nil, errors.New("name is required")
	}
	if project.WorkspaceID == nil || *project.WorkspaceID <= 0 {
		return nil, errors.New("workspace id is required")
	}

	return s.repo.Create(project)
}

func (s *ProjectService) UpdateProject(project *models.Project) error {
	if project.ID <= 0 {
		return errors.New("invalid project id")
	}
	return s.repo.Update(project)
}

func (s *ProjectService) DeleteProject(id int64) error {
	if id <= 0 {
		return errors.New("invalid project id")
	}
	return s.repo.Delete(id)
}
