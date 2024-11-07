package my_jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
)

func GetClaims(userId string, walletId string) string {

	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"walletId": walletId,
	})
	signedToken, err := token.SignedString(secret)
	if err != nil {
		fmt.Println(err)
	}
	return signedToken
}

func ParseToken(receivedToken string) (uuid.UUID, uuid.UUID, error) {

	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(receivedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		fmt.Printf("ParseToken err: %v\n", err)
		return uuid.Nil, uuid.Nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, parseErr := uuid.Parse(claims["userId"].(string))
		if parseErr != nil {
			return uuid.Nil, uuid.Nil, parseErr
		}
		walletId, parseErr := uuid.Parse(claims["walletId"].(string))
		if parseErr != nil {
			return uuid.Nil, uuid.Nil, parseErr
		}

		return userId, walletId, nil
	}
	return uuid.Nil, uuid.Nil, err
}
