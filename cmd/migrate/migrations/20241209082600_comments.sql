-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    creator_id INT NOT NULL,
    thread_id INT NOT NULL,
    parent_comment_id INT DEFAULT NULL, -- for nested comments
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES comments(id) ON DELETE CASCADE -- allows for nested comments
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
