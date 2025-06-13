package db

import (
	"me-pague/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	database, err := gorm.Open(sqlite.Open("payments.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database.AutoMigrate(&models.User{}, &models.Payment{}, &models.Billing{})
	DB = database
}
