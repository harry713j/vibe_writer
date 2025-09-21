package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/service"
	"github.com/harry713j/vibe_writer/internal/utils"
)

type contextKey string

const UserIdKey contextKey = "userID"

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract the access token
			token, err := r.Cookie("access_token")

			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "No access token found in cookie")
				return
			}
			// validate the token
			claims, err := authService.ValidateJwtToken(token.Value)

			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			userIdStr, ok := claims["sub"].(string)

			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			userId, err := uuid.Parse(userIdStr)

			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID in token")
				return
			}
			// add userId to request context
			ctx := context.WithValue(r.Context(), UserIdKey, userId)
			// call next handler with the new request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the user ID from the request context
func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIdKey).(uuid.UUID)
	return userID, ok
}
