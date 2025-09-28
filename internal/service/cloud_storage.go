package service

import (
	"context"
	"errors"
	"log"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/harry713j/vibe_writer/internal/config"
)

// upload to cloud and return cloud url
func UploadToCloud(filePath, fileName string) (string, error) {
	cloud, err := config.NewCloud()

	if err != nil {
		return "", err
	}

	ctx := context.Background()
	resp, err := cloud.Upload.Upload(ctx, filePath, uploader.UploadParams{PublicID: fileName, ResourceType: "image"})

	if err != nil {
		log.Println("Err: ", err)
		return "", err
	}

	log.Println("Cloudinary upload Response: ", resp) // Will Delete later, this line

	if resp.SecureURL == "" {
		return "", errors.New("failed to upload to cloud")
	}

	return resp.SecureURL, nil
}

func DeleteFromCloud(imgUrl string) error {

	cloud, err := config.NewCloud()

	if err != nil {
		return err
	}

	fileName := filepath.Base(imgUrl)
	ctx := context.Background()

	resp, err := cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: fileName})

	log.Println("Cloudinary upload Response: ", resp) // Will Delete later, this line

	return err
}
