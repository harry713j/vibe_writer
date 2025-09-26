package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/service"
	"github.com/harry713j/vibe_writer/internal/utils"
)

type contextKey string

const userIdKey contextKey = "userID"

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract the access token
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "No authorization header found")
				return
			}

			authValues := strings.Split(authHeader, " ")

			if len(authValues) != 2 || authValues[0] != "Bearer" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization header")
				return
			}
			token := authValues[1]
			// validate the token
			claims, err := authService.ValidateJwtToken(token)

			if err != nil {
				log.Println("JWT validation error ", err)

				if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrExpiredToken) {
					utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
					return
				}

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
			ctx := context.WithValue(r.Context(), userIdKey, userId)
			// call next handler with the new request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the user ID from the request context
func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(userIdKey).(uuid.UUID)
	return userID, ok
}
