package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
	"github.com/harry713j/vibe_writer/internal/repo"
)

var (
	ErrInvalidLikeType = errors.New("invalid like type")
)

type LikeService struct {
	likeRepo    *repo.LikeRepository
	userRepo    *repo.UserRepository
	blogRepo    *repo.BlogRepository
	commentRepo *repo.CommentRepository
}

func NewLikeService(likeRepo *repo.LikeRepository, userRepo *repo.UserRepository, blogRepo *repo.BlogRepository,
	commentRepo *repo.CommentRepository) *LikeService {
	return &LikeService{likeRepo: likeRepo, userRepo: userRepo, blogRepo: blogRepo, commentRepo: commentRepo}
}

func (s *LikeService) ToggleBlogLike(userId uuid.UUID, slug string, liketype model.LikeType) (*model.Like, error) {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}

	blog, err := s.blogRepo.GetBlogBySlug(userId, slug)

	if err != nil {
		return nil, ErrBlogNotExists
	}

	if liketype != "like" && liketype != "dislike" {
		return nil, ErrInvalidLikeType
	}

	like, err := s.likeRepo.UpsertBlogLike(userId, blog.Id, liketype)

	if err != nil {
		return nil, err
	}

	return like, nil
}

func (s *LikeService) ToggleCommentLike(userId uuid.UUID, commentId int64, liketype model.LikeType) (*model.Like, error) {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}

	if _, err := s.commentRepo.GetComment(commentId); err != nil {
		return nil, ErrCommentNotExists
	}

	if liketype != "like" && liketype != "dislike" {
		return nil, ErrInvalidLikeType
	}

	like, err := s.likeRepo.UpsertCommentLike(userId, commentId, liketype)

	if err != nil {
		return nil, err
	}

	return like, nil
}

func (s *LikeService) RemoveBlogLike(userId uuid.UUID, slug string) error {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return ErrUserNotExists
	}

	blog, err := s.blogRepo.GetBlogBySlug(userId, slug)

	if err != nil {
		return ErrBlogNotExists
	}

	return s.likeRepo.DeleteBlogLike(userId, blog.Id)
}

func (s *LikeService) RemoveCommentLike(userId uuid.UUID, commentId int64) error {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return ErrUserNotExists
	}

	if _, err := s.commentRepo.GetComment(commentId); err != nil {
		return ErrCommentNotExists
	}

	return s.likeRepo.DeleteCommentLike(userId, commentId)
}
