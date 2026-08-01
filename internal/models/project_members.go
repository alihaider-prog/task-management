package models

import (
	"time"
)

type ProjectMembers struct {
	ID        int       `json:"id"`
	ProjectID int64     `json:"project_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// id
// project_id
// user_id
// created_at
