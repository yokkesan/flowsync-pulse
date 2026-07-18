-- +goose Up

CREATE TABLE tasks (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    project_id BIGINT UNSIGNED NOT NULL,

    ticket_key VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,

    status ENUM(
        'todo',
        'in_progress',
        'review',
        'completed',
        'archived'
    ) NOT NULL DEFAULT 'todo',

    assigned_user_id BIGINT UNSIGNED NULL,

    started_at DATETIME NULL,
    completed_at DATETIME NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    UNIQUE KEY uq_tasks_project_ticket_key (
        project_id,
        ticket_key
    ),

    KEY idx_tasks_project_id (project_id),
    KEY idx_tasks_assigned_user_id (assigned_user_id),
    KEY idx_tasks_status (status),

    CONSTRAINT fk_tasks_project
        FOREIGN KEY (project_id)
        REFERENCES projects (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_tasks_assigned_user
        FOREIGN KEY (assigned_user_id)
        REFERENCES users (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    CONSTRAINT chk_tasks_dates
        CHECK (
            completed_at IS NULL
            OR started_at IS NULL
            OR completed_at >= started_at
        )
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
-- +goose Down

DROP TABLE IF EXISTS tasks;
