package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID              int64  `json:"id"`
	Content         string `json:"content"`
	CreatorId       int64  `json:"creatorId"` // Foreign Key
	ThreadId        int64  `json:"threadId"`  // Foreign Key
	ParentCommentId int64  `json:"parentCommentId"`
	CreatedAt       string `json:"createdAt"`
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (content, creator_id, thread_id, parent_comment_id)
		VALUES (?, ?, ?, ?) RETURNING id
	`

	err := s.db.QueryRowContext(ctx, query, comment.Content, comment.CreatorId, comment.ThreadId, comment.ParentCommentId).Scan(
		&comment.ID,
	)

	return err
}

/*
	SELECT t.id AS thread_id, t.title, t.content, c.id AS comment_id, c.content AS comment_content, c.created_at AS comment_created_at
	FROM threads t
	LEFT JOIN comments c ON c.thread_id = t.id
	WHERE t.id = ?; -- replace with thread id


*/
