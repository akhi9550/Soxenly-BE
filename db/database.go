package db

import (
	"Zhooze/config"
	"Zhooze/domain"
	"Zhooze/helper"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(confg config.Config) (*gorm.DB, error) {
	connectTo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", confg.DBHost, confg.DBUser, confg.DBName, confg.DBPort, confg.DBPassword)
	db, err := gorm.Open(postgres.Open(connectTo), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database:%w", err)
	}
	DB = db
	db.AutoMigrate(&domain.Admin{})
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Product{})
	db.AutoMigrate(&domain.ProductVariant{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Address{})
	db.AutoMigrate(&domain.Cart{})
	db.AutoMigrate(&domain.Orders{})
	db.AutoMigrate(&domain.OrderItem{})
	db.AutoMigrate(&domain.RazerPay{})
	db.AutoMigrate(&domain.PaymentMethod{})
	db.AutoMigrate(&domain.WalletHistory{})
	db.AutoMigrate(&domain.Coupons{})
	db.AutoMigrate(&domain.UsedCoupon{})
	db.AutoMigrate(&domain.ProductOffer{})
	db.AutoMigrate(&domain.CategoryOffer{})
	db.AutoMigrate(&domain.Referral{})
	db.AutoMigrate(&domain.WishList{})
	db.AutoMigrate(&domain.Image{})
	db.AutoMigrate(&domain.Wallet{})
	db.AutoMigrate(&domain.Banner{})
	CheckAndCreateAdmin(confg, db)
	SeedPaymentMethods(db)
	return DB, err
}

func SeedPaymentMethods(db *gorm.DB) {
	methods := []string{"Razorpay"}
	for i, name := range methods {
		var count int64
		db.Model(&domain.PaymentMethod{}).Where("payment_name = ?", name).Count(&count)
		if count == 0 {
			db.Create(&domain.PaymentMethod{
				Model:        gorm.Model{ID: uint(i + 1)},
				Payment_Name: name,
			})
		}
	}
}

func CheckAndCreateAdmin(config config.Config, db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count == 0 {
		password := config.AdminPassword
		hashPassword, err := helper.PasswordHash(password)
		if err != nil {
			return
		}
		admin := domain.User{
			Model:     gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Firstname: "Admin",
			Lastname:  "Soxenly",
			Email:     config.AdminEmail,
			Password:  hashPassword,
			Phone:     "",
			Blocked:   false,
			Isadmin:   true,
		}
		db.Create(&admin)
	}
}
