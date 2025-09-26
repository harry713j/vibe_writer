package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/handler"
)

func CommentRoutes(h *handler.CommentHandler, auth func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/{username}/{slug}", h.HandleGetAllComments)

	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Post("/", h.HandleCreateComment)
		r.Delete("/{commentId}", h.HandleDeleteComment)
	})

	return r
}
