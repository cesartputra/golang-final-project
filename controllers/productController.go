package controllers

import (
	"golang-final-project/models"
	"golang-final-project/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateProduct(c *gin.Context, db *gorm.DB) {
	var request struct {
		Name     string `json:"name" binding:"required"`
		ImageUrl string `json:"imageUrl" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Admin ID
	adminIDString, err := services.ExtractAdminID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error get Admin ID"})
		return
	}

	adminID, err := uuid.Parse(adminIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Admin ID"})
		return
	}

	product := models.Product{
		Name:     request.Name,
		ImageUrl: request.ImageUrl,
		AdminID:  adminID,
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": r})
		}
	}()

	if err := services.CreateProduct(db, &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func GetAllProductsWithPagination(c *gin.Context, db *gorm.DB) {
	page := 1
	pageSize := 10

	if pageParam, exists := c.GetQuery("page"); exists {
		parsedPage, err := strconv.Atoi(pageParam)
		if err == nil {
			page = parsedPage
		}
	}

	if pageSizeParam, exists := c.GetQuery("pageSize"); exists {
		parsedPageSize, err := strconv.Atoi(pageSizeParam)
		if err == nil {
			pageSize = parsedPageSize
		}
	}

	searchName := c.Query("search")
	products, err := services.GetAllProductsWithPaginationAndSearch(db, page, pageSize, searchName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context, db *gorm.DB) {
	idString := c.Param("id")

	if idString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID not provided"})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	product, err := services.GetProductByID(db, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateProductByID(c *gin.Context, db *gorm.DB) {
	var request struct {
		Name     string `json:"name" binding:"required"`
		ImageUrl string `json:"imageUrl" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idString := c.Param("id")

	if idString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID not provided"})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	// Get Admin ID
	adminIDString, err := services.ExtractAdminID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error get Admin ID"})
		return
	}

	adminID, err := uuid.Parse(adminIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Admin ID"})
		return
	}

	existingProduct, err := services.GetProductByID(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if existingProduct.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized with this Admin ID"})
		return
	}

	existingProduct.Name = request.Name
	existingProduct.ImageUrl = request.ImageUrl

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": r})
		}
	}()

	if err := services.UpdateProductByID(db, id, existingProduct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func DeleteProductByID(c *gin.Context, db *gorm.DB) {
	idString := c.Param("id")

	if idString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID not provided"})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	product, err := services.GetProductByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get Admin ID
	adminIDString, err := services.ExtractAdminID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error get Admin ID"})
		return
	}

	adminID, err := uuid.Parse(adminIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Admin ID"})
		return
	}

	if product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized with this Admin ID"})
		return
	}

	variants, _ := services.GetVariantsByProductID(db, id)

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": r})
		}
	}()

	if variants != nil {
		err := services.DeleteVariantsByProductID(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete Variant"})
			panic(err)
		}
	}

	err = services.DeleteProductByID(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
