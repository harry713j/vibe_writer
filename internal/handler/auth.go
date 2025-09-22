package handler

import (
	"encoding/json"
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
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
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
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// set the cookies
	rt := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		MaxAge:   7 * 86400,
		Secure:   false, // true for https
		HttpOnly: true,
	}

	at := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		MaxAge:   900,
		Secure:   false, // true for https
		HttpOnly: true,
	}

	http.SetCookie(w, rt)
	http.SetCookie(w, at)

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Login successful",
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
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// remove the access and refresh token cookies
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "", MaxAge: -1, Path: "/"})
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
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	at := &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Path:     "/",
		MaxAge:   900,
		Expires:  time.Now().Add(15 * time.Minute),
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, at)
	utils.RespondWithJSON(w, http.StatusNoContent, "")
}
