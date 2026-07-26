-- +goose Up

ALTER TABLE work_sessions
    ADD COLUMN match_status ENUM(
        'matched',
        'ticket_not_found',
        'branch_not_matched',
        'ticket_branch_mismatch'
    ) NOT NULL DEFAULT 'ticket_not_found'
        AFTER ticket_key,

    ADD KEY idx_work_sessions_match_status (
        match_status
    );

-- +goose Down

ALTER TABLE work_sessions
    DROP INDEX idx_work_sessions_match_status,
    DROP COLUMN match_status;