package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Comment struct {
	ID              int64  `json:"id"`
	Content         string `json:"content"`
	CreatorId       int64  `json:"creatorId"` // Foreign Key
	ThreadId        int64  `json:"threadId"`  // Foreign Key
	ParentCommentId *int64 `json:"parentCommentId"`
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

func (s *CommentsStore) EditComment(ctx context.Context, commentId int64, content string) error {
	query := `UPDATE comments SET content = ?WHERE id = ?;`

	// Execute the query
	_, err := s.db.Exec(query, content, commentId)
	if err != nil {
		// Return the error if the update fails
		return fmt.Errorf("failed to update comment: %w", err)
	}

	// No error, return nil
	return nil
}

func (s *CommentsStore) DeleteComment(ctx context.Context, commentId int64) error {
	return nil
}

func (s *CommentsStore) LikeComment(ctx context.Context, commentId int64, userId int64) error {
	return nil
}

func (s *CommentsStore) GetComments(ctx context.Context, threadId int64) error {
	return nil
}

func (s *CommentsStore) CheckCommentValid(ctx context.Context, commentId int64) (bool, error) {
	// if a parentComment is provided, then it can be nil.
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM comments
		WHERE id = ?
	);
`

	// Execute the query
	var exists bool
	err := s.db.QueryRowContext(ctx, query, commentId).Scan(&exists)
	if err != nil {
		// Return any error that occurs during the query execution
		return false, fmt.Errorf("error checking comment ID: %w", err)
	}

	// If count is greater than 0, the comment ID exists
	return exists, nil
}

func (s *CommentsStore) CheckCommentCreator(ctx context.Context, commentId int64, userId int64) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM comments
		WHERE id = ? AND creator_id = ?
	);
	`

	// Execute the query
	var exists bool
	err := s.db.QueryRowContext(ctx, query, commentId, userId).Scan(&exists)
	if err != nil {
		// Return any error that occurs during the query execution
		return false, fmt.Errorf("error checking if user created comment: %w", err)
	}

	return exists, nil
}

/*
	SELECT t.id AS thread_id, t.title, t.content, c.id AS comment_id, c.content AS comment_content, c.created_at AS comment_created_at
	FROM threads t
	LEFT JOIN comments c ON c.thread_id = t.id
	WHERE t.id = ?; -- replace with thread id


*/
