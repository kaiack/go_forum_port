package main

import (
	"database/sql"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/kaiack/goforum/internal/env"
	"github.com/kaiack/goforum/internal/store"
	"github.com/kaiack/goforum/utils"
	_ "github.com/mattn/go-sqlite3"
)

const minSecretLen = 32

func main() {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(db)

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 1),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 1),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", ""),
		},
	}

	store := store.NewStorage(db)

	var secret = env.GetString("JWT_SECRET", "pneumonoultramicroscopicsilicovolcanoconiosis")

	if len(secret) < minSecretLen {
		log.Fatalf("Need secret to be 32 chars in len or more")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	app := &application{
		config:     cfg,
		store:      store,
		tokenMaker: *utils.NewJWTMaker(secret),
		validator:  validate,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}

// Read about Go repository pattern
