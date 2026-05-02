package models

type WishListResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Image       []string `json:"image" gorm:"-"`
}
