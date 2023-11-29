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

	port := envPortOr("3000")
	r.Run(port)
}

func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
}
