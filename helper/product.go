package helper

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"sync"

	"Zhooze/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func GetImageMimeType(filename string) string {
	extension := strings.ToLower(strings.Split(filename, ".")[len(strings.Split(filename, "."))-1])

	imageMimeTypes := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"bmp":  "image/bmp",
		"webp": "image/webp",
	}

	if mimeType, ok := imageMimeTypes[extension]; ok {
		return mimeType
	}

	return "application/octet-stream"
}

var (
	cld      *cloudinary.Cloudinary
	cldOnce  sync.Once
	cldError error
)

func AddImageToCloudinary(file *multipart.FileHeader) (string, error) {
	cldOnce.Do(func() {
		cfg, err := config.LoadConfig()
		if err != nil {
			cldError = fmt.Errorf("failed to load config: %w", err)
			return
		}

		cld, cldError = cloudinary.NewFromURL(cfg.CloudinaryURL)
		if cldError != nil {
			cldError = fmt.Errorf("failed to initialize Cloudinary: %w", cldError)
		}
	})

	if cldError != nil {
		return "", cldError
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: "soxenly_uploads",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}
