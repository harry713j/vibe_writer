package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/harry713j/vibe_writer/internal/middleware"
	"github.com/harry713j/vibe_writer/internal/model"
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
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid params")
		return
	}
	// extract query parameters
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

func (u *UserProfileHandler) HandleRemoveAvatar(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := u.profileService.RemoveAvatar(userId)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, "")
}

func (h *UserProfileHandler) HandleGetAllComments(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	slug := chi.URLParam(r, "slug")

	if username == "" || slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Blog parameter required")
		return
	}

	comments, err := h.profileService.GetAllCommentsOfBlog(username, slug)

	if err != nil {
		if errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string][]model.CommentWithStat{"comments": comments})
}

func (h *UserProfileHandler) HandleGetBookmarks(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogs, err := h.profileService.FetchBookmarks(userId)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, blogs)
}

/* follow */
func (h *UserProfileHandler) HandleCreateFollow(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	followingUsername := chi.URLParam(r, "username")
	if followingUsername == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid params")
		return
	}

	err := h.profileService.CreateFollow(userId, followingUsername)
	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrInvalidFollowingUser) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Successful following user"})
}

func (h *UserProfileHandler) HandleRemoveFollow(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	followingUsername := chi.URLParam(r, "username")
	if followingUsername == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid params")
		return
	}

	err := h.profileService.RemoveFollow(userId, followingUsername)
	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrInvalidFollowingUser) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Successful unfollow user"})
}

func (h *UserProfileHandler) HandleFetchFollowers(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	authorUsername := chi.URLParam(r, "username")
	if authorUsername == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid params")
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid query params value")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid query params value")
		return
	}

	followers, err := h.profileService.FetchAllFollower(userId, authorUsername, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrInvalidAuthor) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, followers)
}

func (h *UserProfileHandler) HandleFetchFollowings(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	authorUsername := chi.URLParam(r, "username")
	if authorUsername == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid params")
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid query params value")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid query params value")
		return
	}

	followings, err := h.profileService.FetchAllFollowing(userId, authorUsername, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrInvalidAuthor) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, followings)
}
