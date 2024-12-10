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
