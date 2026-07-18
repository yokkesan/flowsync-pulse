-- +goose Up

CREATE TABLE extension_tokens (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    company_id BIGINT UNSIGNED NOT NULL,

    name VARCHAR(100) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,

    last_used_at DATETIME NULL,
    expires_at DATETIME NULL,
    revoked_at DATETIME NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    UNIQUE KEY uq_extension_tokens_token_hash (token_hash),

    KEY idx_extension_tokens_user_id (user_id),
    KEY idx_extension_tokens_company_id (company_id),
    KEY idx_extension_tokens_expires_at (expires_at),
    KEY idx_extension_tokens_revoked_at (revoked_at),

    CONSTRAINT fk_extension_tokens_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_extension_tokens_company
        FOREIGN KEY (company_id)
        REFERENCES companies (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
-- +goose Down

DROP TABLE IF EXISTS extension_tokens;
