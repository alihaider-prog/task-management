package repository

import (
	"database/sql"
	"errors"
	"task-management/internal/models"
)

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

func (r *WorkspaceRepository) AllWorkspaces() ([]models.Workspace, error) {
	query := `
		SELECT id,
		name,
		owner_id,
		created_at,
		updated_at
		FROM workspaces
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var workspaces []models.Workspace
	for rows.Next() {
		var workspace models.Workspace

		err = rows.Scan(
			&workspace.ID,
			&workspace.Name,
			&workspace.OwnerID,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) WorkspaceByID(id int64) (*models.Workspace, error) {
	query := `
		SELECT id,
		name,
		owner_id,
		created_at,
		updated_at
		FROM workspaces
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var workspace models.Workspace

	if err := row.Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.OwnerID,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("workspace not found")
		}
		return nil, err
	}

	return &workspace, nil
}

func (r *WorkspaceRepository) Create(workspace *models.Workspace) (*int64, error) {
	query := `INSERT INTO workspaces(
			name,
			owner_id
		) VALUES ($1, $2)
		RETURNING id, created_at, updated_at 
	`

	row := r.db.QueryRow(
		query,
		workspace.Name,
		workspace.OwnerID,
	)

	var id int64
	if err := row.Scan(
		&id,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *WorkspaceRepository) Update(workspace *models.Workspace) error {
	query := `UPDATE workspaces SET name = $1 WHERE id = $2 RETURNING updated_at`

	row := r.db.QueryRow(query, workspace.Name, workspace.ID)

	if err := row.Scan(&workspace.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("workspace not found")
		}
		return err
	}

	return nil
}

func (r *WorkspaceRepository) Delete(id int64) error {
	query := `DELETE FROM workspaces WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("workspace not found")
	}

	return nil
}
