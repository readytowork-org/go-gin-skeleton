CREATE TABLE IF NOT EXISTS {{resourcetable}} (
    id INT NOT NULL AUTO_INCREMENT,
    title VARCHAR(20) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
