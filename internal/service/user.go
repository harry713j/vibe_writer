package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
	"github.com/harry713j/vibe_writer/internal/repo"
)

var (
	ErrInvalidFollowingUser = errors.New("user not exist to follow or unfollow")
	ErrInvalidAuthor        = errors.New("author not exists")
)

type UserProfileService struct {
	userRepo    *repo.UserRepository
	profileRepo *repo.UserProfileRepository
	blogRepo    *repo.BlogRepository
	commentRepo *repo.CommentRepository
	followRepo  *repo.FollowRepository
}

func NewUserProfileService(profile *repo.UserProfileRepository, user *repo.UserRepository,
	blog *repo.BlogRepository, comment *repo.CommentRepository, followRepo *repo.FollowRepository) *UserProfileService {
	return &UserProfileService{
		profileRepo: profile,
		userRepo:    user,
		blogRepo:    blog,
		commentRepo: comment,
		followRepo:  followRepo,
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

func (s *UserProfileService) FetchBookmarks(userId uuid.UUID) ([]model.BlogSummary, error) {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}

	blogs, err := s.profileRepo.GetAllBookmarks(userId)

	if err != nil {
		return nil, err
	}

	return blogs, nil
}

/* Follow */
func (s *UserProfileService) CreateFollow(userId uuid.UUID, followingUsername string) error {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return ErrUserNotExists
	}

	followingUser, err := s.userRepo.GetUserByUsername(followingUsername)
	if err != nil {
		return ErrInvalidFollowingUser
	}

	err = s.followRepo.Create(userId, followingUser.Id)
	return err
}

func (s *UserProfileService) RemoveFollow(userId uuid.UUID, followingUsername string) error {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return ErrUserNotExists
	}

	followingUser, err := s.userRepo.GetUserByUsername(followingUsername)
	if err != nil {
		return ErrInvalidFollowingUser
	}

	err = s.followRepo.Delete(userId, followingUser.Id)
	return err
}

func (s *UserProfileService) FetchAllFollower(userId uuid.UUID, followingUsername string, page, limit int) (*model.PaginatedResponse[model.FollowResponse], error) {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}

	followingUser, err := s.userRepo.GetUserByUsername(followingUsername)
	if err != nil {
		return nil, ErrInvalidAuthor
	}

	followers, err := s.followRepo.GetAllFollower(followingUser.Id, page, limit)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (s *UserProfileService) FetchAllFollowing(userId uuid.UUID, followerUsername string, page, limit int) (*model.PaginatedResponse[model.FollowResponse], error) {
	if _, err := s.userRepo.GetUserById(userId); err != nil {
		return nil, ErrUserNotExists
	}

	followerUser, err := s.userRepo.GetUserByUsername(followerUsername)
	if err != nil {
		return nil, ErrInvalidAuthor
	}

	followings, err := s.followRepo.GetAllFollowing(followerUser.Id, page, limit)
	if err != nil {
		return nil, err
	}

	return followings, nil
}
