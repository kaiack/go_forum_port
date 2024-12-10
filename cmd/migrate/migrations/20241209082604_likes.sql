-- +goose Up
-- +goose StatementBegin
CREATE TABLE likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    thread_id INTEGER DEFAULT NULL, -- NULL means like is for a comment
    comment_id INTEGER DEFAULT NULL, -- NULL means like is for a thread
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    CONSTRAINT unique_like UNIQUE (user_id, thread_id, comment_id) -- ensures a user can like a thread/comment only once
);
INSERT INTO likes (user_id, thread_id) VALUES (1, 2);
INSERT INTO likes (user_id, thread_id) VALUES (2, 1);
INSERT INTO likes (user_id, thread_id) VALUES (3, 2);
INSERT INTO likes (user_id, thread_id) VALUES (4, 3);
INSERT INTO likes (user_id, thread_id) VALUES (5, 4);
INSERT INTO likes (user_id, thread_id) VALUES (6, 5);
INSERT INTO likes (user_id, thread_id) VALUES (7, 10);
INSERT INTO likes (user_id, thread_id) VALUES (8, 10);
INSERT INTO likes (user_id, thread_id) VALUES (9, 10);
INSERT INTO likes (user_id, thread_id) VALUES (10, 9);
INSERT INTO likes (user_id, thread_id) VALUES (10, 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE likes
-- +goose StatementEnd
