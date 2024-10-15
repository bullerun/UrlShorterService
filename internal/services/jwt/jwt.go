package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = "nrilsvjn"

func CreateJwt(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

func DecodeJWT(jwtString string) (int64, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := int64(claims["id"].(float64))
		expiration := int64(claims["exp"].(float64))
		if expiration > time.Now().Unix() {
			return id, nil
		} else {
			return 0, ErrTokenExpired
		}
	}
	return 0, ErrTokenInvalid
}
