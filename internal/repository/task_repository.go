package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"task-management/internal/models"
)

// import "task-management/internal/models"

type TaskRepository struct {
	db *sql.DB
}

type TaskFilter struct {
	ProjectID  int64
	Status     string
	Priority   string
	AssigneeID *int64
	Search     string
	SortBy     string
	Order      string
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) AllTasks(filter TaskFilter) ([]models.Task, error) {
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
	`
	params := []interface{}{filter.ProjectID}

	if filter.Status != "" {
		query += " AND status = $" + fmt.Sprint(len(params)+1)
		params = append(params, filter.Status)
	}
	if filter.AssigneeID != nil {
		query += " AND assignee_id = $" + fmt.Sprint(len(params)+1)
		params = append(params, *filter.AssigneeID)
	}
	if filter.Priority != "" {
		query += " AND priority = $" + fmt.Sprint(len(params)+1)
		params = append(params, filter.Priority)
	}
	if filter.Search != "" {
		query += " AND (title ILIKE $" + fmt.Sprint(len(params)+1) + " OR description ILIKE $" + fmt.Sprint(len(params)+2) + ")"
		searchTerm := "%" + strings.ReplaceAll(filter.Search, "%", "\\%") + "%"
		params = append(params, searchTerm, searchTerm)
	}

	sortBy := "created_at"
	allowedSort := map[string]string{
		"due_date":   "due_date",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	if col, ok := allowedSort[filter.SortBy]; ok {
		sortBy = col
	}

	order := "DESC"
	if strings.EqualFold(filter.Order, "asc") {
		order = "ASC"
	}

	query += " ORDER BY " + sortBy + " " + order

	rows, err := r.db.Query(query, params...)
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
	task, err := r.fetchTaskByID(id)
	if err != nil {
		return nil, err
	}

	if task.ParentTaskID != nil {
		parentTask, err := r.fetchTaskByID(*task.ParentTaskID)
		if err != nil {
			return nil, err
		}
		task.ParentTask = parentTask
	}

	subtasks, err := r.fetchSubtasks(id)
	if err != nil {
		return nil, err
	}
	task.Subtasks = subtasks

	return task, nil
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

func (r *TaskRepository) Update(task *models.Task, changedBy int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	currentTask, err := r.fetchTaskByIDTx(tx, task.ID)
	if err != nil {
		return err
	}

	newVersion := currentTask.Version + 1
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

	row := tx.QueryRow(
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

	if err = row.Scan(&task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("task not found")
		}
		return err
	}

	if err = r.insertTaskHistoryTx(tx, currentTask, task, changedBy); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
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

func (r *TaskRepository) UserExists(userID int64) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = $1`
	row := r.db.QueryRow(query, userID)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *TaskRepository) fetchTaskByID(id int64) (*models.Task, error) {
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

func (r *TaskRepository) fetchSubtasks(parentTaskID int64) ([]models.Task, error) {
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
		WHERE parent_task_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, parentTaskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
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
			return nil, err
		}
		subtasks = append(subtasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subtasks, nil
}

func (r *TaskRepository) fetchTaskByIDTx(tx *sql.Tx, id int64) (*models.Task, error) {
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
		FOR UPDATE`
	row := tx.QueryRow(query, id)
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

func (r *TaskRepository) insertTaskHistoryTx(tx *sql.Tx, oldTask, newTask *models.Task, changedBy int64) error {
	query := `INSERT INTO task_history(
		task_id,
		changed_by,
		action,
		field_name,
		old_value,
		new_value,
		old_status,
		new_status,
		old_priority,
		new_priority,
		old_assignee_id,
		new_assignee_id
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	fieldName := "update"
	oldValue := oldTask.Title + "\n" + oldTask.Description
	newValue := newTask.Title + "\n" + newTask.Description

	_, err := tx.Exec(
		query,
		oldTask.ID,
		changedBy,
		"update",
		fieldName,
		oldValue,
		newValue,
		oldTask.Status,
		newTask.Status,
		oldTask.Priority,
		newTask.Priority,
		oldTask.AssigneeID,
		newTask.AssigneeID,
	)
	return err
}
