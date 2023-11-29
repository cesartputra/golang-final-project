package main

import (
	database "golang-final-project/dabatase"
	"golang-final-project/routes"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
)

func main() {
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	if err != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", err)
	}

	db := database.ConnectDB()
	r := gin.Default()

	routes.AuthRoute(r, db)
	routes.ProductRoute(r, db, cld)
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
