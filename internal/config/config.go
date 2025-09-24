package config

import (
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

type DBConfig struct {
	URL string
}

func LoadDBConfig() *DBConfig {
	dbUrl := os.Getenv("DATABASE_URL")
	return &DBConfig{URL: dbUrl}
}

func NewCloud() (*cloudinary.Cloudinary, error) {
	cloudName := os.Getenv("CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	return cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
}
