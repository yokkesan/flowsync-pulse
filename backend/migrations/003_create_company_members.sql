-- +goose Up

CREATE TABLE company_members (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    company_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,

    role ENUM('owner', 'admin', 'member') NOT NULL DEFAULT 'member',
    status ENUM('active', 'invited', 'suspended') NOT NULL DEFAULT 'active',
    joined_at DATETIME NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    UNIQUE KEY uq_company_members_company_user (
        company_id,
        user_id
    ),

    KEY idx_company_members_user_id (user_id),

    CONSTRAINT fk_company_members_company
        FOREIGN KEY (company_id)
        REFERENCES companies (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_company_members_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
-- +goose Down

DROP TABLE IF EXISTS company_members;
