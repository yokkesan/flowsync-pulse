-- +goose Up

CREATE TABLE user_presence (
    user_id BIGINT UNSIGNED NOT NULL,
    office_id BIGINT UNSIGNED NOT NULL,
    area_id BIGINT UNSIGNED NULL,

    status ENUM(
        'online',
        'working',
        'away',
        'offline'
    ) NOT NULL DEFAULT 'offline',

    position_x DECIMAL(10, 2) NOT NULL DEFAULT 0,
    position_y DECIMAL(10, 2) NOT NULL DEFAULT 0,

    last_seen_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id),

    KEY idx_user_presence_office_id (office_id),
    KEY idx_user_presence_area_id (area_id),
    KEY idx_user_presence_status (status),
    KEY idx_user_presence_last_seen_at (last_seen_at),

    CONSTRAINT fk_user_presence_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_user_presence_office
        FOREIGN KEY (office_id)
        REFERENCES offices (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_user_presence_area
        FOREIGN KEY (area_id)
        REFERENCES office_areas (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    CONSTRAINT chk_user_presence_position
        CHECK (
            position_x >= 0
            AND position_y >= 0
        )
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
-- +goose Down

DROP TABLE IF EXISTS user_presence;
