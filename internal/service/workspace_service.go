package service

import (
	"errors"
	"task-management/internal/models"
	"task-management/internal/repository"
)

type WorkspaceService struct {
	repo *repository.WorkspaceRepository
}

func NewWorkspaceService(repo *repository.WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{
		repo: repo,
	}
}

func (s *WorkspaceService) GetWorkspaces() ([]models.Workspace, error) {
	return s.repo.AllWorkspaces()
}

func (s *WorkspaceService) GetWorkspaceByID(id int64) (*models.Workspace, error) {
	if id <= 0 {
		return nil, errors.New("invalid workspace id")
	}

	return s.repo.WorkspaceByID(id)
}

func (s *WorkspaceService) CreateWorkspace(workspace *models.Workspace, userID int64) (*int64, error) {
	if workspace.Name == "" {
		return nil, errors.New("name is required")
	}

	workspace.OwnerID = userID

	return s.repo.Create(workspace)
}

func (s *WorkspaceService) UpdateWorkspace(workspace *models.Workspace) error {
	if workspace.ID <= 0 {
		return errors.New("invalid workspace")
	}
	return s.repo.Update(workspace)
}

func (s *WorkspaceService) DeleteWorkspace(id int64) error {
	if id <= 0 {
		return errors.New("invalid workspace id")
	}
	return s.repo.Delete(id)
}
