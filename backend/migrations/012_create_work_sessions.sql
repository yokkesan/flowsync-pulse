-- +goose Up

CREATE TABLE work_sessions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,

    user_id BIGINT UNSIGNED NOT NULL,
    project_id BIGINT UNSIGNED NOT NULL,
    repository_id BIGINT UNSIGNED NOT NULL,
    task_id BIGINT UNSIGNED NULL,

    branch_name VARCHAR(255) NOT NULL,
    ticket_key VARCHAR(100) NULL,
    workspace_name VARCHAR(255) NULL,

    status ENUM(
        'active',
        'completed',
        'timed_out'
    ) NOT NULL DEFAULT 'active',

    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_heartbeat_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at DATETIME NULL,

    end_reason ENUM(
        'branch_changed',
        'repository_changed',
        'client_closed',
        'timeout',
        'manual'
    ) NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    KEY idx_work_sessions_user_id (user_id),
    KEY idx_work_sessions_project_id (project_id),
    KEY idx_work_sessions_repository_id (repository_id),
    KEY idx_work_sessions_task_id (task_id),
    KEY idx_work_sessions_status (status),
    KEY idx_work_sessions_last_heartbeat_at (last_heartbeat_at),

    CONSTRAINT fk_work_sessions_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_work_sessions_project
        FOREIGN KEY (project_id)
        REFERENCES projects (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_work_sessions_repository
        FOREIGN KEY (repository_id)
        REFERENCES repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_work_sessions_task
        FOREIGN KEY (task_id)
        REFERENCES tasks (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    CONSTRAINT chk_work_sessions_dates
        CHECK (
            ended_at IS NULL
            OR ended_at >= started_at
        )
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
-- +goose Down

DROP TABLE IF EXISTS work_sessions;
