package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: add secret to env
var secretKey = []byte("secret-key")

func CreateToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": id,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims format")
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		return 0, fmt.Errorf("userID not found in token")
	}

	return int(userID), nil
}
