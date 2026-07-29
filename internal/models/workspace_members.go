package models

import "time"

type WorkspaceMemberRole string

const (
	RoleAdmin  WorkspaceMemberRole = "Admin"
	RoleMember WorkspaceMemberRole = "Member"
	RoleViewer WorkspaceMemberRole = "Viewer"
)

type WorkspaceMembers struct {
	ID          int                 `json:"id"`
	WorkspaceID int                 `json:"workspace_id"`
	UserID      int                 `json:"user_id"`
	Role        WorkspaceMemberRole `json:"role"`
	CreatedAt   time.Time           `json:"created_at"`
}

// id
// workspace_id
// user_id
// role (Admin | Member | Viewer)
// created_at
