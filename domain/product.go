package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `json:"-" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	Price       float64   `json:"price"`
	Variants    []ProductVariant `json:"variants" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type ProductVariant struct {
	gorm.Model
	ProductID uint   `json:"product_id"`
	Size      string `json:"size"`
	Stock     int    `json:"stock"`
}

type ProductImages struct {
	gorm.Model
	ProductImageUrl string `json:"product_image_url"`
}

type Image struct {
	gorm.Model
	ProductId uint   `json:"product_id"`
	Url       string `JSON:"url" `
}

type Category struct {
	gorm.Model
	Category  string `json:"category" gorm:"unique; not null"`
	Image     string `json:"image"`
}
