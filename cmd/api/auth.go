package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kaiack/goforum/internal/store"
	"github.com/kaiack/goforum/utils"
)

type UserRegisterReq struct {
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterRes struct {
	Token  string `json:"token"`
	UserId int64  `json:"userId"`
}

type UserLoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type userLoginRes struct {
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

	err := app.validator.Struct(u)
	if err != nil {
		fmt.Println(err)
		utils.HandleValidationError(err, w)
		return
	}

	// TODO: Check if this is the first user in the db or not.
	usersEmpty, err := app.store.Users.IsUsersEmpty(r.Context())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	hashed, err := utils.HashPassword(u.Password)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	u.Password = hashed
	// Empty users ->this is first user. Set them to admin
	setAdmin := usersEmpty

	newUser := store.User{Name: u.Name, Email: u.Email, Password: u.Password, Admin: &setAdmin}
	err = app.store.Users.Create(r.Context(), &newUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error creating new User", http.StatusInternalServerError)
		return
	}

	accessToken, _, err := app.tokenMaker.CreateToken(newUser.ID, newUser.Email, *newUser.Admin)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	fmt.Println(accessToken, newUser.ID)

	res := UserRegisterRes{
		Token:  accessToken,
		UserId: newUser.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var u UserLoginReq

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(u)
	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	fetchedUser, err := app.store.Users.GetUserByEmail(r.Context(), u.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "User doesn't exist", http.StatusBadRequest)
		return
	}

	err = utils.CheckPassword(u.Password, fetchedUser.Password)

	if err != nil {
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	accessToken, _, err := app.tokenMaker.CreateToken(fetchedUser.ID, fetchedUser.Email, *fetchedUser.Admin)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	res := userLoginRes{
		Token:  accessToken,
		UserId: fetchedUser.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}
