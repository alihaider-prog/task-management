package models

import "time"

type TaskAction string

const (
	ActionCreate       TaskAction = "create"
	ActionUpdate       TaskAction = "update"
	ActionDelete       TaskAction = "delete"
	ActionChangeStatus TaskAction = "status_change"
)

type TaskHistory struct {
	ID            int64        `json:"id"`
	TaskID        *int64       `json:"task_id,omitempty"`
	ChangeddBy    *int64       `json:"changed_by,omitempty"`
	Action        TaskAction   `json:"action"`
	FieldName     string       `json:"field_name"`
	OldValue      string       `json:"old_value"`
	NewValue      string       `json:"new_value"`
	OldStatus     TaskStatus   `json:"old_status"`
	NewStatus     TaskStatus   `json:"new_status"`
	OldPriority   TaskPriority `json:"old_priority"`
	NewPriority   TaskPriority `json:"new_priority"`
	OldAssigneeID *int64       `json:"old_assignee_id,omitempty"`
	NewAssigneeID *int64       `json:"new_assignee_id,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}

// id
// task_id
// changed_by
// action(create, update, delete, status_change)
// field_name
// old_value
// new_value
// old_status
// new_status
// old_priority
// new_priority
// old_assignee_id
// new_assignee_id
// created_at
