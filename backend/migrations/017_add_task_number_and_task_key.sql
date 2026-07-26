-- +goose Up

ALTER TABLE tasks
    ADD COLUMN task_number BIGINT UNSIGNED NULL
        AFTER project_id,

    ADD COLUMN task_key VARCHAR(32) NULL
        AFTER task_number;

-- 既存タスクへプロジェクト単位の連番を割り当てる。
-- created_at が同じ場合は id の昇順で採番する。
UPDATE tasks AS target_task
INNER JOIN (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY project_id
            ORDER BY created_at ASC, id ASC
        ) AS generated_task_number
    FROM tasks
) AS numbered_task
    ON numbered_task.id = target_task.id
SET target_task.task_number = numbered_task.generated_task_number;

-- project_key が設定済みの場合のみtask_keyを生成する。
UPDATE tasks AS task
INNER JOIN projects AS project
    ON project.id = task.project_id
SET task.task_key = CONCAT(
    project.project_key,
    '-',
    task.task_number
)
WHERE project.project_key IS NOT NULL;

-- 各プロジェクトの次回採番値を更新する。
UPDATE projects AS project
LEFT JOIN (
    SELECT
        project_id,
        MAX(task_number) AS max_task_number
    FROM tasks
    GROUP BY project_id
) AS task_summary
    ON task_summary.project_id = project.id
SET project.next_task_number = COALESCE(
    task_summary.max_task_number + 1,
    1
);

ALTER TABLE tasks
    MODIFY COLUMN task_number BIGINT UNSIGNED NOT NULL,

    ADD UNIQUE KEY uq_tasks_project_task_number (
        project_id,
        task_number
    ),

    ADD UNIQUE KEY uq_tasks_task_key (
        task_key
    ),

    ADD CONSTRAINT chk_tasks_task_number
        CHECK (
            task_number >= 1
        ),

    ADD CONSTRAINT chk_tasks_task_key
        CHECK (
            task_key IS NULL
            OR task_key REGEXP '^[A-Z][A-Z0-9]{1,9}-[1-9][0-9]*$'
        );

-- +goose Down

ALTER TABLE tasks
    DROP CHECK chk_tasks_task_key,
    DROP CHECK chk_tasks_task_number,
    DROP INDEX uq_tasks_task_key,
    DROP INDEX uq_tasks_project_task_number,
    DROP COLUMN task_key,
    DROP COLUMN task_number;