package main

import (
	"database/sql"
	"log"

	"github.com/kaiack/goforum/internal/env"
	"github.com/kaiack/goforum/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

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

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}

// Read about Go repository pattern
