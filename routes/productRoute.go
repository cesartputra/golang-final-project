package routes

import (
	"golang-final-project/controllers"
	"golang-final-project/middlewares"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductRoute(route *gin.Engine, db *gorm.DB, cld *cloudinary.Cloudinary) {
	route.POST("/api/products", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.CreateProduct(c, db, cld)
	})
	route.GET("/api/products", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.GetAllProductsWithPagination(c, db)
	})
	route.GET("/api/products/:id", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.GetProductByID(c, db)
	})
	route.PUT("/api/products/:id", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.UpdateProductByID(c, db)
	})
	route.DELETE("/api/products/:id", middlewares.AuthenticateJWT(), func(c *gin.Context) {
		controllers.DeleteProductByID(c, db)
	})
}
