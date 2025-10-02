package service

import (
	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
	"github.com/harry713j/vibe_writer/internal/repo"
)

type UserProfileService struct {
	userRepo    *repo.UserRepository
	profileRepo *repo.UserProfileRepository
}

func NewUserProfileService(profile *repo.UserProfileRepository, user *repo.UserRepository) *UserProfileService {
	return &UserProfileService{
		profileRepo: profile,
		userRepo:    user,
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
