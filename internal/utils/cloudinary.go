package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(file *multipart.FileHeader, folder string) (string, string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to init cloudinary: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	ctx := context.Background()
	resp, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload: %w", err)
	}

	return resp.SecureURL, resp.PublicID, nil
}
