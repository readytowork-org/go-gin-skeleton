CREATE TABLE IF NOT EXISTS category (
  `id` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(100) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;


CREATE TABLE IF NOT EXISTS blog (
  `id` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(255) NOT NULL,
  `content` TEXT NOT NULL,
  `thumbnail_image` VARCHAR(255) NOT NULL,
  `is_published` BOOLEAN NOT NULL DEFAULT 0,
  `created_by` INT NOT NULL,
  `updated_by` INT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (id),
  CONSTRAINT `blog_ibfk_created_by` FOREIGN KEY (`created_by`)
    REFERENCES `user`(`id`) ON DELETE CASCADE,
  CONSTRAINT `blog_ibfk_updated_by` FOREIGN KEY (`updated_by`)
    REFERENCES `user`(`id`) ON DELETE CASCADE

) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS blog_categories (
    `id` INT NOT NULL AUTO_INCREMENT,
    `blog_id` INT NOT NULL,
    `category_id`INT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT `blog_categories_ibfk_blog_id` FOREIGN KEY (`blog_id`)
        REFERENCES `blog`(`id`) ON DELETE CASCADE,
    CONSTRAINT `blog_categories_ibfk_category_id` FOREIGN KEY (`category_id`)
        REFERENCES `category`(`id`) ON DELETE CASCADE

) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;