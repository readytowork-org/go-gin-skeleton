ALTER TABLE `blog_categories`
    ADD COLUMN `created_at` DATETIME NOT NULL,
    ADD COLUMN `updated_at` DATETIME NULL,
    ADD COLUMN `deleted_at` DATETIME NULL;
