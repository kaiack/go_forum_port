package store

import (
	"context"
	"database/sql"
)

type Thread struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	Title     string `json:"title"`
	IsPublic  bool   `json:"isPublic"`
	CreatorID int64  `json:"creatorId"`
	CreatedAt string `json:"createdAt"`
	Lock      bool   `json:"lock"`
}

type ThreadsStore struct {
	db *sql.DB
}

func (s *ThreadsStore) Create(ctx context.Context, thread *Thread) error {
	query := `
		INSERT INTO threads (content, title, isPublic, creatorId)
		VALUES (?, ?, ?, ?) RETURNING id
	`

	err := s.db.QueryRowContext(ctx, query, thread.Content, thread.Title, thread.IsPublic, thread.CreatorID).Scan(
		&thread.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// For comments later on
// https://stackoverflow.com/questions/55074867/posts-comments-replies-and-likes-database-schema/55075025
