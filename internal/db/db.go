package db

import (
	"database/sql"
	"log"

	"github.com/harry713j/vibe_writer/internal/config"
)

func ConnectDB(config *config.DBConfig) *sql.DB {
	db, err := sql.Open("postgres", config.URL)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
