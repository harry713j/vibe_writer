package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
	"github.com/harry713j/vibe_writer/internal/repo"
)

var (
	ErrInvalidCommentContent = errors.New("comment content is required")
	ErrCommentNotExists      = errors.New("comment not exists")
	ErrInvalidLikeType       = errors.New("invalid like type")
)

type CommentService struct {
	commentRepo *repo.CommentRepository
	userRepo    *repo.UserRepository
	likeRepo    *repo.LikeRepository
}

func NewCommentService(commentRepo *repo.CommentRepository, userRepo *repo.UserRepository,
	likeRepo *repo.LikeRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		likeRepo:    likeRepo,
		userRepo:    userRepo,
	}
}

func (s *CommentService) DeleteComment(userId uuid.UUID, commentId int64) error {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return ErrUserNotExists
	}

	if _, err := s.commentRepo.GetCommentById(userId, commentId); err != nil {
		return ErrCommentNotExists
	}

	err := s.commentRepo.DeleteCommentById(userId, commentId)

	if err != nil {
		return err
	}

	return nil
}

func (s *CommentService) ToggleCommentLike(userId uuid.UUID, commentId int64, liketype model.LikeType) (*model.Like, error) {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}

	if _, err := s.commentRepo.GetCommentById(userId, commentId); err != nil {
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

func (s *CommentService) RemoveCommentLike(userId uuid.UUID, commentId int64) error {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return ErrUserNotExists
	}

	if _, err := s.commentRepo.GetCommentById(userId, commentId); err != nil {
		return ErrCommentNotExists
	}

	return s.likeRepo.DeleteCommentLike(userId, commentId)
}
