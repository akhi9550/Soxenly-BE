package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"
)

func GetWishList(userID int) ([]models.WishListResponse, error) {
	var wishList []models.WishListResponse
	err := db.DB.Raw("SELECT products.id, products.name, products.description, products.price FROM products INNER JOIN wish_lists ON products.id = wish_lists.product_id WHERE wish_lists.user_id = ? AND wish_lists.deleted_at IS NULL AND products.deleted_at IS NULL", userID).Scan(&wishList).Error
	if err != nil {
		return []models.WishListResponse{}, err
	}
	
	for i := range wishList {
		var images []string
		db.DB.Raw("SELECT url FROM images WHERE product_id = ?", wishList[i].ID).Scan(&images)
		wishList[i].Image = images
	}
	
	return wishList, nil
}

func ProductExistInWishList(productID, userID int) (bool, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM wish_lists WHERE user_id = ? AND product_id = ? AND deleted_at IS NULL", userID, productID).Scan(&count).Error
	if err != nil {
		return false, errors.New("error checking user product already present")
	}
	return count > 0, nil
}

func AddToWishlist(userID, productID int) error {
	wishlist := domain.WishList{
		UserID:    uint(userID),
		ProductID: uint(productID),
	}
	if err := db.DB.Create(&wishlist).Error; err != nil {
		return err
	}
	return nil
}

func RemoveFromWishList(userID, productID int) error {
	err := db.DB.Exec("DELETE FROM wish_lists WHERE user_id = ? AND product_id = ?", userID, productID).Error
	if err != nil {
		return err
	}
	return nil
}
