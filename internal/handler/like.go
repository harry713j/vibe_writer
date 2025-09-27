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

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

type toggleBlogLikeRequest struct {
	Slug     string         `json:"slug"`
	LikeType model.LikeType `json:"like_type"`
}

type toggleCommentLikeRequest struct {
	CommentId int64          `json:"comment_id"`
	LikeType  model.LikeType `json:"like_type"`
}

func (h *LikeHandler) HandleToggleBlogLike(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req toggleBlogLikeRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	like, err := h.likeService.ToggleBlogLike(userId, req.Slug, req.LikeType)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, like)
}

func (h *LikeHandler) HandleToggleCommentLike(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req toggleCommentLikeRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	like, err := h.likeService.ToggleCommentLike(userId, req.CommentId, req.LikeType)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrCommentNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, like)
}

func (h *LikeHandler) HandleRemoveBlogLike(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")

	if slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid data provided")
		return
	}

	err := h.likeService.RemoveBlogLike(userId, slug)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, "")
}

func (h *LikeHandler) HandleRemoveCommentLike(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentIdStr := chi.URLParam(r, "commentId")

	commentId, err := strconv.Atoi(commentIdStr)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid data provided")
		return
	}

	if commentId == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid data provided")
		return
	}

	err = h.likeService.RemoveCommentLike(userId, int64(commentId))

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrCommentNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, "")
}
