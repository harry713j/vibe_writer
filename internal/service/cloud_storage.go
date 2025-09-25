package service

import (
	"context"
	"log"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/harry713j/vibe_writer/internal/config"
)

// upload to cloud and return cloud url
func UploadToCloud(fileName string) (string, error) {
	cloud, err := config.NewCloud()

	if err != nil {
		return "", err
	}

	ctx := context.Background()
	resp, err := cloud.Upload.Upload(ctx, fileName, uploader.UploadParams{PublicID: fileName, ResourceType: "image"})

	if err != nil {
		return "", err
	}

	log.Println("Cloudinary upload Response: ", resp) // Will Delete later, this line

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
