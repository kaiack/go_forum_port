package main

import (
	"net/http"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	// newUser := store.User{Name: "Kai", Email: "kai@example.com", Password: "goodpassword"}
	// err = app.store.Users.Create(r.Context(), &newUser)
	// if err != nil {
	// 	fmt.Println(err)†
	// }
	// fmt.Println(newThread.ID)

	w.Write([]byte("All good"))
}
