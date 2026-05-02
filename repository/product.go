package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"
	"log"
	"strconv"
)

func ShowAllProducts(page int, count int, category string, minPrice, maxPrice float64) ([]models.ProductBrief, error) {
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * count
	var productBrief []models.ProductBrief
	
	query := `SELECT p.*, c.category as category_name FROM products p LEFT JOIN categories c ON p.category_id = c.id WHERE p.deleted_at IS NULL`
	var args []interface{}

	if category != "" && category != "all" {
		query += ` AND c.category ILIKE ?`
		args = append(args, category)
	}
	if minPrice > 0 {
		query += ` AND p.price >= ?`
		args = append(args, minPrice)
	}
	if maxPrice > 0 {
		query += ` AND p.price <= ?`
		args = append(args, maxPrice)
	}

	query += ` LIMIT ? OFFSET ?`
	args = append(args, count, offset)

	err := db.DB.Raw(query, args...).Scan(&productBrief).Error
	if err != nil {
		return nil, err
	}
	
	for i := range productBrief {
		variants, _ := GetVariants(int(productBrief[i].ID))
		productBrief[i].Variants = variants
	}
	
	return productBrief, nil
}

func GetVariants(productID int) ([]models.SizeStock, error) {
	var variants []models.SizeStock
	err := db.DB.Raw("SELECT size, stock FROM product_variants WHERE product_id = ? AND deleted_at IS NULL", productID).Scan(&variants).Error
	return variants, err
}

func ShowAllProductsFromAdmin(page int, count int) ([]models.ProductBrief, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var productBrief []models.ProductBrief
	err := db.DB.Raw(`SELECT p.*, c.category as category_name FROM products p LEFT JOIN categories c ON p.category_id = c.id WHERE p.deleted_at IS NULL limit ? offset ?`, count, offset).Scan(&productBrief).Error
	if err != nil {
		return nil, err
	}
	
	for i := range productBrief {
		variants, _ := GetVariants(int(productBrief[i].ID))
		productBrief[i].Variants = variants
	}
	
	return productBrief, nil
}

func CheckValidateCategory(data map[string]int) error {
	for _, id := range data {
		var count int
		err := db.DB.Raw("SELECT COUNT(*) FROM categories WHERE id=? AND deleted_at IS NULL", id).Scan(&count).Error
		if err != nil {
			return err
		}
		if count < 1 {
			return errors.New("doesn't exist")
		}
	}
	return nil
}

func GetProductFromCategory(id int) ([]models.ProductBrief, error) {
	var product []models.ProductBrief
	err := db.DB.Raw(`SELECT * FROM products JOIN categories ON products.category_id=categories.id WHERE categories.id=? AND products.deleted_at IS NULL`, id).Scan(&product).Error
	if err != nil {
		return []models.ProductBrief{}, err
	}
	return product, nil
}

func GetQuantityFromProductIDAndSize(id int, size string) (int, error) {
	var quantity int
	err := db.DB.Raw("SELECT stock FROM product_variants WHERE product_id= ? AND size = ? AND deleted_at IS NULL", id, size).Scan(&quantity).Error
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

func GetPriceOfProductFromID(prodcut_id int) (float64, error) {
	var productPrice float64

	if err := db.DB.Raw("SELECT price FROM products WHERE id = ? AND deleted_at IS NULL", prodcut_id).Scan(&productPrice).Error; err != nil {
		return 0.0, err
	}
	return productPrice, nil
}

func ProductAlreadyExist(Name string) bool {
	var count int
	if err := db.DB.Raw("SELECT count(*) FROM products WHERE name = ? AND deleted_at IS NULL", Name).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func FindCategoryID(id int) (int, error) {
	var a int
	if err := db.DB.Raw("SELECT category_id FROM products WHERE id = ? AND deleted_at IS NULL", id).Scan(&a).Error; err != nil {
		return 0.0, err
	}
	return a, nil
}

func StockInvalid(Name string) bool {
	var count int
	query := `SELECT COALESCE(SUM(pv.stock), 0) 
              FROM product_variants pv 
              JOIN products p ON pv.product_id = p.id 
              WHERE p.name = ? AND p.deleted_at IS NULL AND pv.deleted_at IS NULL`
	if err := db.DB.Raw(query, Name).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func AddProducts(product models.Product) (domain.Product, error) {
	p := domain.Product{
		Name:        product.Name,
		Description: product.Description,
		CategoryID:  uint(product.CategoryID),
		Price:       product.Price,
	}
	
	if err := db.DB.Create(&p).Error; err != nil {
		log.Println(err.Error())
		return domain.Product{}, err
	}

	for _, v := range product.Variants {
		variant := domain.ProductVariant{
			ProductID: p.ID,
			Size:      v.Size,
			Stock:     v.Stock,
		}
		db.DB.Create(&variant)
	}

	return p, nil
}

func DeleteProducts(id string) error {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM products WHERE id=? AND deleted_at IS NULL", product_id).Scan(&count).Error; err != nil {
		return err
	}
	if count < 1 {
		return errors.New("product for given id does not exist")
	}
	if err := db.DB.Where("id = ?", product_id).Delete(&domain.Product{}).Error; err != nil {
		return err
	}
	// Cascading soft delete for variants and images
	db.DB.Where("product_id = ?", product_id).Delete(&domain.ProductVariant{})
	db.DB.Where("product_id = ?", product_id).Delete(&domain.Image{})
	return nil
}

func CheckProductExist(pid int) (bool, error) {
	var a int
	err := db.DB.Raw("SELECT COUNT(*) FROM products WHERE id=? AND deleted_at IS NULL", pid).Scan(&a).Error
	if err != nil {
		return false, err
	}
	if a == 0 {
		return false, err
	}
	return true, err
}

func UpdateProduct(p models.ProductUpdate) (models.ProductUpdateReciever, error) {
	for _, v := range p.Variants {
		if v.Stock < 0 {
			return models.ProductUpdateReciever{}, errors.New("stock cannot be negative")
		}
	}
	if db.DB == nil {
		return models.ProductUpdateReciever{}, errors.New("database connection is nil")
	}
	
	updateData := map[string]interface{}{
		"name":        p.Name,
		"description": p.Description,
		"category_id": p.CategoryID,
		"price":       p.Price,
	}

	if err := db.DB.Model(&domain.Product{}).Where("id = ?", p.ProductId).Updates(updateData).Error; err != nil {
		return models.ProductUpdateReciever{}, err
	}
	
	// Soft delete old variants and add new ones
	db.DB.Where("product_id = ?", p.ProductId).Delete(&domain.ProductVariant{})
	for _, v := range p.Variants {
		variant := domain.ProductVariant{
			ProductID: uint(p.ProductId),
			Size:      v.Size,
			Stock:     v.Stock,
		}
		db.DB.Create(&variant)
	}

	var newdetails models.ProductUpdateReciever
	newdetails.ProductID = p.ProductId
	return newdetails, nil
}

func DoesProductExist(productID int) (bool, error) {
	var count int
	err := db.DB.Raw("select count(*) from products where id = ? AND deleted_at IS NULL", productID).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func UpdateProductImage(productID int, url string) error {
	err := db.DB.Exec("INSERT INTO images (product_id,url) VALUES ($1,$2) RETURNING * ", productID, url).Error
	if err != nil {
		return errors.New("error while insert image to database")
	}
	return nil
}

func DeleteImage(productID int, url string) error {
	err := db.DB.Where("product_id = ? AND url = ?", productID, url).Delete(&domain.Image{}).Error
	if err != nil {
		return err
	}
	return nil
}

func DisplayImages(productID int) (domain.Product, []domain.Image, error) {
	var product domain.Product
	var image []domain.Image
	err := db.DB.Raw(`SELECT * FROM products WHERE product_id = $1 AND deleted_at IS NULL`, productID).Scan(&product).Error
	if err != nil {
		return domain.Product{}, []domain.Image{}, err
	}
	err = db.DB.Raw(`SELECT * FROM images WHERE product_id = $1 AND deleted_at IS NULL`, productID).Scan(&image).Error
	if err != nil {
		return domain.Product{}, []domain.Image{}, err
	}
	return product, image, nil
}

func GetImage(productID int) ([]string, error) {
	var url []string
	if err := db.DB.Raw(`SELECT url FROM Images WHERE product_id=? AND deleted_at IS NULL`, productID).Scan(&url).Error; err != nil {
		return url, err
	}
	return url, nil

}

func GetInventory(prefix string) ([]models.ProductBrief, error) {
	var productDetails []models.ProductBrief
	query := `
	SELECT i.*
	FROM products i
	LEFT JOIN categories c ON i.category_id = c.id
	WHERE (i.name ILIKE '%' || $1 || '%'
    OR c.category ILIKE '%' || $1 || '%')
	AND i.deleted_at IS NULL;`
	if err := db.DB.Raw(query, prefix).Scan(&productDetails).Error; err != nil {
		return []models.ProductBrief{}, err
	}

	return productDetails, nil
}

func GetProductDetails(id int) (models.ProductBrief, error) {
	var product models.ProductBrief
	query := `SELECT p.*, c.category as category_name FROM products p LEFT JOIN categories c ON p.category_id = c.id WHERE p.id = ? AND p.deleted_at IS NULL`
	err := db.DB.Raw(query, id).Scan(&product).Error
	if err != nil {
		return models.ProductBrief{}, err
	}
	if product.ID == 0 {
		return models.ProductBrief{}, errors.New("product not found")
	}
	
	variants, _ := GetVariants(int(product.ID))
	product.Variants = variants
	
	img, _ := GetImage(int(product.ID))
	product.Image = img
	
	return product, nil
}

func GetNewArrivals(count int) ([]models.ProductBrief, error) {
	var productBrief []models.ProductBrief
	query := `SELECT p.*, c.category as category_name FROM products p LEFT JOIN categories c ON p.category_id = c.id WHERE p.deleted_at IS NULL ORDER BY p.created_at DESC LIMIT ?`
	err := db.DB.Raw(query, count).Scan(&productBrief).Error
	if err != nil {
		return nil, err
	}
	
	for i := range productBrief {
		variants, _ := GetVariants(int(productBrief[i].ID))
		productBrief[i].Variants = variants
		img, _ := GetImage(int(productBrief[i].ID))
		productBrief[i].Image = img
	}
	
	return productBrief, nil
}
