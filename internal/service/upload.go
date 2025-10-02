package service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type UploadService struct{}

var (
	ErrImageNotAllowed = errors.New("image type not allowed")
)

func NewUploadService() *UploadService {
	return &UploadService{}
}

// for cloud upload
func (s *UploadService) Upload(fileData multipart.File, fileName string) (string, error) {
	// create the file in the server
	ext := filepath.Ext(fileName)

	if ext == "" {
		ext = ".png"
	}

	switch strings.ToLower(ext) {
	case ".png", ".jpg", ".jpeg":
		break
	default:
		return "", ErrImageNotAllowed
	}
	// file created
	newFileName := uuid.New().String()
	filePath := filepath.Join("./temp", newFileName)
	safePath := filepath.Clean(filePath)

	if err := os.MkdirAll("./temp", 0755); err != nil {
		return "", err
	}

	file, err := os.Create(safePath)

	if err != nil {
		return "", err
	}

	if _, err := io.Copy(file, fileData); err != nil {
		return "", err
	}

	file.Close()

	// upload to cloud
	imgUrl, err := UploadToCloud(safePath, newFileName)

	if err != nil {
		return "", err
	}

	// remove from the server
	go s.removeImage(safePath)

	return imgUrl, nil
}

func (s *UploadService) removeImage(fileName string) error {
	return os.Remove(fileName)
}
