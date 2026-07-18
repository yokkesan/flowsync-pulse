-- +goose Up

CREATE TABLE office_areas (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    office_id BIGINT UNSIGNED NOT NULL,

    name VARCHAR(150) NOT NULL,
    slug VARCHAR(100) NOT NULL,

    area_type ENUM(
        'main',
        'meeting',
        'executive',
        'break',
        'development',
        'custom'
    ) NOT NULL DEFAULT 'custom',

    background_key VARCHAR(255) NOT NULL,

    map_width INT UNSIGNED NOT NULL DEFAULT 1600,
    map_height INT UNSIGNED NOT NULL DEFAULT 900,

    max_members INT UNSIGNED NULL,

    sort_order INT UNSIGNED NOT NULL DEFAULT 0,
    status ENUM('active', 'archived') NOT NULL DEFAULT 'active',

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    UNIQUE KEY uq_office_areas_office_slug (
        office_id,
        slug
    ),

    KEY idx_office_areas_office_id (office_id),
    KEY idx_office_areas_status (status),

    CONSTRAINT fk_office_areas_office
        FOREIGN KEY (office_id)
        REFERENCES offices (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
-- +goose Down

DROP TABLE IF EXISTS office_areas;
