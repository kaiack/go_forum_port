package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Threads interface {
		Create(context.Context, *Thread) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetUser(context.Context, int64) (*User, error)
		UpdateUser(context.Context, *User) error
		IsUsersEmpty(context.Context) (bool, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Threads: &ThreadsStore{db},
		Users:   &UsersStore{db},
	}
}
