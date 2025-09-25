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
		return "", errors.New("This image file not allowed")
	}
	// file created
	newFileName := uuid.New().String() + ext
	file, err := os.Create("./temp/" + newFileName)

	if err != nil {
		return "", err
	}

	if _, err := io.Copy(file, fileData); err != nil {
		return "", err
	}

	imgLocation := filepath.Join("./temp", newFileName)
	// upload to cloud
	imgUrl, err := UploadToCloud(imgLocation)

	if err != nil {
		return "", err
	}

	// remove from the server
	go removeImage(imgLocation)

	return imgUrl, nil
}

func removeImage(fileName string) error {
	return os.Remove(fileName)
}
