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

var (
	ErrImageNotAllowed = errors.New("this image not allowed")
)

// for cloud upload
func Upload(fileData multipart.File, fileName string) (string, error) {
	// create the file in the server
	ext := filepath.Ext(fileName)

	if ext == "" {
		ext = ".png" // fallback
	}

	switch strings.ToLower(ext) {
	case ".png", ".jpg", ".jpeg":
		break
	default:
		return "", ErrImageNotAllowed
	}
	// file created
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join("./temp", newFileName)
	safePath := filepath.Clean(filePath)
	file, err := os.Create(safePath)

	if err != nil {
		return "", err
	}

	if _, err := io.Copy(file, fileData); err != nil {
		return "", err
	}

	// upload to cloud
	imgUrl, err := UploadToCloud(safePath)

	if err != nil {
		return "", err
	}

	// remove from the server
	go removeImage(safePath)

	return imgUrl, nil
}

func removeImage(fileName string) error {
	return os.Remove(fileName)
}
