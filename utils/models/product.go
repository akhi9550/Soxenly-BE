package models

type Image struct {
	Url string `JSON:"url" `
}

type ProductBrief struct {
	ID              uint             `json:"id" gorm:"unique;not null"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	CategoryID      int              `json:"category_id"`
	Price           float64          `json:"price"`
	DiscountedPrice float64          `json:"discounted_price"`
	ProductStatus   string           `json:"product_status"`
	CategoryName    string           `json:"category_name" gorm:"column:category_name"`
	Image           []string         `json:"image" gorm:"-"`
	Variants        []SizeStock      `json:"variants" gorm:"-"`
}

type SizeStock struct {
	Size  string `json:"size"`
	Stock int    `json:"stock"`
}

type ProductReceiver struct {
	Name        string  `json:"name" `
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	Size        string  `json:"size"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}
type Product struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description" validate:"required"`
	CategoryID  uint        `json:"category_id" validate:"required"`
	Price       float64     `json:"price" validate:"required"`
	Variants    []SizeStock `json:"variants" validate:"required"`
}
type Category struct {
	ID       uint   `json:"id"`
	Category string `json:"category"`
	Image    string `json:"image"`
}
type SetNewName struct {
	Current string `json:"current"`
	New     string `json:"new"`
}
type ProductUpdate struct {
	ProductId   int         `json:"product_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CategoryID  int         `json:"category_id"`
	Price       float64     `json:"price"`
	Variants    []SizeStock `json:"variants"`
}
type ProductUpdateReciever struct {
	ProductID int
}
type SearchItems struct {
	ProductName string `json:"product_name"`
}
