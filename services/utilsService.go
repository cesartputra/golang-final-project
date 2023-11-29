package services

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetPublicImageIDFromCloudinaryURL(url string) string {
	parts := strings.Split(url, "/")
	lastPart := parts[len(parts)-1]
	publicID := strings.Split(lastPart, ".")[0]
	return publicID
}
