package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/harry713j/vibe_writer/internal/middleware"
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
