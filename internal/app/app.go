package app

import (
	"github.com/harry713j/vibe_writer/internal/handler"
	"github.com/harry713j/vibe_writer/internal/service"
)

type App struct {
	AuthService        *service.AuthService
	AuthHandler        *handler.AuthHandler
	UserProfileService *service.UserProfileService
	UserProfileHandler *handler.UserProfileHandler
	BlogService        *service.BlogService
	BlogHandler        *handler.BlogHandler
	CommentService     *service.CommentService
	CommentHandler     *handler.CommentHandler
}
