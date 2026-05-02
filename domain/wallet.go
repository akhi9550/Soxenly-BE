package domain

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID int     `json:"user_id"`
	Users  User    `json:"-" gorm:"foreignkey:UserID"`
	Amount float64 `json:"amount" gorm:"default:0"`
}
type WalletHistory struct {
	gorm.Model
	UserID      int     `json:"user_id"`
	OrderID     int     `json:"order_id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	IsCredited  bool    `json:"is_credited" gorm:"default:true"`
}
