package service

import (
	"errors"
	"task-management/internal/models"
	"task-management/internal/repository"
)

type MemberService struct {
	repo *repository.MemberRepository
}

func NewMemberService(repo *repository.MemberRepository) *MemberService {
	return &MemberService{repo: repo}
}

func (s *MemberService) AddWorkspaceMember(member *models.WorkspaceMembers) (*int64, error) {
	if member.WorkspaceID <= 0 {
		return nil, errors.New("workspace id is required")
	}

	if member.Role == "" || (member.Role != models.RoleAdmin && member.Role != models.RoleMember && member.Role != models.RoleViewer) {
		return nil, errors.New("role is required or invalid")
	}
	return s.repo.AddWorkspaceMember(member)
}

func (s *MemberService) ListWorkspaceMembers(workspaceID int64) ([]models.WorkspaceMembers, error) {
	if workspaceID <= 0 {
		return nil, errors.New("invalid workspace id")
	}
	return s.repo.ListWorkspaceMembers(workspaceID)
}

func (s *MemberService) RemoveWorkspaceMember(id int64) error {
	if id <= 0 {
		return errors.New("invalid member id")
	}
	return s.repo.RemoveWorkspaceMember(id)
}

func (s *MemberService) AddProjectMember(member *models.ProjectMembers) (*int64, error) {
	if member.ProjectID <= 0 {
		return nil, errors.New("project id is required")
	}

	return s.repo.AddProjectMember(member)
}

func (s *MemberService) ListProjectMembers(projectID int64) ([]models.ProjectMembers, error) {
	if projectID <= 0 {
		return nil, errors.New("invalid project id")
	}
	return s.repo.ListProjectMembers(projectID)
}

func (s *MemberService) RemoveProjectMember(id int64) error {
	if id <= 0 {
		return errors.New("invalid member id")
	}
	return s.repo.RemoveProjectMember(id)
}
