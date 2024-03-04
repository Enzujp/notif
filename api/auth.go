package api

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/enzujp/notif/pkg/models"
)

var mySecretKey = os.Getenv("SECRET_KEY")
var secretKey = []byte(mySecretKey)

func GenerateToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	// set expiry token to 24 hours
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}

func ExtractUserIDFromToken(tokenString string) (uint, error){
	token, err := ParseToken(tokenString)
	if err != nil{
		return 0, nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}
	return uint(userID), nil
}
