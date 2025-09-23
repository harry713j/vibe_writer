package service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

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

func (p *UserProfileService) UpdateAvatar(userId uuid.UUID, avatarData multipart.File,
	avatarFileName string) (*model.UserDetails, error) {

	ext := filepath.Ext(avatarFileName)

	if ext == "" {
		ext = ".png"
	}

	switch strings.ToLower(ext) {
	case ".jpeg", ".jpg", ".png":
		break
	default:
		return nil, errors.New("file type not supported " + ext)
	}

	baseFileName := uuid.New().String() // for security reason, attacker send ../../these type of file name
	imgDest := filepath.Join("./temp", baseFileName)
	// create on server
	err := p.createImgFile(imgDest, avatarData)

	if err != nil {
		return nil, err
	}
	//upload to cloud
	avatarImgUrl, err := UploadToCloud(imgDest)

	if err != nil {
		return nil, err
	}
	//remove from server
	err = p.removeImgFile(imgDest)

	if err != nil {
		return nil, err
	}

	err = p.profileRepo.UpdateAvatar(userId, avatarImgUrl)

	if err != nil {
		return nil, err
	}

	userData, err := p.profileRepo.GetUserDetails(userId)
	if err != nil {
		return nil, err
	}

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
