package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/middleware"
	"github.com/harry713j/vibe_writer/internal/service"
	"github.com/harry713j/vibe_writer/internal/utils"
)

type UserProfileHandler struct {
	profileService *service.UserProfileService
	blogService    *service.BlogService
}

func NewUserProfileHandler(profileService *service.UserProfileService, blogService *service.BlogService) *UserProfileHandler {
	return &UserProfileHandler{
		profileService: profileService,
		blogService:    blogService,
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

// shift bottom two functions to user handler
func (u *UserProfileHandler) HandleGetAllBlog(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Inavlid params")
		return
	}
	// extracr query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid query value")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid query value")
		return
	}

	paginatedBlogRes, err := u.blogService.GetAllUserBlog(username, page, limit)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, paginatedBlogRes)
}

func (u *UserProfileHandler) HandleGetBlog(w http.ResponseWriter, r *http.Request) {

	// extract it from parameter
	slug := chi.URLParam(r, "slug")
	username := chi.URLParam(r, "username")
	if slug == "" || username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Blog parameter required")
		return
	}

	blogRes, err := u.blogService.GetBlog(username, slug)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, blogRes)
}
