package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type Cloudinary struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinary(cloudName, apiKey, apiSecret string) *Cloudinary {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		// Handle initialization error in production
		return &Cloudinary{cld: nil}
	}
	return &Cloudinary{cld: cld}
}

func (c *Cloudinary) Upload(ctx context.Context, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	uploadParams := uploader.UploadParams{
		PublicID: uuid.New().String(),
	}
	resp, err := c.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	return resp.SecureURL, nil
}