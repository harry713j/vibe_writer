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

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		service: commentService,
	}
}

type toggleCommentLikeRequest struct {
	LikeType model.LikeType `json:"like_type"`
}

func (h *CommentHandler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	userid, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentIdStr := chi.URLParam(r, "commentId")
	commentId, err := strconv.Atoi(commentIdStr)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Inavlid comment id")
		return
	}

	err = h.service.DeleteComment(userid, int64(commentId))

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "comment deleted successfully"})
}

func (h *CommentHandler) HandleToggleCommentLike(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentId, err := strconv.ParseInt(chi.URLParam(r, "commentId"), 10, 64)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid param")
		return
	}

	var req toggleCommentLikeRequest
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	like, err := h.service.ToggleCommentLike(userId, commentId, req.LikeType)

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

func (h *CommentHandler) HandleRemoveCommentLike(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentId, err := strconv.ParseInt(chi.URLParam(r, "commentId"), 10, 64)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid param")
		return
	}

	if commentId == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid data provided")
		return
	}

	err = h.service.RemoveCommentLike(userId, int64(commentId))

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
