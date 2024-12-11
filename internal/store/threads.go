package store

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

	// -------------------------------------------------------------------------------------------------------

	// To get likes, invert this for comments...
	query = `SELECT user_id FROM likes WHERE thread_id = ? AND comment_id IS NULL`

	likesRows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer likesRows.Close() // defer runs after this function returns.

	likesMap := make(map[string]bool)

	for likesRows.Next() {
		var userID int64
		if err := likesRows.Scan(&userID); err != nil {
			return nil, err
		}
		fmt.Println(userID)
		likesMap[strconv.FormatInt(userID, 10)] = true // Mark the user as having liked the thread
	}

	t.Likes = likesMap

	// -------------------------------------------------------------------------------------------------------

	query = `SELECT userId FROM watching WHERE threadId = ?`
	watchingRows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer watchingRows.Close() // defer runs after this function returns.
	watchingMap := make(map[string]bool)

	for watchingRows.Next() {
		var userID int64
		if err := watchingRows.Scan(&userID); err != nil {
			return nil, err
		}
		watchingMap[strconv.FormatInt(userID, 10)] = true // Mark the user as having liked the thread
	}

	t.Watchees = watchingMap

	return &t, nil
}

func (s *ThreadsStore) GetThreads(ctx context.Context, start int64, userId int64, isAdmin bool) ([]int64, error) {
	// Check If userId is admin
	// If so, get the nth-n+5th posts
	// Else get the n-n+5th posts where the post is public or the post is owned by userId

	var threadIds *sql.Rows
	if isAdmin {
		query := `SELECT id FROM threads ORDER BY id LIMIT 5 OFFSET ?;`
		threadIdsRows, err := s.db.QueryContext(ctx, query, start)
		if err != nil {
			return nil, err
		}
		threadIds = threadIdsRows
	} else {
		query := `SELECT id FROM threads WHERE (creatorId = ? OR isPublic = TRUE) ORDER BY id LIMIT 5 OFFSET ?;`
		threadIdsRows, err := s.db.QueryContext(ctx, query, userId, start)
		if err != nil {
			return nil, err
		}
		threadIds = threadIdsRows
	}

	defer threadIds.Close() // defer runs after this function returns.

	var idsList []int64

	for threadIds.Next() {
		var threadId int64
		if err := threadIds.Scan(&threadId); err != nil {
			return nil, err
		}
		fmt.Println(threadId)
		idsList = append(idsList, threadId)
	}

	return idsList, nil
}
