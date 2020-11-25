-- name: create-users-table
CREATE TABLE `users` (
        `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `firstname` VARCHAR(64) NULL,
        `lastname` VARCHAR(64) NULL,
        `email` VARCHAR(64) NULL,
        `website` VARCHAR(32) NULL,
        `createdAt` DATE NULL,
        `updatedAt` DATE NULL
);

-- name: create-user
INSERT INTO users (firstname, lastname, email, website, createdAt, updatedAt ) VALUES (?, ?, ?, ?, ?, ?)