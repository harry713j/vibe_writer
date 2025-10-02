package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/handler"
)

func UploadRoutes(h *handler.UploadHandler, auth func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(auth)

	r.Post("/avatar", h.HandleUploadAvatar)
	r.Post("/blog", h.HandleUploadBlogImage)

	return r
}
