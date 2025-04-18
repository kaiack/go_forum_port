package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/kaiack/goforum/internal/store"
	"github.com/kaiack/goforum/utils"
)

type application struct {
	config     config
	store      store.Storage
	tokenMaker utils.JWTMaker
	validator  *validator.Validate
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Set up CORS middleware
	r.Use(cors.Handler(cors.Options{
		// Allow all origins for demo purposes
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"X-Custom-Header"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.registerHandler)
		r.Post("/login", app.loginHandler)
	})

	r.Route("/user", func(r chi.Router) {
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Put("/", app.updateUserHandler)
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Get("/", app.getUserHandler)
		r.With(GetAdminMiddleWareFunc(&app.tokenMaker, &app.store)).Put("/admin", app.updateUserAdmin)
	})

	r.Route("/thread", func(r chi.Router) {
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Post("/", app.MakeThreadHandler)
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Get("/", app.GetThreadHandler)
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Put("/", app.EditThreadHandler)
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Delete("/", app.DeleteThreadHandler)
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Put("/like", app.LikeThreadHandler)
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Put("/watch", app.WatchThreadHandler)
	})
	r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Get("/threads", app.GetThreadsHandler)

	r.Route("/comment", func(r chi.Router) {
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Post("/", app.MakeCommentHandler)
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Put("/", app.EditCommentHandler)
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Delete("/", app.DeleteCommentHandler)
		r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Put("/like", app.LikeCommentHandler)
	})
	r.With(GetAuthMiddleWareFunc(&app.tokenMaker)).Get("/comments", app.GetCommentsHandler)

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server has started at %s", app.config.addr)
	return srv.ListenAndServe()
}

// 1:23
