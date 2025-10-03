package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/handler"
)

func BlogRoutes(h *handler.BlogHandler, auth func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(auth)

		r.Post("/", h.HandleCreateBlog)
		r.Put("/{slug}", h.HandleUpdateBlog)
		r.Delete("/{slug}", h.HandleDeleteBlog)
	})

	return r
}
