package service

import (
	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
	"github.com/harry713j/vibe_writer/internal/repo"
)

type UserProfileService struct {
	userRepo    *repo.UserRepository
	profileRepo *repo.UserProfileRepository
	blogRepo    *repo.BlogRepository
	commentRepo *repo.CommentRepository
}

func NewUserProfileService(profile *repo.UserProfileRepository, user *repo.UserRepository,
	blog *repo.BlogRepository, comment *repo.CommentRepository) *UserProfileService {
	return &UserProfileService{
		profileRepo: profile,
		userRepo:    user,
		blogRepo:    blog,
		commentRepo: comment,
	}
}

func (p *UserProfileService) UpdateUserProfile(userId uuid.UUID, fullName, bio string) (*model.UserDetails, error) {

	if _, err := p.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}

	err := p.profileRepo.UpdateProfile(userId, fullName, bio)

	if err != nil {
		return nil, err
	}

	userData, err := p.profileRepo.GetUserDetails(userId)

	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (p *UserProfileService) UpdateAvatar(userId uuid.UUID, avatarUrl string) (*model.UserDetails, error) {
	// get the old avatar url

	if _, err := p.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}

	oldAvatar, err := p.profileRepo.GetAvatarUrl(userId)

	if err != nil {
		return nil, err
	}

	err = p.profileRepo.UpdateAvatar(userId, avatarUrl)

	if err != nil {
		return nil, err
	}

	userData, err := p.profileRepo.GetUserDetails(userId)
	if err != nil {
		return nil, err
	}

	// remove the old avatar from cloud
	go DeleteFromCloud(oldAvatar)

	return userData, nil
}

func (p *UserProfileService) GetProfileDetails(userId uuid.UUID) (*model.UserDetails, error) {
	if _, err := p.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}

	userData, err := p.profileRepo.GetUserDetails(userId)

	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (p *UserProfileService) GetUserDetails(username string) (*model.UserDetails, error) {
	user, err := p.userRepo.GetUserByUsername(username)

	if err != nil {
		return nil, ErrUserNotExists
	}

	userData, err := p.profileRepo.GetUserDetails(user.Id)

	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (p *UserProfileService) RemoveAvatar(userId uuid.UUID) error {

	if _, err := p.userRepo.GetUserById(userId); err != nil {
		return ErrUserNotExists
	}

	err := p.profileRepo.DeleteAvatarUrl(userId)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserProfileService) GetAllCommentsOfBlog(username, slug string) ([]model.CommentWithStat, error) {

	user, err := s.userRepo.GetUserByUsername(username)

	if err != nil {
		return nil, ErrUserNotExists
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
