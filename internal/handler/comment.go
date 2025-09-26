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

type createCommentRequest struct {
	Slug     string `json:"slug"` // slug of the blog
	Content  string `json:"content"`
	ParentId int64  `json:"parent_id"` // could not be present
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		service: commentService,
	}
}

func (h *CommentHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	userid, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req createCommentRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Provide a valid json")
		return
	}

	comment, err := h.service.CreateComment(userid, req.Slug, req.ParentId, req.Content)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrBlogNotExists) ||
			errors.Is(err, service.ErrInvalidCommentContent) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, comment)
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

func (h *CommentHandler) HandleGetAllComments(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	username := chi.URLParam(r, "username")
	if slug == "" || username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Blog parameter required")
		return
	}

	comments, err := h.service.GetAllCommentsOfBlog(username, slug)

	if err != nil {
		if errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string][]*model.Comment{"comments": comments})
}
