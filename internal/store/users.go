package store

import (
	"context"
	"database/sql"
	"strings"
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
		INSERT INTO users(name, password, email, admin) VALUES(?, ?, ?, ?) RETURNING id
	`

	err := s.db.QueryRowContext(ctx, query, user.Name, user.Password, user.Email, user.Admin).Scan(&user.ID)

	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) GetUserById(ctx context.Context, id int64) (*User, error) {
	var u User
	query := `SELECT id, name, password, email, admin, image FROM users where id=$1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Password, &u.Email, &u.Admin, &u.Image)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UsersStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	query := `SELECT id, name, password, email, admin, image FROM users where email=?`
	err := s.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Name, &u.Password, &u.Email, &u.Admin, &u.Image)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// https://stackoverflow.com/a/70025947
func (s *UsersStore) UpdateUser(ctx context.Context, user *User) error {

	query := "UPDATE users SET "
	var args []interface{}
	// Empty interface denotes "any type" since every type implements the empty interface,
	// This slice can hold any type, allowing us to collect all the different types for the
	// user and pass them as args to the DB call.
	setClauses := []string{}

	if user.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, user.Name)
	}
	if user.Password != "" {
		setClauses = append(setClauses, "password = ?")
		args = append(args, user.Password)
	}
	if user.Email != "" {
		setClauses = append(setClauses, "email = ?")
		args = append(args, user.Email)
	}
	if user.Image.Valid && user.Image.String != "" {
		setClauses = append(setClauses, "image = ?")
		args = append(args, user.Image.String)
	}

	if len(setClauses) == 0 {
		return nil // Nothing to update, do nothing. Not really an error.
	}

	query += strings.Join(setClauses, ", ")
	query += " WHERE id = ?"
	args = append(args, user.ID)

	_, err := s.db.ExecContext(ctx, query, args...) //.Scan(&u.ID, &u.Name, &u.Password, &u.Email, &u.Admin, &u.Image)

	return err
}

func (s *UsersStore) IsUsersEmpty(ctx context.Context) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	err := s.db.QueryRowContext(ctx, query).Scan(&count)

	return count == 0, err
}
