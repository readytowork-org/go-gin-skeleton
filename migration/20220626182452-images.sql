
-- +migrate Up
CREATE TABLE IF NOT EXISTS `images` (
  `id` BINARY(16) NOT NULL, 
  `name` VARCHAR(100) NOT NULL,
  `img` VARCHAR(500) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`)
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE IF EXISTS `images`;