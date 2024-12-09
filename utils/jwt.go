package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secret string
}

func NewJWTMaker(secret string) *JWTMaker {
	return &JWTMaker{secret}
}

func (maker *JWTMaker) CreateToken(id int64, email string, isAdmin bool) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, isAdmin)
	if err != nil {
		return "", nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(maker.secret))
	if err != nil {
		return "", nil, fmt.Errorf("errror signing token: %w", err)
	}

	return tokenStr, claims, nil
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		//verify signing method
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(maker.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
