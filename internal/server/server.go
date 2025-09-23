package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/harry713j/vibe_writer/internal/app"
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

func NewServer(config *ServerConfig, app *app.App) *http.Server {
	router := chi.NewRouter()

	cors := cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value for preflight request
	})

	router.Use(cors)
	v1Router := route.RegisterRoutes(app)
	router.Mount("/api/v1", v1Router)

	server := &http.Server{
		Addr:              ":" + config.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second, // for slowloris attack
	}

	return server
}
