package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/morzisorn/practicum_final_project/config"
)

func GetToken(password string) (string, error) {
	if password != config.RealPass {
		return "", fmt.Errorf("incorrect password")
	}

	secretKey := []byte(config.RealPass)

	hash := sha256.Sum256([]byte(config.RealPass))
	hashString := hex.EncodeToString(hash[:])

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hashString,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return "", err
	}
	return tokenString, nil
}

func isTokenValid(tokenString string) bool {
	realPass := os.Getenv("TODO_PASSWORD")
	secretKey := []byte(realPass)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false
	}
	hash := sha256.Sum256([]byte(realPass))
	if token.Claims.(jwt.MapClaims)["hash"] != hex.EncodeToString(hash[:]) {
		return false
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return false
	}

	return true
}

func Auth(c *fiber.Ctx) error {
	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) == 0 {
		fmt.Println("Password not set")
		return c.Status(fiber.StatusInternalServerError).SendString("Server misconfiguration: password is not set")
	}

	token := c.Cookies("token")
	if token == "" {
		fmt.Println("Token missing")
		return c.Status(fiber.StatusUnauthorized).SendString("Token is missing")
	}

	if !isTokenValid(token) {
		fmt.Println("Invalid token")
		return c.Status(fiber.StatusUnauthorized).SendString("Authentication required")
	}

	return c.Next()
}
