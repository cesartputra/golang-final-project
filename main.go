package main

import (
	database "golang-final-project/dabatase"
	"golang-final-project/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.ConnectDB()

	r := gin.Default()

	routes.AuthRoute(r, db)
	routes.ProductRoute(r, db)
	routes.VariantRoutes(r, db)

	r.Run(os.Getenv("PORT"))
}
