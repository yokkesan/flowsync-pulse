-- +goose Up

CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    display_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    avatar_type ENUM('preset', 'image') NOT NULL DEFAULT 'preset',
    avatar_key VARCHAR(100) NULL,
    avatar_image_url VARCHAR(500) NULL,

    status ENUM('active', 'suspended') NOT NULL DEFAULT 'active',
    last_login_at DATETIME NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE KEY uq_users_email (email)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
-- +goose Down

DROP TABLE IF EXISTS users;
