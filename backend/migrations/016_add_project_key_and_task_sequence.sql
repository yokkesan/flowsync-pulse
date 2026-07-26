-- +goose Up

ALTER TABLE projects
    ADD COLUMN project_key VARCHAR(10) NULL
        AFTER slug,

    ADD COLUMN next_task_number BIGINT UNSIGNED NOT NULL DEFAULT 1
        AFTER project_key,

    ADD UNIQUE KEY uq_projects_company_project_key (
        company_id,
        project_key
    ),

    ADD CONSTRAINT chk_projects_project_key
        CHECK (
            project_key IS NULL
            OR project_key REGEXP '^[A-Z][A-Z0-9]{1,9}$'
        ),

    ADD CONSTRAINT chk_projects_next_task_number
        CHECK (
            next_task_number >= 1
        );

-- +goose Down

ALTER TABLE projects
    DROP CHECK chk_projects_next_task_number,
    DROP CHECK chk_projects_project_key,
    DROP INDEX uq_projects_company_project_key,
    DROP COLUMN next_task_number,
    DROP COLUMN project_key;