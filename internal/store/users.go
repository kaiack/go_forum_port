package store

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Password  string         `json:"password"`
	Email     string         `json:"email"`
	Image     sql.NullString `json:"image"`
	Admin     bool           `json:"admin"`
	CreatedAt string         `json:"createdAt"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users(name, password, email) VALUES($1, $2, $3) RETURNING id
	`

	err := s.db.QueryRowContext(ctx, query, user.Name, user.Password, user.Email).Scan(&user.ID)

	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) GetUser(ctx context.Context, id int64) error {
	var u User
	query := `SELECT id, name, password, email, admin, image FROM users where id=$1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Password, &u.Email, &u.Admin, &u.Image)
	if err != nil {
		return err
	}
	fmt.Println(u)

	return nil
}
