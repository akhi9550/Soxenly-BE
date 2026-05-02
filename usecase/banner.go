package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

func AddBanner(banner models.BannerRequest, file *multipart.FileHeader) (domain.Banner, error) {
	// Handle file upload
	fileName := fmt.Sprintf("banner_%d%s", os.Getpid(), filepath.Ext(file.Filename))
	filePath := filepath.Join("uploads", fileName)

	// Create uploads directory if not exists
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	// This is a simplified mock of file saving logic usually found in handlers,
	// but we'll define the domain object here.
	newBanner := domain.Banner{
		Title1:    banner.Title1,
		Title2:    banner.Title2,
		Subtitle1: banner.Subtitle1,
		Subtitle2: banner.Subtitle2,
		Link:      banner.Link,
		Image:     filePath,
	}

	return repository.AddBanner(newBanner)
}

func GetBannersForAdmin() ([]models.BannerResponse, error) {
	return repository.GetBanners()
}

func GetBannersForUser() ([]models.BannerResponse, error) {
	banners, err := repository.GetActiveBanners()
	if err != nil {
		return nil, err
	}

	// Fallback to default if no banners are active
	if len(banners) == 0 {
		return []models.BannerResponse{
			{
				Title1:    "SOCKS FOR THE",
				Title2:    "CONCRETE.",
				Subtitle1: "/// COLLECTION 01 — SS26",
				Subtitle2: "N° 9550",
				Image:     "/24265_7.jpeg",
				Link:      "/shop",
				IsActive:  true,
			},
		}, nil
	}

	return banners, nil
}

func DeleteBanner(id int) error {
	return repository.DeleteBanner(id)
}

func ToggleBannerStatus(id int) error {
	return repository.ToggleBannerStatus(id)
}
