package route

import (
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/app"
	"github.com/harry713j/vibe_writer/internal/handler"
	"github.com/harry713j/vibe_writer/internal/middleware"
)

func RegisterRoutes(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/health", handler.HandleHealth)

	r.Mount("/auth", AuthRoutes(app.AuthHandler, middleware.AuthMiddleware(app.AuthService)))

	return r
}
