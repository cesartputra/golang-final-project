package services

import (
	"golang-final-project/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateVariant(db *gorm.DB, variant *models.Variant) error {
	return db.Create(&variant).Error
}

func GetAllVariantsWithPaginationAndSearch(db *gorm.DB, page, pageSize int, searchName string) ([]models.Variant, error) {
	var variants []models.Variant

	offset := (page - 1) * pageSize

	query := db.Offset(offset).Limit(pageSize)

	if searchName != "" {
		query = query.Where("variant_name LIKE ?", "%"+searchName+"%")
	}

	if err := query.Find(&variants).Error; err != nil {
		return nil, err
	}

	return variants, nil
}

func GetVariantByID(db *gorm.DB, id uuid.UUID) (*models.Variant, error) {
	var variant models.Variant
	err := db.First(&variant, id).Error
	return &variant, err
}

func GetVariantsByProductID(db *gorm.DB, productID uuid.UUID) (*models.Variant, error) {
	var variant models.Variant
	err := db.Where("product_id = ?", productID).First(&variant).Error
	return &variant, err
}

func UpdateVariantByID(db *gorm.DB, id uuid.UUID, variant *models.Variant) error {
	return db.Model(&models.Variant{}).Where("id = ?", id).Updates(variant).Error
}

func DeleteVariantsByProductID(db *gorm.DB, productID uuid.UUID) error {
	res := db.Where("product_id = ?", productID).Delete(models.Variant{}).Error

	if res != nil {
		return res
	}

	return nil
}

func DeleteVariantByID(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&models.Variant{}, id).Error
}
