package models

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusOverdue    TaskStatus = "overdue"
)

type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
	PriorityUrgent TaskPriority = "urgent"
)

type Task struct {
	ID           int64        `json:"id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Status       TaskStatus   `json:"status"`
	Priority     TaskPriority `json:"priority"`
	DueDate      *time.Time   `json:"due_date"`
	AssigneeID   *int64       `json:"assignee_id"`
	ProjectID    *int64       `json:"project_id"`
	ParentTaskID *int64       `json:"parent_task_id"`
	ParentTask   *Task        `json:"parent_task,omitempty"`
	Subtasks     []Task       `json:"subtasks,omitempty"`
	Version      int          `json:"version"`
	CreatedBy    int64        `json:"created_by"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// id
// title
// description
// status  (todo, in_progress, completed, overdue)
// priority (low, medium, high, urgent)
// due_date
// assignee_id
// parent_task_id
// project_id
// created_by
// created_at
// updated_at
