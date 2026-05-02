package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
)

func AddBanner(banner domain.Banner) (domain.Banner, error) {
	if err := db.DB.Create(&banner).Error; err != nil {
		return domain.Banner{}, err
	}
	return banner, nil
}

func GetBanners() ([]models.BannerResponse, error) {
	var banners []models.BannerResponse
	err := db.DB.Raw("SELECT id, title1, title2, subtitle1, subtitle2, image, link, is_active FROM banners WHERE deleted_at IS NULL").Scan(&banners).Error
	return banners, err
}

func GetActiveBanners() ([]models.BannerResponse, error) {
	var banners []models.BannerResponse
	err := db.DB.Raw("SELECT id, title1, title2, subtitle1, subtitle2, image, link, is_active FROM banners WHERE is_active = true AND deleted_at IS NULL").Scan(&banners).Error
	return banners, err
}

func DeleteBanner(id int) error {
	if err := db.DB.Exec("UPDATE banners SET deleted_at = NOW() WHERE id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func ToggleBannerStatus(id int) error {
	if err := db.DB.Exec("UPDATE banners SET is_active = NOT is_active WHERE id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
