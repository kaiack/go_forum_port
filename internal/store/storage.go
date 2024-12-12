package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Threads interface {
		CreateThread(context.Context, *Thread) error
		GetThread(context.Context, int64) (*Thread, error)
		GetThreads(context.Context, int64, int64, bool) ([]int64, error)
		UpdateThread(ctx context.Context, thread *Thread) error
		DeleteThread(ctx context.Context, threadId int64) error
		LikeThread(ctx context.Context, threadId int64, userId int64) error
		ValidateThreadId(ctx context.Context, id int64) error
		IsThreadLocked(ctx context.Context, id int64) (bool, error)
		IsThreadOwner(ctx context.Context, userId int64, threadId int64) (bool, error)
		IsThreadPublic(ctx context.Context, id int64) (bool, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetUserById(context.Context, int64) (*User, error)
		GetUserByEmail(context.Context, string) (*User, error)
		UpdateUser(context.Context, *User) error
		IsUsersEmpty(context.Context) (bool, error)
		IsUserAdmin(context.Context, int64) (bool, error)
	}
	Comments interface {
		Create(context.Context, *Comment)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Threads: &ThreadsStore{db},
		Users:   &UsersStore{db},
	}
}
