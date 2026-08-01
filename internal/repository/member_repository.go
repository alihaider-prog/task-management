package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"task-management/internal/models"
)

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) AddWorkspaceMember(member *models.WorkspaceMembers) (*int64, error) {
	query := `INSERT INTO workspace_members(
		workspace_id,
		user_id,
		role
	) VALUES ($1, $2, $3)
	RETURNING id, created_at`

	row := r.db.QueryRow(
		query,
		member.WorkspaceID,
		member.UserID,
		member.Role,
	)

	var id int64
	if err := row.Scan(&id, &member.CreatedAt); err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *MemberRepository) ListWorkspaceMembers(workspaceID int64) ([]models.WorkspaceMembers, error) {
	query := `SELECT id, workspace_id, user_id, role, created_at FROM workspace_members WHERE workspace_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.WorkspaceMembers
	for rows.Next() {
		var member models.WorkspaceMembers
		if err := rows.Scan(
			&member.ID,
			&member.WorkspaceID,
			&member.UserID,
			&member.Role,
			&member.CreatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *MemberRepository) RemoveWorkspaceMember(id int64) error {
	query := `DELETE FROM workspace_members WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("workspace member not found")
	}

	return nil
}

func (r *MemberRepository) AddProjectMember(member *models.ProjectMembers) (*int64, error) {
	query := `INSERT INTO project_members(
		project_id,
		user_id
	) VALUES ($1, $2)
	RETURNING id, created_at`

	row := r.db.QueryRow(
		query,
		member.ProjectID,
		member.UserID,
	)

	var id int64
	if err := row.Scan(&id, &member.CreatedAt); err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *MemberRepository) ListProjectMembers(projectID int64) ([]models.ProjectMembers, error) {
	query := `SELECT id, project_id, user_id, created_at FROM project_members WHERE project_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectMembers
	for rows.Next() {
		var member models.ProjectMembers
		if err := rows.Scan(
			&member.ID,
			&member.ProjectID,
			&member.UserID,
			&member.CreatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *MemberRepository) RemoveProjectMember(id int64) error {
	query := `DELETE FROM project_members WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("project member not found")
	}

	return nil
}

func (r *MemberRepository) IsProjectMember(projectID, userID int64) (bool, error) {
	query := `SELECT 1 FROM project_members WHERE project_id = $1 AND user_id = $2`
	row := r.db.QueryRow(query, projectID, userID)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *MemberRepository) ModelExists(modelID int64, tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = $1", tableName)
	row := r.db.QueryRow(query, modelID)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
