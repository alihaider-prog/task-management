package repository

import (
	"database/sql"
	"errors"
	"task-management/internal/models"
)

// import "task-management/internal/models"

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) AllTasks(projectID int64) ([]models.Task, error) {
	query := `SELECT id,
		title,
		description,
		status,
		priority,
		due_date,
		assignee_id,
		parent_task_id,
		version,
		created_by,
		created_at,
		updated_at
		FROM tasks
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.AssigneeID,
			&task.ParentTaskID,
			&task.Version,
			&task.CreatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil

}

func (r *TaskRepository) TaskByID(id int64) (*models.Task, error) {
	query := `SELECT id,
		title,
		description,
		status,
		priority,
		due_date,
		assignee_id,
		parent_task_id,
		version,
		project_id,
		created_by,
		created_at,
		updated_at
		FROM tasks
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var task models.Task

	if err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.AssigneeID,
		&task.ParentTaskID,
		&task.Version,
		&task.ProjectID,
		&task.CreatedBy,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Create(task *models.Task) (*int64, error) {
	query := `INSERT INTO tasks(
		title,
		description,
		status,
		priority,
		due_date,
		assignee_id,
		parent_task_id,
		version,
		project_id,
		created_by
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	RETURNING id, created_at, updated_at`

	row := r.db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.AssigneeID,
		task.ParentTaskID,
		task.Version,
		task.ProjectID,
		task.CreatedBy,
	)

	var id int64
	if err := row.Scan(&id, &task.CreatedAt, &task.UpdatedAt); err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	// increment version and update fields
	newVersion := task.Version + 1

	query := `UPDATE tasks SET
		title = $1,
		description = $2,
		status = $3,
		priority = $4,
		due_date = $5,
		assignee_id = $6,
		parent_task_id = $7,
		version = $8,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING updated_at`

	row := r.db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.AssigneeID,
		task.ParentTaskID,
		newVersion,
		task.ID,
	)

	if err := row.Scan(&task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("task not found")
		}
		return err
	}

	task.Version = newVersion
	return nil
}

func (r *TaskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *TaskRepository) ProjectExists(projectID int64) (bool, error) {
	query := `SELECT 1 FROM projects WHERE id = $1`
	row := r.db.QueryRow(query, projectID)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
