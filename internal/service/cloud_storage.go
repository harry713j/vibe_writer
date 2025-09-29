package service

import (
	"context"
	"errors"
	"log"
	"path/filepath"
	"strings"

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
	publicId := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	log.Printf("imgUrl: %v , fileName: %v\n", imgUrl, publicId)
	ctx := context.Background()

	resp, err := cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicId, ResourceType: "image"})

	if err != nil {
		return err
	}

	log.Println("Cloudinary delete Response: ", resp) // Will Delete later, this line

	return nil
}
