package db

import (
	"database/sql"

	"github.com/harry713j/vibe_writer/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(config *config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.URL)

	if err != nil {
		return nil, err
	}
	// test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
