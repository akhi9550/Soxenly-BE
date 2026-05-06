package usecase

import (
	"Zhooze/domain"
	"Zhooze/helper"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"mime/multipart"
	"strings"
	"sync"
)

func ShowAllProducts(page int, count int, category string, minPrice, maxPrice float64) ([]models.ProductBrief, error) {
	productDetails, err := repository.ShowAllProducts(page, count, category, minPrice, maxPrice)
	if err != nil {
		return []models.ProductBrief{}, err
	}

	for i := range productDetails {
		p := &productDetails[i]
		inStock := false
		for _, v := range p.Variants {
			if v.Stock > 0 {
				inStock = true
				break
			}
		}
		if inStock {
			p.ProductStatus = "in stock"
		} else {
			p.ProductStatus = "out of stock"
		}
		catName, _ := repository.GetCategoryNameByID(p.CategoryID)
		p.CategoryName = catName
	}
	//loop inside products and then calculate discounted price of each then return
	for j := range productDetails {
		discount_percentage, err := repository.FindDiscountPercentageForProduct(int(productDetails[j].ID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64
		if discount_percentage > 0 {
			discount = (productDetails[j].Price * float64(discount_percentage)) / 100
		}
		productDetails[j].DiscountedPrice = productDetails[j].Price - discount

		discount_percentageCategory, err := repository.FindDiscountPercentageForCategory(int(productDetails[j].CategoryID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var categorydiscount float64
		if discount_percentageCategory > 0 {
			categorydiscount = (productDetails[j].Price * float64(discount_percentageCategory)) / 100
		}

		productDetails[j].DiscountedPrice = productDetails[j].DiscountedPrice - categorydiscount
	}
	var updatedproductDetails []models.ProductBrief
	for _, p := range productDetails {
		img, err := repository.GetImage(int(p.ID))
		if err != nil {
			return nil, err
		}
		p.Image = img
		updatedproductDetails = append(updatedproductDetails, p)
	}

	return updatedproductDetails, nil
}

func ShowAllProductsFromAdmin(page int, count int) ([]models.ProductBrief, error) {
	productDetails, err := repository.ShowAllProductsFromAdmin(page, count)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	for i := range productDetails {
		p := &productDetails[i]
		inStock := false
		for _, v := range p.Variants {
			if v.Stock > 0 {
				inStock = true
				break
			}
		}
		if inStock {
			p.ProductStatus = "in stock"
		} else {
			p.ProductStatus = "out of stock"
		}
		catName, _ := repository.GetCategoryNameByID(p.CategoryID)
		p.CategoryName = catName
	}
	for j := range productDetails {
		discount_percentage, err := repository.FindDiscountPercentageForProduct(int(productDetails[j].ID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64
		if discount_percentage > 0 {
			discount = (productDetails[j].Price * float64(discount_percentage)) / 100
		}
		productDetails[j].DiscountedPrice = productDetails[j].Price - discount

		discount_percentageCategory, err := repository.FindDiscountPercentageForCategory(int(productDetails[j].CategoryID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var categorydiscount float64
		if discount_percentageCategory > 0 {
			categorydiscount = (productDetails[j].Price * float64(discount_percentageCategory)) / 100
		}

		productDetails[j].DiscountedPrice = productDetails[j].DiscountedPrice - categorydiscount
	}
	var updatedproductDetails []models.ProductBrief
	for _, p := range productDetails {
		img, err := repository.GetImage(int(p.ID))
		if err != nil {
			return nil, err
		}
		p.Image = img
		updatedproductDetails = append(updatedproductDetails, p)
	}

	return updatedproductDetails, nil

}

func FilterCategory(data map[string]int) ([]models.ProductBrief, error) {
	err := repository.CheckValidateCategory(data)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	var ProductFromCategory []models.ProductBrief
	for _, id := range data {
		product, err := repository.GetProductFromCategory(id)
		if err != nil {
			return []models.ProductBrief{}, err
		}
		for _, products := range product {
			inStock := false
			for _, v := range products.Variants {
				if v.Stock > 0 {
					inStock = true
					break
				}
			}
			if inStock {
				products.ProductStatus = "in stock"
			} else {
				products.ProductStatus = "out of stock"
			}
			if products.ID != 0 {
				ProductFromCategory = append(ProductFromCategory, products)
			}
		}

	}
	for j := range ProductFromCategory {
		discount_percentage, err := repository.FindDiscountPercentageForProduct(int(ProductFromCategory[j].ID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64
		if discount_percentage > 0 {
			discount = (ProductFromCategory[j].Price * float64(discount_percentage)) / 100
		}
		ProductFromCategory[j].DiscountedPrice = ProductFromCategory[j].Price - discount

		discount_percentageCategory, err := repository.FindDiscountPercentageForCategory(int(ProductFromCategory[j].CategoryID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var categorydiscount float64
		if discount_percentageCategory > 0 {
			categorydiscount = (ProductFromCategory[j].Price * float64(discount_percentageCategory)) / 100
		}

		ProductFromCategory[j].DiscountedPrice = ProductFromCategory[j].DiscountedPrice - categorydiscount
	}
	updatedproductDetails := make([]models.ProductBrief, 0)
	for _, p := range ProductFromCategory {
		img, err := repository.GetImage(int(p.ID))
		if err != nil {
			return nil, err
		}
		p.Image = img
		updatedproductDetails = append(updatedproductDetails, p)
	}

	return updatedproductDetails, nil
}

func AddProducts(product models.Product) (domain.Product, error) {
	exist := repository.ProductAlreadyExist(product.Name)
	if exist {
		return domain.Product{}, errors.New("product already exist")
	}
	productResponse, err := repository.AddProducts(product)
	if err != nil {
		return domain.Product{}, err
	}

	return productResponse, nil
}

func DeleteProducts(id string) error {
	err := repository.DeleteProducts(id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProduct(p models.ProductUpdate) (models.ProductUpdateReciever, error) {
	result, err := repository.CheckProductExist(p.ProductId)
	if err != nil {
		return models.ProductUpdateReciever{}, err
	}
	if !result {
		return models.ProductUpdateReciever{}, errors.New("there is no product as you mentioned")
	}
	newcat, err := repository.UpdateProduct(p)
	if err != nil {
		return models.ProductUpdateReciever{}, err
	}
	return newcat, err

}

func DeleteImage(productID int, url string) error {
	err := repository.DeleteImage(productID, url)
	if err != nil {
		return err
	}
	return nil
}
func UpdateProductImage(id int, files []*multipart.FileHeader) error {
	var wg sync.WaitGroup
	urlChan := make(chan string, len(files))
	errChan := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(f *multipart.FileHeader) {
			defer wg.Done()
			url, err := helper.AddImageToCloudinary(f)
			if err != nil {
				errChan <- err
				return
			}
			urlChan <- url
		}(file)
	}

	wg.Wait()
	close(urlChan)
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	var urls []string
	for url := range urlChan {
		urls = append(urls, url)
	}

	if len(urls) > 0 {
		err := repository.UpdateProductImage(id, urls)
		if err != nil {
			return err
		}
	}

	return nil
}

func SearchProductsOnPrefix(prefix string) ([]models.ProductBrief, error) {

	inventoryList, err := repository.GetInventory(prefix)

	if err != nil {
		return nil, err
	}

	var filteredProducts []models.ProductBrief

	for _, product := range inventoryList {
		if strings.HasPrefix(strings.ToLower(product.Name), strings.ToLower(prefix)) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	if len(filteredProducts) == 0 {
		return nil, errors.New("no items matching your keyword")
	}

	return filteredProducts, nil
}

func GetProductDetails(id int) (models.ProductBrief, error) {
	product, err := repository.GetProductDetails(id)
	if err != nil {
		return models.ProductBrief{}, err
	}

	inStock := false
	for _, v := range product.Variants {
		if v.Stock > 0 {
			inStock = true
			break
		}
	}
	if inStock {
		product.ProductStatus = "in stock"
	} else {
		product.ProductStatus = "out of stock"
	}

	return product, nil
}

func GetNewArrivals(count int) ([]models.ProductBrief, error) {
	productDetails, err := repository.GetNewArrivals(count)
	if err != nil {
		return []models.ProductBrief{}, err
	}

	for i := range productDetails {
		p := &productDetails[i]
		inStock := false
		for _, v := range p.Variants {
			if v.Stock > 0 {
				inStock = true
				break
			}
		}
		if inStock {
			p.ProductStatus = "in stock"
		} else {
			p.ProductStatus = "out of stock"
		}
	}

	return productDetails, nil
}
