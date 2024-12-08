-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    admin BOOLEAN DEFAULT FALSE NOT NULL,
    image TEXT,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);
INSERT INTO users (name, password, email, admin) VALUES ('Alice', 'password123', 'alice@example.com', 0);
INSERT INTO users (name, password, email, admin) VALUES ('Bob', 'password456', 'bob@example.com', 1);
INSERT INTO users (name, password, email, admin) VALUES ('Charlie', 'password789', 'charlie@example.com', 0);
INSERT INTO users (name, password, email, admin) VALUES ('David', 'password000', 'david@example.com', 1);
INSERT INTO users (name, password, email, admin) VALUES ('Eve', 'mypassword1', 'eve@example.com', 0);
INSERT INTO users (name, password, email, admin) VALUES ('Frank', 'frankpassword', 'frank@example.com', 0);
INSERT INTO users (name, password, email, admin) VALUES ('Grace', 'grace1234', 'grace@example.com', 0);
INSERT INTO users (name, password, email, admin) VALUES ('Heidi', 'heidipassword', 'heidi@example.com', 0);
INSERT INTO users (name, password, email, admin) VALUES ('Ivan', 'ivansupersecure', 'ivan@example.com', 0);
INSERT INTO users (name, password, email, admin) VALUES ('Judy', 'judy123', 'judy@example.com', 0);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;
-- +goose StatementEnd