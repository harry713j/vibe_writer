package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/harry713j/vibe_writer/internal/middleware"
	"github.com/harry713j/vibe_writer/internal/model"
	"github.com/harry713j/vibe_writer/internal/service"
	"github.com/harry713j/vibe_writer/internal/utils"
)

type BlogHandler struct {
	blogService *service.BlogService
}

func NewBlogHandler(service *service.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: service,
	}
}

type createBlogRequest struct {
	Title     string   `json:"title"`
	Slug      string   `json:"slug"`
	Content   string   `json:"content"`
	PhotoUrls []string `json:"photo_urls"`
}

type updateBlogRequest struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	PhotoUrls []string `json:"photo_urls"`
}

type createCommentRequest struct {
	Content  string `json:"content"`
	ParentId int64  `json:"parent_id"` // could not be present
}

type toggleBlogLikeRequest struct {
	LikeType model.LikeType `json:"like_type"`
}

func (h *BlogHandler) HandleCreateBlog(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req createBlogRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Provide required data")
		return
	}

	if req.Title == "" || req.Slug == "" || req.Content == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Title, slug and content are required")
		return
	}

	blogData, err := h.blogService.CreateBlog(userId, req.Title, req.Slug, req.Content, req.PhotoUrls)

	if err != nil {
		if errors.Is(err, service.ErrTitleExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, blogData)
}

func (h *BlogHandler) HandleUpdateBlog(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req updateBlogRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	slug := chi.URLParam(r, "slug")

	if err != nil || slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Please Provide required data")
		return
	}

	updatedBlogData, err := h.blogService.UpdateBlog(userId, slug, req.Title, req.Content, req.PhotoUrls)

	if err != nil {
		if errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedBlogData)
}

func (h *BlogHandler) HandleDeleteBlog(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")

	if slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Please provide required data")
		return
	}

	err := h.blogService.DeleteBlog(userId, slug)

	if err != nil {
		if errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, "")
}

func (h *BlogHandler) HandleGetBlogs(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
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

	blogs, err := h.blogService.GetAllBlog(userId, page, limit)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, blogs)
}

func (h *BlogHandler) HandleChangeBlogVisibility(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")

	if slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "empty params")
		return
	}

	blog, err := h.blogService.ChangeBlogVisibility(userId, slug)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, blog)
}

func (h *BlogHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	userid, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")

	if slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid params")
		return
	}

	var req createCommentRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Provide a valid json")
		return
	}

	comment, err := h.blogService.CreateComment(userid, slug, req.ParentId, req.Content)

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

func (h *BlogHandler) HandleToggleBlogLike(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")

	if slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid params")
		return
	}

	var req toggleBlogLikeRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	like, err := h.blogService.ToggleBlogLike(userId, slug, req.LikeType)

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

func (h *BlogHandler) HandleRemoveBlogLike(w http.ResponseWriter, r *http.Request) {
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

	err := h.blogService.RemoveBlogLike(userId, slug)

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

func (h *BlogHandler) HandleCreateBookmark(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid params")
		return
	}

	err := h.blogService.CreateBookmark(userId, slug)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Bookmark created successfully"})
}

func (h *BlogHandler) HandleRemoveBookmark(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid params")
		return
	}

	err := h.blogService.RemoveBookmark(userId, slug)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Bookmark removed successfully"})
}
