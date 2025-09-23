package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/handler"
)

func UserProfileRoute(h *handler.UserProfileHandler, auth func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(auth)

	r.Patch("/profile", h.HandleUpdateProfile)
	r.Patch("/avatar", h.HandleUpdateAvatar)
	r.Get("/me", h.HandleGetUserDetails)

	return r
}
