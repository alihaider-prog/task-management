package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"task-management/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

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
	RETURNING id`

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

	var pgErr *pgconn.PgError
	var id int64
	err := row.Scan(&id)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // Foreign key violation
			if strings.Contains(pgErr.Message, "tasks_project_id_fkey") {
				return nil, errors.New("project does not exist")
			}
			if strings.Contains(pgErr.Message, "tasks_assignee_id_fkey") {
				return nil, errors.New("assignee user does not exist")
			}
		}
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
		WHERE id = $9`

	res, err := tx.Exec(
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
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // Foreign key violation
			if strings.Contains(pgErr.Message, "tasks_assignee_id_fkey") {
				return errors.New("assignee user does not exist")
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("task not found")
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

func (r *TaskRepository) GetProjectIDByTaskID(taskID int64) (*int64, error) {
	query := `SELECT project_id FROM tasks WHERE id = $1`
	row := r.db.QueryRow(query, taskID)
	var projectID int64
	if err := row.Scan(&projectID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &projectID, nil
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

	var changedByValue interface{}
	if changedBy > 0 {
		changedByValue = changedBy
	} else if oldTask.CreatedBy > 0 {
		changedByValue = oldTask.CreatedBy
	} else {
		changedByValue = nil
	}

	_, err := tx.Exec(
		query,
		oldTask.ID,
		changedByValue,
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

func (r *TaskRepository) TaskHistory(taskID int64) ([]models.TaskHistory, error) {
	query := `SELECT id,
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
		new_assignee_id,
		created_at
		FROM task_history
		WHERE task_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.TaskHistory
	for rows.Next() {
		var th models.TaskHistory
		if err := rows.Scan(
			&th.ID,
			&th.TaskID,
			&th.ChangedBy,
			&th.Action,
			&th.FieldName,
			&th.OldValue,
			&th.NewValue,
			&th.OldStatus,
			&th.NewStatus,
			&th.OldPriority,
			&th.NewPriority,
			&th.OldAssigneeID,
			&th.NewAssigneeID,
			&th.CreatedAt,
		); err != nil {
			return nil, err
		}
		history = append(history, th)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

func (r *TaskRepository) MarkOverdueTasks(batchSize int, changedBy int64) (int, error) {
	if batchSize <= 0 {
		batchSize = 100
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	rows, err := tx.Query(`SELECT id,
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
        WHERE due_date < NOW() AND status != 'completed' AND status != 'overdue'
        FOR UPDATE SKIP LOCKED
        LIMIT $1`, batchSize)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer rows.Close()

	var toProcess []*models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.DueDate,
			&t.AssigneeID,
			&t.ParentTaskID,
			&t.Version,
			&t.ProjectID,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			tx.Rollback()
			return 0, err
		}
		toProcess = append(toProcess, &t)
	}
	if err := rows.Err(); err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(toProcess) == 0 {
		if err := tx.Commit(); err != nil {
			return 0, err
		}
		return 0, nil
	}

	for _, old := range toProcess {
		newVersion := old.Version + 1
		// update task status and version
		var updatedAt time.Time
		row := tx.QueryRow(`UPDATE tasks SET status = $1, version = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING updated_at`, models.StatusOverdue, newVersion, old.ID)
		if err := row.Scan(&updatedAt); err != nil {
			tx.Rollback()
			return 0, err
		}

		newTask := *old
		newTask.Status = models.StatusOverdue
		newTask.Version = newVersion
		newTask.UpdatedAt = updatedAt

		if err := r.insertTaskHistoryTx(tx, old, &newTask, changedBy); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return len(toProcess), nil
}
