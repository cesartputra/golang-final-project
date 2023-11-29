package routes

import (
	"golang-final-project/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoute(route *gin.Engine, db *gorm.DB) {
	route.POST("/api/auth/register", func(c *gin.Context) {
		controllers.Register(c, db)
	})

	route.POST("api/auth/login", func(c *gin.Context) {
		controllers.Login(c, db)
	})
}
