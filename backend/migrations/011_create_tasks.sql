-- +goose Up

CREATE TABLE tasks (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    project_id BIGINT UNSIGNED NOT NULL,

    name VARCHAR(255) NOT NULL,
    description TEXT NULL,

    assignee_user_id BIGINT UNSIGNED NOT NULL,

    status ENUM(
        'not_started',
        'in_progress',
        'completed',
        'suspended'
    ) NOT NULL DEFAULT 'not_started',

    priority ENUM(
        'high',
        'medium',
        'low'
    ) NOT NULL DEFAULT 'medium',

    start_date DATE NULL,
    due_date DATE NULL,
    completed_at DATETIME NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    KEY idx_tasks_project_id (project_id),
    KEY idx_tasks_assignee_user_id (assignee_user_id),
    KEY idx_tasks_status (status),
    KEY idx_tasks_priority (priority),
    KEY idx_tasks_due_date (due_date),

    CONSTRAINT fk_tasks_project
        FOREIGN KEY (project_id)
        REFERENCES projects (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_tasks_assignee_user
        FOREIGN KEY (assignee_user_id)
        REFERENCES users (id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT chk_tasks_dates
        CHECK (
            due_date IS NULL
            OR start_date IS NULL
            OR due_date >= start_date
        )
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

-- +goose Down

DROP TABLE IF EXISTS tasks;