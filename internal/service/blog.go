package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
	"github.com/harry713j/vibe_writer/internal/repo"
	"github.com/harry713j/vibe_writer/internal/utils"
)

type BlogService struct {
	blogRepo *repo.BlogRepository
	userRepo *repo.UserRepository
}

var (
	ErrBlogNotExists = errors.New("no blog exists with this slug")
	ErrTitleExists   = errors.New("blog with this title already exists")
	ErrUserNotExists = errors.New("user not found")
)

func NewBlogService(blogRepo *repo.BlogRepository, userRepo *repo.UserRepository) *BlogService {
	return &BlogService{
		blogRepo: blogRepo,
		userRepo: userRepo,
	}
}

// create blog
func (r *BlogService) CreateBlog(userId uuid.UUID, title, slug, content string, photoUrls []string) (*model.BlogResponse, error) {

	err := r.blogRepo.GetBlogByTitle(userId, title)

	if err == nil {
		return nil, ErrTitleExists
	}

	randomHex, err := utils.RandomHex(16) // give a random hex of length 16
	if err != nil {
		return nil, err
	}

	transformedSlug := slug + "-" + randomHex
	// create the blog
	blogId, err := r.blogRepo.CreateBlog(userId, title, transformedSlug, content)

	if err != nil {
		return nil, err
	}
	// store the blog images
	for _, url := range photoUrls {
		if _, err := r.blogRepo.CreateBlogImage(blogId, url); err != nil {
			return nil, err
		}
	}

	// get that blog
	blog, err := r.blogRepo.GetBlogById(userId, blogId)

	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (r *BlogService) UpdateBlog(userId uuid.UUID, slug, title, content string, photoUrls []string) (*model.BlogResponse, error) {

	// check blog exists or not
	if _, err := r.blogRepo.GetBlogBySlug(userId, slug); err != nil {
		return nil, ErrBlogNotExists
	}

	blogId, err := r.blogRepo.UpdateBlog(userId, slug, title, content)

	if err != nil {
		return nil, err
	}

	existingPhotoUrls, err := r.blogRepo.GetPhotoUrls(blogId)

	if err != nil {
		return nil, err
	}

	removedUrls := r.differentUrls(existingPhotoUrls, photoUrls) // old photo urls to be remove
	addUrls := r.differentUrls(photoUrls, existingPhotoUrls)     // new photo urls to be add

	// remove from table
	err = r.blogRepo.DeleteBlogPhotosByURLs(blogId, removedUrls)

	if err != nil {
		return nil, err
	}
	// create new photoUrls
	for _, url := range addUrls {
		if _, err := r.blogRepo.CreateBlogImage(blogId, url); err != nil {
			return nil, err
		}
	}

	// get that blog
	blog, err := r.blogRepo.GetBlogById(userId, blogId)

	if err != nil {
		return nil, err
	}

	// remove the old photo url from cloud
	for _, url := range removedUrls {
		go DeleteFromCloud(url)
	}

	return blog, err
}

func (r *BlogService) GetAllUserBlog(username string, page, limit int) (*model.PaginatedResponse[model.BlogSummary], error) {
	user, err := r.userRepo.GetUserByUsername(username)

	if err != nil {
		return nil, ErrUserNotExists
	}

	blogs, err := r.blogRepo.GetAllBlog(user.Id, page, limit)

	if err != nil {
		return nil, err
	}

	return blogs, nil
}

// return `BlogDetails` with error
func (r *BlogService) GetBlog(username string, slug string) (*model.BlogResponse, error) {
	// get the user details by username
	user, err := r.userRepo.GetUserByUsername(username)

	if err != nil {
		return nil, ErrUserNotExists
	}

	blog, err := r.blogRepo.GetBlogBySlug(user.Id, slug)

	if err != nil {
		return nil, ErrBlogNotExists
	}

	return blog, nil
}

func (r *BlogService) DeleteBlog(userId uuid.UUID, slug string) error {
	// check the blog exists or not
	if _, err := r.blogRepo.GetBlogBySlug(userId, slug); err != nil {
		return ErrBlogNotExists
	}

	return r.blogRepo.DeleteBlog(userId, slug)
}

// returns string in a that are not present in b
func (r *BlogService) differentUrls(a, b []string) []string {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	var diff []string
	for _, item := range a {
		if !m[item] {
			diff = append(diff, item)
		}
	}

	return diff
}
