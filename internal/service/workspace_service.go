package service

import "task-management/internal/repository"

type WorkspaceService struct {
	repo *repository.WorkspaceRepository
}

func NewWorkspaceService(repo *repository.WorkspaceRepository) *WorkspaceService{
	return &WorkspaceService{
		repo: repo,
	}
}