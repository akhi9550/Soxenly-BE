package usecase

import (
	"Zhooze/domain"
	"Zhooze/helper"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"mime/multipart"
)

func GetCategory() ([]domain.Category, error) {
	category, err := repository.GetCategory()
	if err != nil {
		return []domain.Category{}, err
	}
	return category, nil

}

func AddCategory(category models.Category) (domain.Category, error) {
	exists, err := repository.CheckIfCategoryAlreadyExists(category.Category)
	if err != nil {
		return domain.Category{}, err
	}

	if exists {
		return domain.Category{}, errors.New("category already exists")
	}
	categories, err := repository.AddCategory(category)
	if err != nil {
		return domain.Category{}, err
	}
	return categories, nil
}

func UpdateCategory(category models.Category) (domain.Category, error) {
	newCate, err := repository.UpdateCategory(category)
	if err != nil {
		return domain.Category{}, err
	}
	return newCate, nil
}

func DeleteCategory(id int) error {
	err := repository.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCategoryImage(id int, file *multipart.FileHeader) (string, error) {
	url, err := helper.AddImageToCloudinary(file)
	if err != nil {
		return "", err
	}
	err = repository.UpdateCategoryImage(id, url)
	if err != nil {
		return "", err
	}
	return url, nil
}
