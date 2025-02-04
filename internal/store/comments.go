package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Comment struct {
	ID              int64   `json:"id"`
	Content         string  `json:"content"`
	CreatorId       int64   `json:"creatorId"` // Foreign Key
	ThreadId        int64   `json:"threadId"`  // Foreign Key
	ParentCommentId *int64  `json:"parentCommentId"`
	CreatedAt       string  `json:"createdAt"`
	Likes           []int64 `json:"likes"`
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
	query := `UPDATE comments SET content = ? WHERE id = ?;`

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
	query := "DELETE FROM comments WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, commentId)
	return err
}

func (s *CommentsStore) LikeComment(ctx context.Context, commentId int64, userId int64, turnOn bool) error {
	if turnOn {
		query := "INSERT INTO likes (user_id, comment_id, thread_id) VALUES (?, ?, ?)"
		_, err := s.db.ExecContext(ctx, query, userId, commentId, -1)
		return err
	} else {
		query := "DELETE FROM likes WHERE user_id = ? AND comment_id = ?"
		_, err := s.db.ExecContext(ctx, query, userId, commentId)
		return err
	}
}

func (s *CommentsStore) GetComments(ctx context.Context, threadId int64) ([]Comment, error) {
	commentsQuery := `SELECT id, creator_id, thread_id, parent_comment_id, content, created_at
	FROM comments WHERE thread_id = ?`
	rows, err := s.db.QueryContext(ctx, commentsQuery, threadId)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.CreatorId, &comment.ThreadId, &comment.ParentCommentId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		// Initialize the Likes map for the comment
		comment.Likes = make([]int64, 0)

		// Fetch likes for the current comment
		likesQuery := `SELECT user_id FROM likes WHERE comment_id = ? AND thread_id IS -1`
		likesRows, err := s.db.Query(likesQuery, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch likes for comment %d: %w", comment.ID, err)
		}
		defer likesRows.Close()

		// Populate the Likes map
		for likesRows.Next() {
			var userID int64
			err := likesRows.Scan(&userID)
			if err != nil {
				return nil, fmt.Errorf("failed to scan like for comment %d: %w", comment.ID, err)
			}
			comment.Likes = append(comment.Likes, userID)
		}

		// Check for errors while reading likes
		if err := likesRows.Err(); err != nil {
			return nil, fmt.Errorf("error while reading likes for comment %d: %w", comment.ID, err)
		}

		// Append the comment with likes
		comments = append(comments, comment)
	}

	// Check for errors while reading comments
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading comments: %w", err)
	}

	return comments, nil
}

func (s *CommentsStore) CheckCommentValid(ctx context.Context, commentId *int64, canBeNil bool) error {
	if commentId == nil && canBeNil {
		return nil
	}

	// if a parentComment is provided, then it can be nil.
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM comments
		WHERE id = ?
	);`

	// Execute the query
	var exists bool
	err := s.db.QueryRowContext(ctx, query, *commentId).Scan(&exists)
	if err != nil {
		// Return any error that occurs during the query execution
		return fmt.Errorf("error checking comment ID: %w", err)
	}

	// IF a parent comment's id is provided it can be null.
	if !exists {
		return fmt.Errorf("Comment Id Invalid: %w", err)
	}

	// If count is greater than 0, the comment ID exists
	return nil
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

func (s *CommentsStore) GetThreadFromComment(ctx context.Context, commentId int64) (int64, error) {
	query := `SELECT thread_id from comments where id = ?`
	var threadId int64
	err := s.db.QueryRowContext(ctx, query, commentId).Scan(&threadId)

	return threadId, err
}

// func (s *CommentsStore) GetCommentLikes(ctx context.Context, comment *Comment, commentId int64) error {
// 	query := `SELECT user_id FROM likes WHERE comment_id = ? AND thread_id IS NULL`

// 	likesRows, err := s.db.QueryContext(ctx, query, commentId)
// 	if err != nil {
// 		return err
// 	}
// 	defer likesRows.Close() // defer runs after this function returns.

// 	likesMap := make(map[int64]bool)

// 	for likesRows.Next() {
// 		var userID int64
// 		if err := likesRows.Scan(&userID); err != nil {
// 			return err
// 		}

// 		likesMap[userID] = true // Mark the user as having liked the thread
// 	}

// 	comment.Likes = likesMap

// 	return nil
// }

/*
	SELECT t.id AS thread_id, t.title, t.content, c.id AS comment_id, c.content AS comment_content, c.created_at AS comment_created_at
	FROM threads t
	LEFT JOIN comments c ON c.thread_id = t.id
	WHERE t.id = ?; -- replace with thread id


*/
