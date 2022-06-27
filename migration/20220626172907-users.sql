
-- +migrate Up
CREATE TABLE IF NOT EXISTS `users` (
  `id` BINARY(16) NOT NULL, 
  `full_name` VARCHAR(100) NOT NULL,
  `email` VARCHAR(256) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT email_unique UNIQUE(email)
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE IF EXISTS `users`;