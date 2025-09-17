package server

import (
	"log"
	"net/http"
	"os"

	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type ServerConfig struct {
	Port string
}

func LoadServerConfig() *ServerConfig {
	port := os.Getenv("PORT")

	if port == "" {
		log.Println("Please Add Port value to .env file or environment")
		os.Exit(1)
	}

	return &ServerConfig{Port: port}
}

func NewServer(config *ServerConfig) *http.Server {
	router := chi.NewRouter()

	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)

	server := &http.Server{
		Addr:    config.Port,
		Handler: router,
	}

	return server
}
