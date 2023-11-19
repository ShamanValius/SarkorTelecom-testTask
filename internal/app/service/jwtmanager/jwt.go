package jwtmanager

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

// TODO: Это необходимо вынести в конфигурацию
// TODO: Переделать на использование DI
const (
	secretKey       = "your-secret-key" // Ключ используется для подписи токена
	tokenExpiration = time.Hour * 1     // Время жизни токена
)

// NewToken создает новый токен для пользователя с логином login и id userId в payload
func NewToken(login string, userId int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = login
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(tokenExpiration).Unix()

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", errors.New("error creating token")
	}

	return tokenString, nil
}

// IsTokenValid проверяет, что токен валиден
func IsTokenValid(tokenString string) bool {
	_, err := parseToken(tokenString)
	return err == nil
}

// GetPayload возвращает payload токена в виде map[string]interface{}
func GetPayload(tokenString string) (map[string]interface{}, error) {
	// Парсим токен
	claims, err := parseToken(tokenString)
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, errors.New("error getting payload")
	}

	return claims, nil
}

// GetUserID возвращает payload токена
func parseToken(tokenString string) (jwt.MapClaims, error) {
	// Парсим токен встроенной функцией Parse в golang-jwt
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	// Проверяем, что токен валиден
	if err != nil || !token.Valid {
		log.Printf("Token is not valid: %v", err)
		return nil, errors.New("token is not valid")
	}

	// Получаем payload токена
	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
