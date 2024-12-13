package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kaiack/goforum/internal/store"
	"github.com/kaiack/goforum/utils"
)

type MakeThreadReq struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content"`
	IsPublic bool   `json:"isPublic" validate:"required"`
}

type MakeThreadRes struct {
	Id int64 `json:"id"`
}

type GetThreadReq struct {
	Id int64 `json:"id" validate:"required"`
}

type GetThreadRes = store.Thread

type GetThreadsReq struct {
	Start int64 `json:"count" validate:"required"`
}

type GetThreadsRes []int64

type EditThreadReq struct {
	Id       *int64 `json:"id" validate:"required,gte=1"` // Is a pointer because default value is 0, which we dont really want and isnt a valid id in SQL.
	Title    string `json:"title"`
	IsPublic *bool  `json:"isPublic"`
	Content  string `json:"content"`
	Lock     *bool  `json:"lock"`
}

type DeleteThreadReq struct {
	Id *int64 `json:"id" validate:"required,gte=1"`
}

type LikeThreadReq struct {
	Id     *int64 `json:"id" validate:"required,gte=1"`
	TurnOn *bool  `json:"turnOn" validate:"required"`
}

type WatchThreadReq struct {
	Id     *int64 `json:"id" validate:"required,gte=1"`
	TurnOn *bool  `json:"turnOn" validate:"required"`
}

// id, title, isPublic, content, lock
func (app *application) MakeThreadHandler(w http.ResponseWriter, r *http.Request) {
	var t MakeThreadReq

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(t)

	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id
	lock := false
	newThread := store.Thread{
		Title:     t.Title,
		Content:   t.Content,
		IsPublic:  &t.IsPublic,
		CreatorID: userId,
		Lock:      &lock,
	}

	err = app.store.Threads.CreateThread(r.Context(), &newThread)

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

	err := app.validator.Struct(t)
	if err != nil {
		utils.HandleValidationError(err, w)
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

func (app *application) GetThreadsHandler(w http.ResponseWriter, r *http.Request) {
	var t GetThreadsReq

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(t)
	if err != nil {
		utils.HandleValidationError(err, w)
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	isAdmin, err := app.store.Users.IsUserAdmin(r.Context(), userId)
	if err != nil {
		http.Error(w, "Error fetching user details", http.StatusInternalServerError)
		return
	}

	threads, err := app.store.Threads.GetThreads(r.Context(), t.Start, userId, isAdmin)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching threads: %s", err), http.StatusInternalServerError)
		return
	}

	res := GetThreadsRes(threads)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)
}

func (app *application) EditThreadHandler(w http.ResponseWriter, r *http.Request) {
	var t EditThreadReq

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(t)

	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	// await assertValidThread(id);
	err = app.store.Threads.ValidateThreadId(r.Context(), *t.Id)
	if err != nil {
		http.Error(w, "thread id not valid", http.StatusInternalServerError)
		return
	}
	// await assertUnlockedThread(id);
	locked, err := app.store.Threads.IsThreadLocked(r.Context(), *t.Id)

	if err != nil {
		http.Error(w, "error fetching thread for lock", http.StatusInternalServerError)
		return
	}

	if locked {
		http.Error(w, "thread is locked", http.StatusInternalServerError)
		return
	}
	// await assertEditPermissionOfThread(authUserId, id);
	isCreator, err := app.store.Threads.IsThreadOwner(r.Context(), userId, *t.Id)
	if err != nil {
		http.Error(w, "error fetching thread for creator", http.StatusInternalServerError)
		return
	}
	isAdmin, err := app.store.Users.IsUserAdmin(r.Context(), userId)
	if err != nil {
		http.Error(w, "error fetching if user admin", http.StatusInternalServerError)
		return
	}

	if !(isCreator || isAdmin) {
		http.Error(w, "Permisson Denied", http.StatusForbidden)
		return
	}

	newThread := store.Thread{
		ID:       *t.Id,
		Content:  t.Content,
		Title:    t.Title,
		IsPublic: t.IsPublic,
		Lock:     t.Lock,
	}

	err = app.store.Threads.UpdateThread(r.Context(), &newThread)

	if err != nil {
		http.Error(w, "Error updating thread", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(make(map[string]string)) // send empty object?
}

func (app *application) DeleteThreadHandler(w http.ResponseWriter, r *http.Request) {
	var t DeleteThreadReq

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(t)

	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	// Check given threadID present in db
	err = app.store.Threads.ValidateThreadId(r.Context(), *t.Id)
	if err != nil {
		http.Error(w, "thread id not valid", http.StatusInternalServerError)
		return
	}

	// Check if user has permission to delete this thread.
	isCreator, err := app.store.Threads.IsThreadOwner(r.Context(), userId, *t.Id)
	if err != nil {
		http.Error(w, "error fetching thread for creator", http.StatusInternalServerError)
		return
	}
	isAdmin, err := app.store.Users.IsUserAdmin(r.Context(), userId)
	if err != nil {
		http.Error(w, "error fetching if user admin", http.StatusInternalServerError)
		return
	}

	if !(isCreator || isAdmin) {
		http.Error(w, "Permisson Denied", http.StatusForbidden)
		return
	}

	app.store.Threads.DeleteThread(r.Context(), *t.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(make(map[string]string)) // send empty object?
}

func (app *application) LikeThreadHandler(w http.ResponseWriter, r *http.Request) {
	var t LikeThreadReq

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(t)

	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	// Check given threadID present in db
	err = app.store.Threads.ValidateThreadId(r.Context(), *t.Id)
	if err != nil {
		http.Error(w, "thread id not valid", http.StatusInternalServerError)
		return
	}

	// Check if user has permission to like this thread.
	isCreator, err := app.store.Threads.IsThreadOwner(r.Context(), userId, *t.Id)
	if err != nil {
		http.Error(w, "error fetching thread for creator", http.StatusInternalServerError)
		return
	}
	isAdmin, err := app.store.Users.IsUserAdmin(r.Context(), userId)
	if err != nil {
		http.Error(w, "error fetching if user admin", http.StatusInternalServerError)
		return
	}
	isPublic, err := app.store.Threads.IsThreadPublic(r.Context(), *t.Id)
	if err != nil {
		http.Error(w, "error fetching thread data", http.StatusInternalServerError)
		return
	}
	fmt.Println(isPublic)
	if !isPublic && !(isAdmin || isCreator) {
		http.Error(w, "Permisson Denied", http.StatusForbidden)
		return
	}

	app.store.Threads.LikeThread(r.Context(), *t.Id, userId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(make(map[string]string)) // send empty object?
}

func (app *application) WatchThreadHandler(w http.ResponseWriter, r *http.Request) {
	var t WatchThreadReq

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		http.Error(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(t)

	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	// Check given threadID present in db
	err = app.store.Threads.ValidateThreadId(r.Context(), *t.Id)
	if err != nil {
		http.Error(w, "thread id not valid", http.StatusInternalServerError)
		return
	}

	// Check if user has permission to like this thread.
	isCreator, err := app.store.Threads.IsThreadOwner(r.Context(), userId, *t.Id)
	if err != nil {
		http.Error(w, "error fetching thread for creator", http.StatusInternalServerError)
		return
	}
	isAdmin, err := app.store.Users.IsUserAdmin(r.Context(), userId)
	if err != nil {
		http.Error(w, "error fetching if user admin", http.StatusInternalServerError)
		return
	}
	isPublic, err := app.store.Threads.IsThreadPublic(r.Context(), *t.Id)
	if err != nil {
		http.Error(w, "error fetching thread data", http.StatusInternalServerError)
		return
	}

	if !isPublic && !(isAdmin || isCreator) {
		http.Error(w, "Permisson Denied", http.StatusForbidden)
		return
	}

	app.store.Threads.WatchThread(r.Context(), *t.Id, userId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(make(map[string]string)) // send empty object?
}

// TODO:
// Add validation for route bodies where required...
