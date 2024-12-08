package main

import (
	"fmt"
	"net/http"

	"github.com/kaiack/goforum/internal/store"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	// newThread := store.Thread{Content: "Hello, World!", Title: "NEW POST", IsPublic: true, CreatorID: 1}
	// err := app.store.Threads.Create(r.Context(), &newThread)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(newThread.ID)

	// newUser := store.User{Name: "Kai", Email: "kai@example.com", Password: "goodpassword"}
	// err = app.store.Users.Create(r.Context(), &newUser)
	// if err != nil {
	// 	fmt.Println(err)†
	// }
	// fmt.Println(newThread.ID)

	// err := app.store.Users.GetUser(r.Context(), 10)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("yay")
	updatedUser := store.User{ID: 10, Name: "Kai", Email: "Kai@example.com", Password: "helloworld"}
	err := app.store.Users.UpdateUser(r.Context(), &updatedUser)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("yay")

	w.Write([]byte("All good"))
}
