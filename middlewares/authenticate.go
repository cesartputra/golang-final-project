package middlewares

import (
	"golang-final-project/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		if bearerToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bearer Token not provided"})
			c.Abort()
			return
		}

		if !strings.Contains(bearerToken, "Bearer") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bearer not provided"})
			c.Abort()
			return
		}

		strArr := strings.Split(bearerToken, " ")
		tokenString := strArr[1]

		if tokenString == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
			c.Abort()
			return
		}

		claims := &services.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return services.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("adminID", claims.ID)
		c.Next()
	}
}
