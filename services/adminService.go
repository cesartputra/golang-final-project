package services

import (
	"golang-final-project/models"

	"gorm.io/gorm"
)

func CreateAdmin(db *gorm.DB, admin *models.Admin) error {
	return db.Create(&admin).Error
}

func GetAdminByEmail(db *gorm.DB, email string) (*models.Admin, error) {
	var admin models.Admin
	err := db.Where("email = ?", email).First(&admin).Error

	return &admin, err
}
