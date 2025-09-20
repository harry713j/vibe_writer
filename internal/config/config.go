package config

import (
	"os"
)

type DBConfig struct {
	URL string
}

func LoadDBConfig() *DBConfig {
	dbUrl := os.Getenv("DATABASE_URL")
	return &DBConfig{URL: dbUrl}
}
