package domain

import "gorm.io/gorm"

type Banner struct {
	gorm.Model
	Title1    string `json:"title1"`    // e.g. "SOCKS FOR THE"
	Title2    string `json:"title2"`    // e.g. "CONCRETE."
	Subtitle1 string `json:"subtitle1"` // e.g. "/// COLLECTION 01 — SS26"
	Subtitle2 string `json:"subtitle2"` // e.g. "N° 9550"
	Image     string `json:"image"`
	Link      string `json:"link"`      // Optional link to a category or product
	IsActive  bool   `json:"is_active" gorm:"default:true"`
}
