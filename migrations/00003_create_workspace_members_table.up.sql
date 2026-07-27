CREATE TABLE workspace_members (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    workspace_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'member', 'viewer')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(workspace_id, user_id)
);

-- id
-- workspace_id
-- user_id
-- role (Admin | Member | Viewer)
-- created_at