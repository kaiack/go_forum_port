package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/kaiack/goforum/internal/store"
	"github.com/kaiack/goforum/utils"
)

type authKey struct{}

// Should this be a method of the application struct or be a method that is passed the tokenMaker??
func GetAuthMiddleWareFunc(tokenMaker *utils.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read authorization header
			// Verify the token
			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				http.Error(w, fmt.Sprintf("error verifying token %v", err), http.StatusUnauthorized)
				return
			}
			// pass the payload/claims down the context
			ctx := context.WithValue(r.Context(), authKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAdminMiddleWareFunc(tokenMaker *utils.JWTMaker, userStorage *store.UsersStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read authorization header
			// Verify the token
			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				http.Error(w, fmt.Sprintf("error verifying token %v", err), http.StatusUnauthorized)
				return
			}

			isAdmin, err := userStorage.IsUserAdmin(r.Context(), claims.Id)

			if err != nil {
				http.Error(w, fmt.Sprintf("user not found??, big error call someone %v", err), http.StatusInternalServerError)
			}

			if !isAdmin {
				http.Error(w, "USer is not admin", http.StatusUnauthorized)
				return
			}
			// pass the payload/claims down the context
			ctx := context.WithValue(r.Context(), authKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifyClaimsFromAuthHeader(r *http.Request, tokenMaker *utils.JWTMaker) (*utils.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header missing")
	}

	fields := strings.Fields(authHeader) // Format is Bearer {Token}
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("malformed authorization header")
	}

	token := fields[1]
	claims, err := tokenMaker.VerifyToken(token)

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
