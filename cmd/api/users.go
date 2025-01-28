package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
		utils.SendError(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	// fmt.Println(u)
	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	// fmt.Println(claims.Email, claims.Id)

	var updatedUser = store.User{ID: claims.Id, Email: u.Email, Name: u.Name, Password: u.Password, Image: u.Image}

	err := app.store.Users.UpdateUser(r.Context(), &updatedUser)

	if err != nil {
		fmt.Println(err)
		utils.SendError(w, "error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(make(map[string]string)) // send empty object?
}

func (app *application) updateUserAdmin(w http.ResponseWriter, r *http.Request) {
	var u UserAdminUpdateReq

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		utils.SendError(w, "Request body invalid", http.StatusBadRequest)
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
		utils.SendError(w, "error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(make(map[string]string)) // send empty object?
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// var u UserGetReq

	// if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
	// 	fmt.Println(err)
	// 	utils.SendError(w, "Request body invalid", http.StatusBadRequest)
	// 	return
	// }

	// err := app.validator.Struct(u)
	// if err != nil {
	// 	utils.HandleValidationError(err, w)
	// 	return
	// }

	query := r.URL.Query()

	// Get the 'term' query parameter (required)
	userIdString := query.Get("userId")
	if userIdString == "" {
		utils.SendError(w, "Missing required 'userId' query parameter", http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseInt(userIdString, 10, 64)

	if err != nil || userId < 1 {
		utils.SendError(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	userInfo, err := app.store.Users.GetUserById(r.Context(), userId)

	if err != nil {
		utils.SendError(w, "error fetching user", http.StatusInternalServerError)
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
