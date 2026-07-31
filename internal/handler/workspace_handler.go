package handler

import "task-management/internal/service"

type WorkspaceHandler struct {
	service *service.WorkspaceService
}

func NewWorkspaceHandler(service *service.WorkspaceService) *WorkspaceHandler{
	return &WorkspaceHandler{
		service: service,
	}
}

