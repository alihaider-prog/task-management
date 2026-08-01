package models

import "time"

type WorkspaceMemberRole string

const (
	RoleAdmin  WorkspaceMemberRole = "admin"
	RoleMember WorkspaceMemberRole = "member"
	RoleViewer WorkspaceMemberRole = "viewer"
)

type WorkspaceMembers struct {
	ID          int                 `json:"id"`
	WorkspaceID int64               `json:"workspace_id"`
	UserID      int64               `json:"user_id"`
	Role        WorkspaceMemberRole `json:"role"`
	CreatedAt   time.Time           `json:"created_at"`
}

// id
// workspace_id
// user_id
// role (Admin | Member | Viewer)
// created_at
