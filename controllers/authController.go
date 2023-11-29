package controllers

import (
	"golang-final-project/models"
	"golang-final-project/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(c *gin.Context, db *gorm.DB) {
	var request struct {
		Name          string `json:"name" binding:"required"`
		Email         string `json:"email" binding:"required"`
		Password      string `json:"password" binding:"required"`
		PasswordCheck string `json:"passwordCheck" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Password != request.PasswordCheck {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password not match"})
		return
	}

	encryptedPassword, err := services.EncryptPassword(request.Password)

	if err != nil {
		return
	}

	admin := models.Admin{
		Name:     request.Name,
		Email:    request.Email,
		Password: encryptedPassword,
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": r})
		}
	}()

	if err := services.CreateAdmin(db, &admin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Admin created successfully"})
}

func Login(c *gin.Context, db *gorm.DB) {
	var request struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := services.GetAdminByEmail(db, request.Email)

	if err != nil {
		return
	}

	if err := services.CheckPassword(request.Password, admin.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.GenerateJWT(admin.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Generate token failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"adminID": admin.ID,
		"token":   token,
	})
}
