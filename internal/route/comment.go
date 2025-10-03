package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/handler"
)

func CommentRoutes(h *handler.CommentHandler, auth func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Delete("/{commentId}", h.HandleDeleteComment)
		r.Post("/{commentId}/reactions", h.HandleToggleCommentLike)
		r.Delete("/{comment}/reactions", h.HandleRemoveCommentLike)
	})

	return r
}
