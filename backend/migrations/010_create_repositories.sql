-- +goose Up

CREATE TABLE repositories (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    project_id BIGINT UNSIGNED NOT NULL,

    name VARCHAR(255) NOT NULL,
    remote_url VARCHAR(500) NOT NULL,
    provider ENUM(
        'github',
        'gitlab',
        'bitbucket',
        'other'
    ) NOT NULL DEFAULT 'github',

    default_branch VARCHAR(255) NOT NULL DEFAULT 'main',

    status ENUM(
        'active',
        'archived',
        'disconnected'
    ) NOT NULL DEFAULT 'active',

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    UNIQUE KEY uq_repositories_project_url (
        project_id,
        remote_url
    ),

    KEY idx_repositories_project_id (project_id),
    KEY idx_repositories_status (status),

    CONSTRAINT fk_repositories_project
        FOREIGN KEY (project_id)
        REFERENCES projects (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
-- +goose Down

DROP TABLE IF EXISTS repositories;
