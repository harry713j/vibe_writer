package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/handler"
)

func UserProfileRoutes(h *handler.UserProfileHandler, auth func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Patch("/profile", h.HandleUpdateProfile)
		r.Patch("/avatar", h.HandleUpdateAvatar)
		r.Get("/me", h.HandleGetOwnDetails)
		r.Delete("/avatar", h.HandleRemoveAvatar)
		r.Get("/bookmarks", h.HandleGetBookmarks)
		r.Post("/{username}/follow", h.HandleCreateFollow)
		r.Delete("/{username}/follow", h.HandleRemoveFollow)
		r.Get("/{username}/followings", h.HandleFetchFollowings)
		r.Get("/{username}/followers", h.HandleFetchFollowers)
	})
	r.Get("/{username}", h.HandleGetUserDetails)
	r.Get("/{username}/blogs", h.HandleGetAllBlog)
	r.Get("/{username}/blogs/{slug}", h.HandleGetBlog)
	r.Get("/{username}/blogs/{slug}/comments", h.HandleGetAllComments)

	return r
}
