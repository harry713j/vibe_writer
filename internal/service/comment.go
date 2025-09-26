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
)

type CommentService struct {
	commentRepo *repo.CommentRepository
	blogRepo    *repo.BlogRepository
	userRepo    *repo.UserRepository
}

func NewCommentService(commentRepo *repo.CommentRepository, blogRepo *repo.BlogRepository,
	userRepo *repo.UserRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		blogRepo:    blogRepo,
		userRepo:    userRepo,
	}
}

func (s *CommentService) CreateComment(userId uuid.UUID, slug string, parentId int64, content string) (*model.Comment, error) {
	// check user exists or not
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}
	// check blog with blog id exists or not
	blog, err := s.blogRepo.GetBlogBySlug(userId, slug)

	if err != nil {
		return nil, ErrBlogNotExists
	}

	if content == "" {
		return nil, ErrInvalidCommentContent
	}

	comment, err := s.commentRepo.CreateComment(userId, blog.Id, parentId, content)

	if err != nil {
		return nil, err
	}

	return comment, nil
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

func (s *CommentService) GetAllCommentsOfBlog(username, slug string) ([]*model.Comment, error) {
	user, err := s.userRepo.GetUserByUsername(username)

	if err != nil {
		return nil, ErrBlogNotExists
	}

	blog, err := s.blogRepo.GetBlogBySlug(user.Id, slug)

	if err != nil {
		return nil, ErrBlogNotExists
	}

	comments, err := s.commentRepo.GetCommentsByBlogId(blog.Id)

	if err != nil {
		return nil, err
	}

	return comments, nil
}
