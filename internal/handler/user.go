package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/middleware"
	"github.com/harry713j/vibe_writer/internal/service"
	"github.com/harry713j/vibe_writer/internal/utils"
)

type UserProfileHandler struct {
	profileService *service.UserProfileService
}

func NewUserProfileHandler(service *service.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		profileService: service,
	}
}

type updateProfileRequest struct {
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
}

type updateAvatarRequest struct {
	AvatarUrl string `json:"avatar_url"`
}

// update profile
func (u *UserProfileHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req updateProfileRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	userDetails, err := u.profileService.UpdateUserProfile(userId, req.FullName, req.Bio)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, userDetails)
}

// update avatar
func (u *UserProfileHandler) HandleUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var request updateAvatarRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Avatar url is required")
		return
	}

	userDetails, err := u.profileService.UpdateAvatar(userId, request.AvatarUrl)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, userDetails)
}

// get user details
func (u *UserProfileHandler) HandleGetOwnDetails(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userDetails, err := u.profileService.GetProfileDetails(userId)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, userDetails)
}

func (u *UserProfileHandler) HandleGetUserDetails(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	username := chi.URLParam(r, "username")

	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid username")
		return
	}

	userDetails, err := u.profileService.GetUserDetails(username)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, userDetails)
}
