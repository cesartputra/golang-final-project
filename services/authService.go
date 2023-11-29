package services

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	AdminID uuid.UUID `json:"adminID"`
	jwt.RegisteredClaims
}

func GenerateJWT(adminID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		AdminID: adminID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func ExtractAdminID(c *gin.Context) (string, error) {
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		adminID, ok := claims["adminID"].(string)
		if !ok {
			return "", fmt.Errorf("adminID not found in token")
		}
		return adminID, nil
	}

	return "", fmt.Errorf("invalid token")
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")

	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
