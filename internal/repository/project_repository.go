package repository

import (
	"database/sql"
	"errors"
	"task-management/internal/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) AllProjects(workspaceID int64) ([]models.Project, error) {
	query := `SELECT id,
		workspace_id,
		name,
		description,
		created_by,
		created_at,
		updated_at
		FROM projects
		WHERE workspace_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		var workspaceIDValue int64
		if err := rows.Scan(
			&project.ID,
			&workspaceIDValue,
			&project.Name,
			&project.Description,
			&project.CreatedBy,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, err
		}
		project.WorkspaceID = &workspaceIDValue
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) ProjectByID(id int64) (*models.Project, error) {
	query := `SELECT id,
	 	workspace_id,
	 	name,
		description,
		created_by,
		created_at,
		updated_at
		FROM projects
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)

	var project models.Project
	var workspaceIDValue int64
	if err := row.Scan(&project.ID,
		&workspaceIDValue,
		&project.Name,
		&project.Description,
		&project.CreatedBy,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	project.WorkspaceID = &workspaceIDValue

	return &project, nil
}

func (r *ProjectRepository) Create(project *models.Project) (*int64, error) {
	query := `INSERT INTO projects (
		workspace_id,
		name,
		description,
		created_by
	) VALUES ($1, $2, $3, $4)
	 RETURNING id, created_by, created_at, updated_at`

	row := r.db.QueryRow(
		query,
		*project.WorkspaceID,
		project.Name,
		project.Description,
		project.CreatedBy,
	)

	var id int64
	if err := row.Scan(
		&id,
		&project.CreatedBy,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *ProjectRepository) Update(project *models.Project) error {
	query := `UPDATE projects SET
		name = $1,
		description = $2
		WHERE id = $3 
		RETURNING updated_at
	`

	row := r.db.QueryRow(query,
		project.Name,
		project.Description,
		project.ID,
	)

	if err := row.Scan(&project.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("project not found")
		}
		return err
	}

	return nil
}

func (r *ProjectRepository) Delete(id int64) error {
	query := `DELETE FROM projects WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("project not found")
	}

	return nil
}

func (r *ProjectRepository) WorkspaceExists(workspaceID int64) (bool, error) {
	query := `SELECT 1 FROM workspaces WHERE id = $1`
	row := r.db.QueryRow(query, workspaceID)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ProjectRepository) GetWorkspaceIDByProjectID(projectID int64) (*int64, error) {
	query := `SELECT workspace_id FROM projects WHERE id = $1`
	row := r.db.QueryRow(query, projectID)
	var workspaceID int64
	if err := row.Scan(&workspaceID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	return &workspaceID, nil
}
