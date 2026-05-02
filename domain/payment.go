package domain

import "gorm.io/gorm"

type PaymentMethod struct {
	gorm.Model
	Payment_Name string         `json:"payment_name" gorm:"unique; not null"`
}

type RazerPay struct {
	gorm.Model
	OrderID   string `json:"order_id" `
	Order     Orders  `json:"-" gorm:"foreignkey:OrderID"`
	RazorID   string `json:"razor_id"`
	PaymentID string `json:"payment_id"`
}
