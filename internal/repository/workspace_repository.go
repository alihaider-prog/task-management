package repository

import "database/sql"

type WorkspaceRepository struct {
	db *sql.DB
}

func NewWorkspaceRepository(db *sql.DB) *WorkspaceRepository {
	return &WorkspaceRepository{
		db: db,
	}
}

func (r *WorkspaceRepository) GetUserRole(workspaceID, userID int64) (string, error) {
	query := `
		SELECT role
		FROM workspace_members
		WHERE workspace_id $1 AND user_id = $2
	`

	var role string

	err := r.db.QueryRow(
		query, workspaceID, userID,
	).Scan(
		&role,
	)

	return role, err
}
