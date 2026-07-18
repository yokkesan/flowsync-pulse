-- +goose Up

CREATE TABLE projects (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    company_id BIGINT UNSIGNED NOT NULL,
    team_id BIGINT UNSIGNED NULL,

    name VARCHAR(150) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT NULL,

    status ENUM(
        'planned',
        'active',
        'completed',
        'archived'
    ) NOT NULL DEFAULT 'active',

    started_at DATETIME NULL,
    ended_at DATETIME NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    UNIQUE KEY uq_projects_company_slug (
        company_id,
        slug
    ),

    KEY idx_projects_company_id (company_id),
    KEY idx_projects_team_id (team_id),
    KEY idx_projects_status (status),

    CONSTRAINT fk_projects_company
        FOREIGN KEY (company_id)
        REFERENCES companies (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_projects_team
        FOREIGN KEY (team_id)
        REFERENCES teams (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    CONSTRAINT chk_projects_dates
        CHECK (
            ended_at IS NULL
            OR started_at IS NULL
            OR ended_at >= started_at
        )
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
-- +goose Down

DROP TABLE IF EXISTS projects;
