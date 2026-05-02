package domain

import (
	"Zhooze/utils/models"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Firstname string `json:"firstname" gorm:"validate:required"`
	Lastname  string `json:"lastname" gorm:"validate:required"`
	Email     string `json:"email" gorm:"validate:required"`
	Password  string `json:"password" gorm:"validate:required"`
}
type TokenAdmin struct {
	Admin models.AdminDetailsResponse
	Token string
}
