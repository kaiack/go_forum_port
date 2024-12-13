package store

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Thread struct {
	ID        int64           `json:"id"`
	Content   string          `json:"content"`
	Title     string          `json:"title"`
	IsPublic  *bool           `json:"isPublic"`
	CreatorID int64           `json:"creatorId"`
	CreatedAt string          `json:"createdAt"`
	Lock      *bool           `json:"lock"`
	Likes     map[string]bool `json:"likes"`
	Watchees  map[string]bool `json:"watchees"`
}

type ThreadsStore struct {
	db *sql.DB
}

func (s *ThreadsStore) CreateThread(ctx context.Context, thread *Thread) error {
	query := `
		INSERT INTO threads (content, title, isPublic, creatorId)
		VALUES (?, ?, ?, ?) RETURNING id
	`

	err := s.db.QueryRowContext(ctx, query, thread.Content, thread.Title, *thread.IsPublic, thread.CreatorID).Scan(
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

func (s *ThreadsStore) UpdateThread(ctx context.Context, thread *Thread) error {
	query := "UPDATE threads SET "
	var args []interface{}
	// Empty interface denotes "any type" since every type implements the empty interface,
	// This slice can hold any type, allowing us to collect all the different types for the
	// user and pass them as args to the DB call.
	setClauses := []string{}

	if thread.Title != "" {
		setClauses = append(setClauses, "title = ?")
		args = append(args, thread.Title)
	}
	if thread.Content != "" {
		setClauses = append(setClauses, "content = ?")
		args = append(args, thread.Content)
	}
	if thread.IsPublic != nil {
		setClauses = append(setClauses, "isPublic = ?")
		args = append(args, *thread.IsPublic)
	}
	if thread.Lock != nil {
		setClauses = append(setClauses, "lock = ?")
		args = append(args, *thread.Lock)
	}

	if len(setClauses) == 0 {
		return nil // Nothing to update, do nothing. Not really an error.
	}

	query += strings.Join(setClauses, ", ")
	query += " WHERE id = ?"
	args = append(args, thread.ID)
	_, err := s.db.ExecContext(ctx, query, args...)
	fmt.Println(err)
	return err
}

func (s *ThreadsStore) DeleteThread(ctx context.Context, threadId int64) error {
	query := "DELETE FROM threads WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, threadId)
	return err
}

func (s *ThreadsStore) LikeThread(ctx context.Context, threadId int64, userId int64) error {
	query := "INSERT INTO likes (user_id, thread_id) VALUES (?, ?)"
	_, err := s.db.ExecContext(ctx, query, userId, threadId)
	return err
}

func (s *ThreadsStore) WatchThread(ctx context.Context, threadId int64, userId int64) error {
	query := "INSERT INTO watching (userId, threadId) VALUES (?, ?)"
	_, err := s.db.ExecContext(ctx, query, userId, threadId)
	return err
}

func (s *ThreadsStore) ValidateThreadId(ctx context.Context, id int64) error {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM threads
		WHERE id = ?
	);`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		// Return any error that occurs during the query execution
		return fmt.Errorf("error checking if thread exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("Thread not found: %v", err)
	}

	return nil
}

func (s *ThreadsStore) IsThreadLocked(ctx context.Context, id int64) (bool, error) {
	// Query to check if the given thread is locked
	query := "SELECT lock FROM threads WHERE id = ?"

	var lock bool
	err := s.db.QueryRowContext(ctx, query, id).Scan(&lock)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %v", err)
	}

	return lock, nil
}

func (s *ThreadsStore) IsThreadOwner(ctx context.Context, userId int64, threadId int64) (bool, error) {

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM threads
		WHERE id = ? AND creatorId = ?
	);`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, threadId, userId).Scan(&exists)
	if err != nil {
		// Return any error that occurs during the query execution
		return false, fmt.Errorf("error checking if user created thread: %w", err)
	}

	return exists, nil
}

func (s *ThreadsStore) IsThreadPublic(ctx context.Context, id int64) (bool, error) {
	// return public, nil
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM threads
		WHERE id = ? AND isPublic = true
	);`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		// Return any error that occurs during the query execution
		return false, fmt.Errorf("error checking if thread public: %w", err)
	}

	return exists, nil
}
