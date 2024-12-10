-- +goose Up
-- +goose StatementBegin
CREATE TABLE watching (
    userId INTEGER NOT NULL,                      -- 'userId' is the ID of the user watching the thread
    threadId INTEGER NOT NULL,                    -- 'threadId' is the ID of the thread being watched
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP, -- 'createdAt' stores the timestamp of when the user started watching the thread
    PRIMARY KEY (userId, threadId),               -- The combination of 'userId' and 'threadId' is unique (a user can only watch a thread once)
    FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE,   -- Foreign key reference to 'users' table
    FOREIGN KEY (threadId) REFERENCES threads(id) ON DELETE CASCADE -- Foreign key reference to 'threads' table
);
INSERT INTO watching (userId, threadId) VALUES (1, 2);
INSERT INTO watching (userId, threadId) VALUES (2, 1);
INSERT INTO watching (userId, threadId) VALUES (3, 2);
INSERT INTO watching (userId, threadId) VALUES (4, 3);
INSERT INTO watching (userId, threadId) VALUES (5, 4);
INSERT INTO watching (userId, threadId) VALUES (6, 5);
INSERT INTO watching (userId, threadId) VALUES (7, 6);
INSERT INTO watching (userId, threadId) VALUES (8, 10);
INSERT INTO watching (userId, threadId) VALUES (9, 10);
INSERT INTO watching (userId, threadId) VALUES (10, 9);
INSERT INTO watching (userId, threadId) VALUES (10, 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE watching
-- +goose StatementEnd
