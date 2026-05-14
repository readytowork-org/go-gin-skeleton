-- Create "users" table
CREATE TABLE `users`
(
    `id`         INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `full_name`  VARCHAR(45)  NOT NULL,
    `phone`      VARCHAR(15)  NOT NULL,
    `gender`     VARCHAR(15)  NOT NULL,
    `email`      VARCHAR(100) NOT NULL,
    `password`   VARCHAR(100) NOT NULL,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME     NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `UQ_user_email` (`email`),
    UNIQUE INDEX `UQ_user_phone` (`phone`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
