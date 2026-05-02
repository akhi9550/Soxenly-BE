package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"
)

func GetCategory() ([]domain.Category, error) {
	var category []domain.Category
	err := db.DB.Raw("SELECT * FROM categories WHERE deleted_at IS NULL").Scan(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func CheckIfCategoryAlreadyExists(category string) (bool, error) {
	var count int64
	err := db.DB.Raw("SELECT COUNT(*) FROM categories WHERE category = $1 AND deleted_at IS NULL", category).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func AddCategory(category models.Category) (domain.Category, error) {
	newCat := domain.Category{
		Category: category.Category,
		Image:    category.Image,
	}
	if err := db.DB.Create(&newCat).Error; err != nil {
		return domain.Category{}, err
	}
	return newCat, nil
}

func DeleteCategory(id int) error {
	var count int64
	if err := db.DB.Model(&domain.Category{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error; err != nil {
		return err
	}
	if count < 1 {
		return errors.New("category for given id does not exist")
	}

	if err := db.DB.Where("id = ?", id).Delete(&domain.Category{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateCategory(category models.Category) (domain.Category, error) {
	if db.DB == nil {
		return domain.Category{}, errors.New("database connection is nil")
	}
	
	updateData := map[string]interface{}{
		"category": category.Category,
		"image":    category.Image,
	}

	if err := db.DB.Model(&domain.Category{}).Where("id = ?", category.ID).Updates(updateData).Error; err != nil {
		return domain.Category{}, err
	}
	
	var updatedCat domain.Category
	if err := db.DB.First(&updatedCat, category.ID).Error; err != nil {
		return domain.Category{}, err
	}
	return updatedCat, nil
}

func CheckCategory(current string) (bool, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM categories WHERE category=? AND deleted_at IS NULL", current).Scan(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, err
	}
	return true, err
}
func GetCategoryNameByID(id int) (string, error) {
	var categoryName string
	err := db.DB.Raw("SELECT category FROM categories WHERE id = ? AND deleted_at IS NULL", id).Scan(&categoryName).Error
	return categoryName, err
}
func UpdateCategoryImage(id int, url string) error {
	err := db.DB.Model(&domain.Category{}).Where("id = ?", id).Update("image", url).Error
	if err != nil {
		return err
	}
	return nil
}
