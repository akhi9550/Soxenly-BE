package domain

import "gorm.io/gorm"

type WishList struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	Users     User    `json:"-" gorm:"foreignkey:UserID"`
	ProductID uint    `json:"product_id"`
	Products  Product `json:"-" gorm:"foreignkey:ProductID"`
}
