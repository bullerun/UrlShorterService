package user_service

import (
	"UrlShorterService/internal/entity"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func CreateUser(username, email, password string) (*entity.User, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return nil, err
	}
	pwd := hashPassword(password, salt)
	return &entity.User{
		Name:     username,
		Password: pwd,
		Salt:     salt,
		Email:    email,
	}, nil
}
func generateSalt(saltSize int) (string, error) {
	var salt = make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}
func hashPassword(password string, salt string) string {
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()
	passwordBytes = append(passwordBytes, salt...)
	sha512Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}
func IsCorrectPassword(password, salt, hash string) bool {
	return hashPassword(password, salt) == hash
}
