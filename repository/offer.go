package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"
	"time"
)

func AddProductOffer(productOffer models.ProductOfferReceiver) error {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM product_offers WHERE offer_name = ? AND product_id = ? AND deleted_at IS NULL", productOffer.OfferName, productOffer.ProductID).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the offer already exists")
	}

	// if there is any other offer for this product delete that before adding this one
	count = 0
	err = db.DB.Raw("SELECT COUNT(*) FROM product_offers WHERE product_id = ? AND deleted_at IS NULL", productOffer.ProductID).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		err = db.DB.Where("product_id = ?", productOffer.ProductID).Delete(&domain.ProductOffer{}).Error
		if err != nil {
			return err
		}
	}

	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 5)
	err = db.DB.Exec("INSERT INTO product_offers (product_id, offer_name, discount_percentage, start_date, end_date) VALUES (?, ?, ?, ?, ?)", productOffer.ProductID, productOffer.OfferName, productOffer.DiscountPercentage, startDate, endDate).Error
	if err != nil {
		return err
	}

	return nil

}

func GetOffers() ([]domain.ProductOffer, error) {
	var model []domain.ProductOffer
	err := db.DB.Raw("SELECT * FROM product_offers WHERE deleted_at IS NULL").Scan(&model).Error
	if err != nil {
		return []domain.ProductOffer{}, err
	}

	return model, nil
}

func FindDiscountPercentageForProduct(id int) (int, error) {
	var percentage int
	err := db.DB.Raw("SELECT discount_percentage FROM product_offers WHERE product_id= $1 AND deleted_at IS NULL", id).Scan(&percentage).Error
	if err != nil {
		return 0, err
	}

	return percentage, nil
}

func AddCategoryOffer(categoryOffer models.CategoryOfferReceiver) error {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM category_offers WHERE offer_name = ? AND deleted_at IS NULL", categoryOffer.OfferName).Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the offer already exists")
	}
	// if there is any other offer for this category delete that before adding this one
	count = 0
	err = db.DB.Raw("SELECT COUNT(*) FROM category_offers WHERE category_id = ? AND deleted_at IS NULL", categoryOffer.CategoryID).Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		err = db.DB.Where("category_id = ?", categoryOffer.CategoryID).Delete(&domain.CategoryOffer{}).Error
		if err != nil {
			return err
		}
	}
	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 5)
	err = db.DB.Exec("INSERT INTO category_offers (category_id, offer_name, discount_percentage, start_date, end_date) VALUES (?, ?, ?, ?, ?)", categoryOffer.CategoryID, categoryOffer.OfferName, categoryOffer.DiscountPercentage, startDate, endDate).Error
	if err != nil {
		return err
	}
	return nil
}

func GetCategoryOffer() ([]domain.CategoryOffer, error) {
	var model []domain.CategoryOffer
	err := db.DB.Raw("SELECT * FROM category_offers WHERE deleted_at IS NULL").Scan(&model).Error
	if err != nil {
		return []domain.CategoryOffer{}, err
	}
	return model, nil
}

func FindDiscountPercentageForCategory(id int) (int, error) {
	var percentage int
	err := db.DB.Raw("SELECT discount_percentage FROM category_offers WHERE category_id= $1 AND deleted_at IS NULL", id).Scan(&percentage).Error
	if err != nil {
		return 0, err
	}

	return percentage, nil
}
