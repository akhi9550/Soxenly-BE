package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"

	"gorm.io/gorm"
)

func DisplayCart(userID int) ([]models.Cart, error) {

	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&count).Error; err != nil {
		return []models.Cart{}, err
	}

	if count == 0 {
		return []models.Cart{}, nil
	}

	var cartResponse []models.Cart

	if err := db.DB.Raw("SELECT carts.user_id,users.firstname as user_name,carts.product_id,products.name as product_name,(SELECT url FROM images WHERE images.product_id = carts.product_id LIMIT 1) as image,carts.size,carts.quantity,carts.total_price from carts INNER JOIN users ON carts.user_id = users.id INNER JOIN products ON carts.product_id = products.id WHERE carts.user_id = ? AND carts.deleted_at IS NULL AND users.deleted_at IS NULL AND products.deleted_at IS NULL", userID).Scan(&cartResponse).Error; err != nil {
		return []models.Cart{}, err
	}
	return cartResponse, nil

}

func GetTotalPrice(userID int) (models.CartTotal, error) {

	var cartTotal models.CartTotal
	err := db.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&cartTotal.TotalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = db.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&cartTotal.FinalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = db.DB.Raw("SELECT firstname as user_name FROM users WHERE id = ? AND deleted_at IS NULL", userID).Scan(&cartTotal.UserName).Error
	if err != nil {
		return models.CartTotal{}, err
	}

	return cartTotal, nil

}

func CartExist(userID int) (bool, error) {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}

func EmptyCart(userID int) error {
	if err := db.DB.Exec("DELETE FROM carts WHERE  user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}

func CheckProduct(product_id int) (bool, string, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM products WHERE id = ? AND deleted_at IS NULL", product_id).Scan(&count).Error
	if err != nil {
		return false, "", err
	}
	if count > 0 {
		var category string
		err := db.DB.Raw("SELECT categories.category FROM categories INNER JOIN products ON products.category_id = categories.id WHERE products.id = ? AND products.deleted_at IS NULL AND categories.deleted_at IS NULL", product_id).Scan(&category).Error

		if err != nil {
			return false, "", err
		}
		return true, category, nil
	}
	return false, "", nil
}

func QuantityOfProductInCart(userId int, productId int, size string) (int, error) {
	var productQty int
	err := db.DB.Raw("SELECT quantity FROM carts WHERE user_id = ? AND product_id = ? AND size = ? AND deleted_at IS NULL", userId, productId, size).Scan(&productQty).Error
	if err != nil {
		return 0, err
	}
	return productQty, nil
}

func AddItemIntoCart(userId int, productId int, Quantity int, productprice float64, size string) error {
	cart := domain.Cart{
		UserID:     uint(userId),
		ProductID:  uint(productId),
		Size:       size,
		Quantity:   float64(Quantity),
		TotalPrice: productprice,
	}
	if err := db.DB.Create(&cart).Error; err != nil {
		return err
	}
	return nil
}

func TotalPriceForProductInCart(userID int, productID int, size string) (float64, error) {

	var totalPrice float64
	if err := db.DB.Raw("SELECT SUM(total_price) as total_price FROM carts  WHERE user_id = ? AND product_id = ? AND size = ? AND deleted_at IS NULL", userID, productID, size).Scan(&totalPrice).Error; err != nil {
		return 0.0, err
	}
	return totalPrice, nil
}

func UpdateCart(quantity int, price float64, userID int, product_id int, size string) error {
	if err := db.DB.Model(&domain.Cart{}).Where("user_id = ? and product_id = ? and size = ?", userID, product_id, size).Updates(map[string]interface{}{
		"quantity":    quantity,
		"total_price": price,
	}).Error; err != nil {
		return err
	}
	return nil
}

func ProductExist(userID int, productID int, size string) (bool, error) {
	var count int
	if err := db.DB.Raw("SELECT count(*) FROM carts  WHERE carts.user_id = ? AND product_id = ? AND size = ? AND deleted_at IS NULL", userID, productID, size).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}

func GetQuantityAndProductDetails(userId int, productId int, size string, cartDetails struct {
	Quantity   int
	TotalPrice float64
}) (struct {
	Quantity   int
	TotalPrice float64
}, error) {
	if err := db.DB.Raw("SELECT quantity,total_price FROM carts WHERE user_id = ? AND product_id = ? AND size = ? AND deleted_at IS NULL", userId, productId, size).Scan(&cartDetails).Error; err != nil {
		return struct {
			Quantity   int
			TotalPrice float64
		}{}, err
	}
	return cartDetails, nil
}

func RemoveProductFromCart(userID int, product_id int, size string) error {

	if err := db.DB.Exec("DELETE FROM carts WHERE user_id = ? AND product_id = ? AND size = ?", uint(userID), uint(product_id), size).Error; err != nil {
		return err
	}

	return nil
}

func UpdateCartDetails(cartDetails struct {
	Quantity   int
	TotalPrice float64
}, userId int, productId int, size string) error {
	if err := db.DB.Model(&domain.Cart{}).Where("user_id = ? AND product_id = ? AND size = ?", userId, productId, size).Updates(map[string]interface{}{
		"quantity":    cartDetails.Quantity,
		"total_price": cartDetails.TotalPrice,
	}).Error; err != nil {
		return err
	}
	return nil
}

func CartAfterRemovalOfProduct(user_id int) ([]models.Cart, error) {
	var cart []models.Cart
	if err := db.DB.Raw("SELECT carts.product_id,products.name as product_name,(SELECT url FROM images WHERE images.product_id = carts.product_id LIMIT 1) as image,carts.size,carts.quantity,carts.total_price FROM carts INNER JOIN products on carts.product_id = products.id WHERE carts.user_id = ? AND carts.deleted_at IS NULL AND products.deleted_at IS NULL", user_id).Scan(&cart).Error; err != nil {
		return []models.Cart{}, err
	}
	return cart, nil
}

func GetAllItemsFromCart(userID int) ([]models.Cart, error) {
	var count int
	var cartResponse []models.Cart
	err := db.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&count).Error
	if err != nil {
		return []models.Cart{}, err
	}
	if count == 0 {
		return []models.Cart{}, nil
	}
	err = db.DB.Raw("SELECT carts.user_id,users.firstname as user_name,carts.product_id,products.name as product_name,(SELECT url FROM images WHERE images.product_id = carts.product_id LIMIT 1) as image,carts.size,carts.quantity,carts.total_price from carts INNER JOIN users on carts.user_id = users.id INNER JOIN products ON carts.product_id = products.id where carts.user_id = ? AND carts.deleted_at IS NULL AND users.deleted_at IS NULL AND products.deleted_at IS NULL", userID).Scan(&cartResponse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if len(cartResponse) == 0 {
				return []models.Cart{}, nil
			}
			return []models.Cart{}, err
		}
		return []models.Cart{}, err
	}
	return cartResponse, nil
}

func GetTotalPriceFromCart(userID int) (float64, error) {
	var totalPrice float64
	err := db.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&totalPrice).Error
	if err != nil {
		return 0.0, err
	}
	return totalPrice, nil

}
