-- +goose Up
-- +goose StatementBegin
CREATE TABLE threads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT,
    title TEXT NOT NULL,
    isPublic BOOLEAN NOT NULL,
    creatorId INTEGER NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    lock BOOLEAN DEFAULT FALSE NOT NULL
);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('This is the first thread content.', 'First Thread', TRUE, 1, '2024-12-08 12:00:00', FALSE);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('Content for the second thread goes here.', 'Second Thread', FALSE, 2, '2024-12-08 14:30:00', TRUE);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('Here is the content of thread number three.', 'Third Thread', TRUE, 1, '2024-12-08 16:45:00', FALSE);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('A locked thread with secret content.', 'Locked Thread', FALSE, 3, '2024-12-08 18:00:00', TRUE);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('Content for the public thread with some updates.', 'Updated Public Thread', TRUE, 2, '2024-12-08 20:15:00', FALSE);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('This thread is private, and its content is hidden from others.', 'Private Thread', FALSE, 4, '2024-12-09 09:00:00', TRUE);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('Discussing tech advancements and AI in this thread.', 'Tech Thread', TRUE, 5, '2024-12-09 11:30:00', FALSE);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('Another locked thread with some ongoing discussion.', 'Another Locked Thread', FALSE, 6, '2024-12-09 13:45:00', TRUE);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('A public thread for sharing memes and jokes.', 'Meme Thread', TRUE, 7, '2024-12-09 15:30:00', FALSE);

INSERT INTO threads (content, title, isPublic, creatorId, createdAt, lock) 
VALUES ('This is a sensitive topic, hence locked.', 'Sensitive Discussion', FALSE, 8, '2024-12-09 17:00:00', TRUE);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE threads
-- +goose StatementEnd
