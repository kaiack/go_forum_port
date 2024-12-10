package store

import (
	"context"
	"database/sql"
)

type Thread struct {
	ID        int64           `json:"id"`
	Content   string          `json:"content"`
	Title     string          `json:"title"`
	IsPublic  bool            `json:"isPublic"`
	CreatorID int64           `json:"creatorId"`
	CreatedAt string          `json:"createdAt"`
	Lock      bool            `json:"lock"`
	Likes     map[string]bool `json:"likes"`
	Watchees  map[string]bool `json:"watchees"`
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

func (s *ThreadsStore) GetThread(ctx context.Context, id int64) (*Thread, error) {
	var t Thread
	query := `SELECT content, title, creatorId, isPublic, createdAt, lock FROM threads WHERE id = ?`

	err := s.db.QueryRowContext(ctx, query, id).Scan(&t.Content, &t.Title, &t.CreatorID, &t.IsPublic, &t.CreatedAt, &t.Lock)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
