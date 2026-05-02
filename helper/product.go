package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
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

func AddImageToS3(file *multipart.FileHeader) (string, error) {
	// Create uploads directory if not exists
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Create unique filename using timestamp to avoid collisions
	filename := fmt.Sprintf("%d_%s", os.Getpid(), strings.ReplaceAll(file.Filename, " ", "_"))
	dst := fmt.Sprintf("%s/%s", uploadDir, filename)
	
	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy content
	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}
