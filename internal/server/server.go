package server

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/route"
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
	v1Router := route.RegisterRoutes()
	router.Mount("/api/v1", v1Router)

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	return server
}
