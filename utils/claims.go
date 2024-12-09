package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func NewUserClaims(id int64, email string, admin bool) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error during token id generation %w", err)
	}

	return &UserClaims{
		Email: email,
		ID:    id,
		Admin: admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       tokenID.String(),
			Subject:  email,
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil
}
