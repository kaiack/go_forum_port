package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kaiack/goforum/internal/store"
	"github.com/kaiack/goforum/utils"
)

type UserUpdateReq struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Image    string `json:"image"`
}

type UserAdminUpdateReq struct {
	Id     int64 `json:"userId" validate:"required"`
	TurnOn *bool `json:"turnon" vaidate:"required"`
}

type UserGetReq struct {
	Id int64 `json:"userId" validate:"required"`
}

type UserGetRes struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Image    string `json:"image"`
	Admin    bool   `json:"admin"`
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

	var u UserUpdateReq

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	// fmt.Println(u)
	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	// fmt.Println(claims.Email, claims.Id)

	var updatedUser = store.User{ID: claims.Id, Email: u.Email, Name: u.Name, Password: u.Password, Image: u.Image}

	err := app.store.Users.UpdateUser(r.Context(), &updatedUser)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) updateUserAdmin(w http.ResponseWriter, r *http.Request) {
	var u UserAdminUpdateReq

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

	var updatedUser = store.User{ID: u.Id, Admin: u.TurnOn}

	err = app.store.Users.UpdateUser(r.Context(), &updatedUser)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	var u UserGetReq

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

	userInfo, err := app.store.Users.GetUserById(r.Context(), u.Id)

	if err != nil {
		http.Error(w, "error fetching user", http.StatusInternalServerError)
		return
	}

	var res UserGetRes = UserGetRes{
		Id:       userInfo.ID,
		Email:    userInfo.Email,
		Password: "", // Password returned from getUserById should also be empty, but just to be explicit here...
		Admin:    *userInfo.Admin,
		Image:    userInfo.Image,
		Name:     userInfo.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
