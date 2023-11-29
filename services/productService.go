package services

import (
	"golang-final-project/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB, product *models.Product) error {
	return db.Create(&product).Error
}

func GetAllProductsWithPaginationAndSearch(db *gorm.DB, page, pageSize int, searchName string) ([]models.Product, error) {
	var products []models.Product

	offset := (page - 1) * pageSize

	query := db.Offset(offset).Limit(pageSize)

	if searchName != "" {
		query = query.Where("name LIKE ?", "%"+searchName+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func GetAllProducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product

	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func GetProductByID(db *gorm.DB, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := db.First(&product, id).Error
	return &product, err
}

func UpdateProductByID(db *gorm.DB, id uuid.UUID, product *models.Product) error {
	return db.Model(&models.Product{}).Where("id = ?", id).Updates(product).Error
}

func DeleteProductByID(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&models.Product{}, id).Error
}
