package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(name []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": name,
			"exp":      time.Now().Add(time.Hour * 12).Unix(),
		})

	tkString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tkString, nil
}

func VerifyToken(name string) error {
	token, err := jwt.Parse(name, func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("SECRET_KEY"), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
