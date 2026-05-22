CREATE TABLE IF NOT EXISTS `users`
(
    `id`         INT UNSIGNED AUTO_INCREMENT NOT NULL,
    `full_name`  VARCHAR(45)                 NOT NULL,
    `phone`      VARCHAR(15)                 NOT NULL,
    `gender`     VARCHAR(15)                 NOT NULL,
    `email`      VARCHAR(100)                NOT NULL,
    `password`   VARCHAR(100)                NOT NULL,
    `created_at` DATETIME                    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME                    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME                    NULL,
    PRIMARY KEY (id),
    CONSTRAINT `UQ_user_email` UNIQUE (`email`),
    CONSTRAINT `UQ_user_phone` UNIQUE (`phone`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- Refresh tokens are stored as SHA-256 hashes so a DB leak doesn't yield
-- usable tokens. Rotation: each /auth/refresh issues a new pair and marks the
-- previous row revoked_at = NOW(). Logout marks revoked_at on a single row.
CREATE TABLE IF NOT EXISTS `refresh_tokens`
(
    `id`         BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
    `user_id`    INT UNSIGNED                   NOT NULL,
    `token_hash` VARCHAR(255)                   NOT NULL,
    `expires_at` DATETIME                       NOT NULL,
    `revoked_at` DATETIME                       NULL,
    `created_at` DATETIME                       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT `UQ_refresh_token_hash` UNIQUE (`token_hash`),
    INDEX `IDX_refresh_user` (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
