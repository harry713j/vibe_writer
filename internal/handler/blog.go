package handler

import (
	"encoding/json"
	"errors"
	"net/http"

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

func (h *BlogHandler) HandleGetAllBlog(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogDatas, err := h.blogService.GetAllUserBlog(userId)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string][]*model.BlogDetails{"blogs": blogDatas})
}

func (h *BlogHandler) HandleGetBlog(w http.ResponseWriter, r *http.Request) {

	// extract it from parameter
	slug := chi.URLParam(r, "slug")
	username := chi.URLParam(r, "username")
	if slug == "" || username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Blog parameter required")
		return
	}

	BlogData, err := h.blogService.GetBlog(username, slug)

	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) || errors.Is(err, service.ErrBlogNotExists) {
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, BlogData)
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
