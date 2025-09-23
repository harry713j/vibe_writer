package service

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloud(fileName string) (string, error) {
	cloudName := os.Getenv("CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cloud, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)

	if err != nil {
		return "", err
	}

	ctx := context.Background()
	resp, err := cloud.Upload.Upload(ctx, fileName, uploader.UploadParams{PublicID: fileName})

	if err != nil {
		return "", err
	}

	log.Println("Cloudinary Response: ", resp) // Will Delete later, this line

	return resp.SecureURL, nil
}
