package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kaiack/goforum/internal/store"
	"github.com/kaiack/goforum/utils"
)

type MakeThreadReq struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	IsPublic bool   `json:"isPublic"`
}

type MakeThreadRes struct {
	Id int64 `json:"id"`
}

type GetThreadReq struct {
	Id int64 `json:"id"`
}

type GetThreadRes = store.Thread

func (app *application) MakeThreadHandler(w http.ResponseWriter, r *http.Request) {
	var t MakeThreadReq

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	newThread := store.Thread{
		Title:     t.Title,
		Content:   t.Content,
		IsPublic:  t.IsPublic,
		CreatorID: userId,
		Lock:      false,
	}

	err := app.store.Threads.Create(r.Context(), &newThread)

	if err != nil {
		// Will this cover the case in which creatorId is not present in the DB? -> Should be enforced by foreign key.1
		http.Error(w, fmt.Sprintf("Error creating thread: %s", err), http.StatusBadRequest)
		return
	}

	res := MakeThreadRes{
		Id: newThread.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (app *application) GetThreadHandler(w http.ResponseWriter, r *http.Request) {
	var t GetThreadReq

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	thread, err := app.store.Threads.GetThread(r.Context(), t.Id)
	thread.ID = t.Id

	if err != nil {
		http.Error(w, fmt.Sprintf("Error finding thread: %s", err), http.StatusBadRequest)
		return
	}

	// Make sure users gets an empty map in json rather than null.
	if thread.Likes == nil {
		thread.Likes = make(map[string]bool)
	}
	if thread.Watchees == nil {
		thread.Watchees = make(map[string]bool)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&thread)
}
