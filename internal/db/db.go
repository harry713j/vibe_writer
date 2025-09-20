package db

import (
	"database/sql"
	"log"

	"github.com/harry713j/vibe_writer/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(config *config.DBConfig) *sql.DB {
	db, err := sql.Open("pgx", config.URL)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
