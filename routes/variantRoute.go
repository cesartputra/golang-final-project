package routes

import (
	"golang-final-project/controllers"
	"golang-final-project/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VariantRoutes(route *gin.Engine, db *gorm.DB) {
	route.POST("/api/products/variants", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.CreateVariant(c, db)
	})
	route.GET("/api/products/variants", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.GetAllVariantsWithPagination(c, db)
	})
	route.GET("/api/products/variants/:id", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.GetVariantByID(c, db)
	})
	route.PUT("/api/products/variants/:id", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.UpdateVariantByID(c, db)
	})
	route.DELETE("/api/products/variants/:id", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.DeleteVariantByID(c, db)
	})
}
