CREATE TABLE task_history (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    task_id INTEGER,
    changed_by INTEGER,
    action VARCHAR(20) NOT NULL CHECK (action IN ('create', 'update', 'delete', 'status_change')),
    field_name VARCHAR(100),
    old_value TEXT,
    new_value TEXT,
    old_status VARCHAR(20) NOT NULL CHECK (old_status IN ('todo', 'in_progress', 'completed', 'overdue')),
    new_status VARCHAR(20) NOT NULL CHECK (new_status IN ('todo', 'in_progress', 'completed', 'overdue')),
    old_priority VARCHAR(20) NOT NULL CHECK (old_priority IN ('low', 'medium', 'high', 'urgent')),
    new_priority VARCHAR(20) NOT NULL CHECK (new_priority IN ('low', 'medium', 'high', 'urgent')),
    old_assignee_id INTEGER,
    new_assignee_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (changed_by) REFERENCES users(id),
    FOREIGN KEY (old_assignee_id) REFERENCES users(id),
    FOREIGN KEY (old_assignee_id) REFERENCES users(id)
);


-- id
-- task_id
-- changed_by
-- action(create, update, delete, status_change)
-- field_name
-- old_value
-- new_value
-- created_at
-- old_status
-- new_status
-- old_priority
-- new_priority
-- old_assignee
-- new_assignee
