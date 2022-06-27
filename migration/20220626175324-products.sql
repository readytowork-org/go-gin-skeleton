
-- +migrate Up
CREATE TABLE IF NOT EXISTS `products` (
  `id` BINARY(16) NOT NULL, 
  `name` VARCHAR(100) NOT NULL,
  `price` VARCHAR(256) NOT NULL,
  `description` VARCHAR(500) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`)
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE IF EXISTS `products`;