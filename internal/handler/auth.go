package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/harry713j/vibe_writer/internal/middleware"
	"github.com/harry713j/vibe_writer/internal/service"
	"github.com/harry713j/vibe_writer/internal/utils"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Email, Username & password are required")
		return
	}

	user, err := h.service.RegisterUser(req.Username, req.Email, req.Password)

	if err != nil {
		if errors.Is(err, service.ErrUsernameExists) || errors.Is(err, service.ErrEmailExists) ||
			errors.Is(err, utils.ErrShortPassword) || errors.Is(err, utils.ErrInvalidEmail) || errors.Is(err, utils.ErrInvalidUsername) ||
			errors.Is(err, utils.ErrNoLowerCase) || errors.Is(err, utils.ErrNoNumber) || errors.Is(err, utils.ErrNoSpecialCharacter) ||
			errors.Is(err, utils.ErrNoUpperCase) || errors.Is(err, utils.ErrShortEmail) || errors.Is(err, utils.ErrShortUsername) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	res := registerResponse{
		ID:       user.Id.String(),
		Email:    user.Email,
		Username: user.Username,
	}

	utils.RespondWithJSON(w, http.StatusCreated, res)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Identifier == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Identifier & password are required")
		return
	}

	accessToken, refreshToken, err := h.service.LoginUser(req.Identifier, req.Password)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrWrongPassword) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// set the cookies
	rt := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		MaxAge:   7 * 86400,
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // true for https
		HttpOnly: true,
	}

	http.SetCookie(w, rt)

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message":      "Login successful",
		"access_token": accessToken,
	})
}

// must have authorization
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// auth middleware should inject user context into request
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized Request")
		return
	}

	err := h.service.LogoutUser(userId)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", MaxAge: -1, Path: "/"})

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Log out successful"})
}

// refresh access token
func (h *AuthHandler) HandleRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("refresh_token")

	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	newAccessToken, err := h.service.RefreshAccessToken(token.Value)

	if err != nil {

		if errors.Is(err, service.ErrExpiredRefreshToken) || errors.Is(err, service.ErrInvalidRefreshToken) ||
			errors.Is(err, service.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"access_token": newAccessToken,
	})
}
