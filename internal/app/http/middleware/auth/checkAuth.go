package auth

import (
	"SarkorTelecom-testTask/internal/app/service/jwtmanager"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// CheckAuth - проверяет авторизацию пользователя по токену в куке
func CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Получаем куку SESSTOKEN из запроса
		tokenString, err := c.Cookie("SESSTOKEN")
		if err != nil {
			// Кука отсутствует, возвращаем ошибку
			log.Printf("Cookie not found in CheckAuth: %v", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized, cookie not found")
		}

		if !jwtmanager.IsTokenValid(tokenString.Value) {
			// Токен не валиден, возвращаем ошибку авторизации
			log.Printf("Token is not valid in CheckAuth: %v", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized, token is not valid")
		}

		return next(c)
	}
}
