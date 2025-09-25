package service

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
	"github.com/harry713j/vibe_writer/internal/repo"
)

type UserProfileService struct {
	profileRepo *repo.UserProfileRepository
}

func NewUserProfileService(profile *repo.UserProfileRepository) *UserProfileService {
	return &UserProfileService{
		profileRepo: profile,
	}
}

func (p *UserProfileService) UpdateUserProfile(userId uuid.UUID, fullName, bio string) (*model.UserDetails, error) {
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

	// reove the old avatar from cloud
	go DeleteFromCloud(oldAvatar)

	return userData, nil
}

func (p *UserProfileService) GetUserDetails(userId uuid.UUID) (*model.UserDetails, error) {
	userData, err := p.profileRepo.GetUserDetails(userId)

	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (p *UserProfileService) createImgFile(imgLocation string, fileData multipart.File) error {
	// #nosec G304 -- safe: filename sanitized and stored only in ./temp
	file, err := os.Create(imgLocation)

	if err != nil {
		return err
	}

	if _, err := io.Copy(file, fileData); err != nil {
		return err
	}

	return nil
}

func (p *UserProfileService) removeImgFile(imgLocation string) error {
	return os.Remove(imgLocation)
}
