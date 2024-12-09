package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kaiack/goforum/internal/store"
)

type UserRegisterReq struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserRegisterRes struct {
	Token  string `json:"token"`
	UserId int64  `json:"userId"`
}

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	var u UserRegisterReq

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	// TODO: Check if this is the first user in the db or not.
	usersEmpty, err := app.store.Users.IsUsersEmpty(r.Context())
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
	}

	newUser := store.User{Name: u.Name, Email: u.Email, Password: u.Password, Admin: !usersEmpty}
	err = app.store.Users.Create(r.Context(), &newUser)
	if err != nil {
		http.Error(w, "error creating new User", http.StatusInternalServerError)
		return
	}

	accessToken, _, err := app.tokenMaker.CreateToken(newUser.ID, newUser.Email, newUser.Admin, 15*time.Minute)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	res := UserRegisterRes{
		Token:  accessToken,
		UserId: newUser.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
