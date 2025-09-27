package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/handler"
)

func LikeRoutes(h *handler.LikeHandler, auth func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(auth)

	r.Route("/comments", func(r chi.Router) {
		r.Post("/", h.HandleToggleCommentLike)
		r.Delete("/{commentId}", h.HandleRemoveCommentLike)
	})
	r.Route("/blogs", func(r chi.Router) {
		r.Post("/", h.HandleToggleBlogLike)
		r.Delete("/{slug}", h.HandleRemoveBlogLike)
	})

	return r
}
