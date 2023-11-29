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

func CreateVariant(c *gin.Context, db *gorm.DB) {
	var request struct {
		VariantName string    `json:"variantName" binding:"required"`
		Quantity    int       `json:"quantity" binding:"required"`
		ProductID   uuid.UUID `json:"productID" binding:"required"`
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

	product, _ := services.GetProductByID(db, request.ProductID)

	if product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized with this Admin ID"})
		return
	}

	variant := models.Variant{
		VariantName: request.VariantName,
		Quantity:    request.Quantity,
		ProductID:   request.ProductID,
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": r})
		}
	}()

	if err := services.CreateVariant(db, &variant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Variant created successfully"})
}

func GetAllVariantsWithPagination(c *gin.Context, db *gorm.DB) {
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
	variants, err := services.GetAllVariantsWithPaginationAndSearch(db, page, pageSize, searchName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variants)
}

func GetVariantByID(c *gin.Context, db *gorm.DB) {
	idString := c.Param("id")

	if idString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Variant ID not provided"})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Variant ID"})
		return
	}

	variant, err := services.GetVariantByID(db, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

func UpdateVariantByID(c *gin.Context, db *gorm.DB) {
	var request struct {
		VariantName string `json:"variantName" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idString := c.Param("id")

	if idString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Variant ID not provided"})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Variant ID"})
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

	existingVariant, err := services.GetVariantByID(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	product, _ := services.GetProductByID(db, existingVariant.ProductID)

	if product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized with this Admin ID"})
		return
	}

	existingVariant.VariantName = request.VariantName
	existingVariant.Quantity = request.Quantity

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": r})
		}
	}()

	if err := services.UpdateVariantByID(db, id, existingVariant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Variant updated successfully"})
}

func DeleteVariantByID(c *gin.Context, db *gorm.DB) {
	idString := c.Param("id")

	if idString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Variant ID not provided"})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Variant ID"})
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

	variant, err := services.GetVariantByID(db, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Variant not found"})
		return
	}

	product, _ := services.GetProductByID(db, variant.ProductID)

	if product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized with this Admin ID"})
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": r})
		}
	}()

	err = services.DeleteVariantByID(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
