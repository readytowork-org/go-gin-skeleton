CREATE TABLE IF NOT EXISTS `user_profile` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `address` VARCHAR(100) NOT NULL,
    `contact` VARCHAR(16) NOT NULL,
    `user_id` INT NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NULL,
    `deleted_at` DATETIME NULL,
    PRIMARY KEY (id),
    CONSTRAINT `UQ_user_profile_contact` UNIQUE (`contact`),
    CONSTRAINT `user_profile_ibfk_user_id` FOREIGN KEY (`user_id`)
        REFERENCES `users`(`id`) ON DELETE CASCADE,
    CONSTRAINT `UQ_user_profile_user_id` UNIQUE (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;