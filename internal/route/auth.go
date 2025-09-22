package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/handler"
)

func AuthRoutes(h *handler.AuthHandler, auth func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", h.HandleSignup)
	r.Post("/login", h.HandleLogin)
	r.Post("/refresh", h.HandleRefreshAccessToken)

	r.Group(func(protected chi.Router) {
		protected.Use(auth)
		protected.Get("/logout", h.HandleLogout)
	})

	return r
}
