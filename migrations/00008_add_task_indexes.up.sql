CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_tasks_project_status_priority_assignee_due_date
ON tasks (project_id, status, priority, assignee_id, due_date);

CREATE INDEX IF NOT EXISTS idx_tasks_project_due_date
ON tasks (project_id, due_date);

CREATE INDEX IF NOT EXISTS idx_tasks_title_description_trgm
ON tasks USING gin ((title || ' ' || coalesce(description, '')) gin_trgm_ops);
