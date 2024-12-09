package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kaiack/goforum/utils"
)

type UserUpdateReq struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Image    string `json:"image"`
}

type UserUpdateRes struct{}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

	var u UserUpdateReq

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	fmt.Println(u)
	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	fmt.Println(claims.Email, claims.Id)
	// newUser := store.User{Name: "Kai", Email: "kai@example.com", Password: "goodpassword"}
	// err = app.store.Users.Create(r.Context(), &newUser)
	// if err != nil {
	// 	fmt.Println(err)†
	// }
	// fmt.Println(newThread.ID)

	// w.Write([]byte("All good"))
	w.WriteHeader(http.StatusOK)
}

// "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTEsImVtYWlsIjoia2FpQGV4YW1wbGUuY29tIiwiYWRtaW4iOnRydWUsInN1YiI6ImthaUBleGFtcGxlLmNvbSIsImlhdCI6MTczMzc0MjI2MiwianRpIjoiYjI0MGMwMzctYjQ1My00YTE5LWExMmMtYjc5NTI3OTQ2ZGI3In0.Zm82-tJrDegtPL7YByWX3aJg843MYv5YSsoGBTq11Zs"
