package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kaiack/goforum/internal/store"
	"github.com/kaiack/goforum/utils"
)

type GetCommentsReq struct {
	ThreadId *int64 `json:"threadId" validate:"required,gte=1"`
}

type GetCommentsRes []store.Comment

type CreateCommentsReq struct {
	ThreadId        *int64 `json:"threadId" validate:"required,gte=1"`
	ParentCommentId *int64 `json:"parentCommentId" validate:"omitempty,gte=1"`
	Content         string `json:"content"`
}

type CreateCommentsRes struct {
	Id int64 `json:"id" validate:"required,gte=1"`
}

type EditCommentReq struct {
	Id      *int64 `json:"id" validate:"required,gte=1"`
	Content string `json:"content" validate:"required"`
}

type DeleteCommentReq struct {
	Id *int64 `json:"id" validate:"required,gte=1"`
}

type LikeCommentReq struct {
	Id     *int64 `json:"id" validate:"required,gte=1"`
	TurnOn bool   `json:"turnon" validate:"required"`
}

func (app *application) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// var c GetCommentsReq

	// if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
	// 	fmt.Println(err)
	// 	utils.SendError(w, "Request body invalid", http.StatusBadRequest)
	// 	return
	// }

	// err := app.validator.Struct(c)
	// if err != nil {
	// 	utils.HandleValidationError(err, w)
	// 	return
	// }

	// "126160": {
	// 	"id": 126160,
	// 	"creatorId": 35546,
	// 	"threadId": 808403,
	// 	"parentCommentId": null,
	// 	"content": "REAL NEW COMMENT!",
	// 	"createdAt": "2024-03-20T12:05:36.084Z",
	// 	"likes": {
	// 	  "35546": true
	// 	}
	//   }
	query := r.URL.Query()

	// Get the 'term' query parameter (required)
	theadIdString := query.Get("threadId")
	if theadIdString == "" {
		utils.SendError(w, "Missing required 'threadId' query parameter", http.StatusBadRequest)
		return
	}

	threadId, err := strconv.ParseInt(theadIdString, 10, 64)

	if err != nil || threadId < 0 {
		utils.SendError(w, "Invalid threadId ID", http.StatusBadRequest)
		return
	}

	err = app.store.Threads.ValidateThreadId(r.Context(), threadId)
	if err != nil {
		utils.SendError(w, "thread id not valid", http.StatusInternalServerError)
		return
	}

	comments, err := app.store.Comments.GetComments(r.Context(), threadId)

	if err != nil {
		utils.SendError(w, fmt.Sprintf("Error fetching threads: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(comments) == 0 {
		json.NewEncoder(w).Encode(make([]int64, 0))
	} else {
		json.NewEncoder(w).Encode(&comments)
	}
}

// As it stands, this lets you comment on locked threads and make a comment a child of a comment from another thread...
// will i fix it? nahhhhhhhhhhh cbs
func (app *application) MakeCommentHandler(w http.ResponseWriter, r *http.Request) {
	var c CreateCommentsReq

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		fmt.Println(err)
		utils.SendError(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(c)
	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	err = app.store.Threads.ValidateThreadId(r.Context(), *c.ThreadId)
	if err != nil {
		utils.SendError(w, "thread id not valid", http.StatusInternalServerError)
		return
	}

	err = app.store.Comments.CheckCommentValid(r.Context(), c.ParentCommentId, true)
	if err != nil {
		utils.SendError(w, "comment id not valid", http.StatusInternalServerError)
		return
	}

	newComment := store.Comment{
		ParentCommentId: c.ParentCommentId,
		Content:         c.Content,
		ThreadId:        *c.ThreadId,
		CreatorId:       userId,
	}
	err = app.store.Comments.Create(r.Context(), &newComment)

	if err != nil {
		utils.SendError(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	res := CreateCommentsRes{
		Id: newComment.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (app *application) EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	var c EditCommentReq

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		fmt.Println(err)
		utils.SendError(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(c)
	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	err = app.store.Comments.CheckCommentValid(r.Context(), c.Id, true)
	if err != nil {
		utils.SendError(w, "comment id not valid", http.StatusBadRequest)
		return
	}

	// Check user is creator of comment or admin to be allowed to change it.
	isCreator, err := app.store.Comments.CheckCommentCreator(r.Context(), *c.Id, userId)
	if err != nil {
		utils.SendError(w, "error fetching if comment creator is user", http.StatusInternalServerError)
		return
	}

	isAdmin, err := app.store.Users.IsUserAdmin(r.Context(), userId)
	if err != nil {
		utils.SendError(w, "error fetching if user admin", http.StatusInternalServerError)
		return
	}

	if !(isCreator || isAdmin) {
		utils.SendError(w, "Permisson Denied", http.StatusForbidden)
		return
	}

	app.store.Comments.EditComment(r.Context(), *c.Id, c.Content)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(make(map[string]string)) // send empty object?
}
func (app *application) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	var c DeleteCommentReq

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		fmt.Println(err)
		utils.SendError(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(c)
	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	err = app.store.Comments.CheckCommentValid(r.Context(), c.Id, true)
	if err != nil {
		utils.SendError(w, "comment id not valid", http.StatusBadRequest)
		return
	}

	// Check user is creator of comment or admin to be allowed to change it.
	isCreator, err := app.store.Comments.CheckCommentCreator(r.Context(), *c.Id, userId)
	if err != nil {
		utils.SendError(w, "error fetching if comment creator is user", http.StatusInternalServerError)
		return
	}

	isAdmin, err := app.store.Users.IsUserAdmin(r.Context(), userId)
	if err != nil {
		utils.SendError(w, "error fetching if user admin", http.StatusInternalServerError)
		return
	}

	if !(isCreator || isAdmin) {
		utils.SendError(w, "Permisson Denied", http.StatusForbidden)
		return
	}

	app.store.Comments.DeleteComment(r.Context(), *c.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(make(map[string]string)) // send empty object?
}

func (app *application) LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	var c LikeCommentReq

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		fmt.Println(err)
		utils.SendError(w, "Request body invalid", http.StatusBadRequest)
		return
	}

	err := app.validator.Struct(c)
	if err != nil {
		utils.HandleValidationError(err, w)
		return
	}

	claims := r.Context().Value(authKey{}).(*utils.UserClaims)
	userId := claims.Id

	err = app.store.Comments.CheckCommentValid(r.Context(), c.Id, true)
	if err != nil {
		utils.SendError(w, "comment id not valid", http.StatusBadRequest)
		return
	}

	threadId, err := app.store.Comments.GetThreadFromComment(r.Context(), *c.Id)
	if err != nil {
		utils.SendError(w, "error fetching thread for comment", http.StatusInternalServerError)
		return
	}
	// Check if user has permission to like this thread.
	isCreator, err := app.store.Threads.IsThreadOwner(r.Context(), userId, threadId)
	if err != nil {
		utils.SendError(w, "error fetching thread for creator", http.StatusInternalServerError)
		return
	}
	isAdmin, err := app.store.Users.IsUserAdmin(r.Context(), userId)
	if err != nil {
		utils.SendError(w, "error fetching if user admin", http.StatusInternalServerError)
		return
	}
	isPublic, err := app.store.Threads.IsThreadPublic(r.Context(), threadId)
	if err != nil {
		utils.SendError(w, "error fetching thread data", http.StatusInternalServerError)
		return
	}
	// fmt.Println(isPublic)
	if !(isPublic || isAdmin || isCreator) {
		utils.SendError(w, "Permisson Denied", http.StatusForbidden)
		return
	}

	app.store.Comments.LikeComment(r.Context(), *c.Id, userId, c.TurnOn)

	// check thread is public or user is admin or user is creator of the thread.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(make(map[string]string)) // send empty object?
}
