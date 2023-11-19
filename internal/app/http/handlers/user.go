package handlers

import (
	"SarkorTelecom-testTask/internal/app/DTO"
	"SarkorTelecom-testTask/internal/app/service/jwtmanager"
	"SarkorTelecom-testTask/internal/app/service/passHash"
	"SarkorTelecom-testTask/internal/storage"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

// HandlUser обрабатывает запросы, связанные с UserDTO
type HandlUser struct {
	s storage.Storage
	h passHash.Hasher
}

// NewUserHandler создает новый обработчик пользователей с хранилищем s и сервисом хеширования h
func NewUserHandler(s storage.Storage, h passHash.Hasher) *HandlUser {
	return &HandlUser{s: s, h: h}
}

// RegisterUser обрабатывает POST /user и сохраняет пользователя в базе данных, возвращая StatusCreated в случае успеха
func (h HandlUser) RegisterUser(c echo.Context) error {
	user := new(DTO.UserDTO)
	// Привязываем данные из формы к структуре UserDTO
	if err := c.Bind(user); err != nil {
		log.Printf("Invalid data from form: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid data from form"})
	}

	// Хешируем пароль с помощью сервиса passHash и той реализации, которая передана в h
	hashedPassword, err := h.h.Password(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "error hashing password"})
	}

	// Перезаписываем пароль в структуре UserDTO на хешированный
	user.Password = string(hashedPassword)

	// Сохраняем пользователя в базе данных через то хранилище, которое передано в h
	err = h.s.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

// AuthenticateUser обрабатывает POST /user/auth и аутентифицирует пользователя, возвращая токен в случае успеха
func (h HandlUser) AuthenticateUser(c echo.Context) error {
	//TODO: Обработчик получился слишком большим, его нужно выделить в отдельный сервис AuthenticateUser

	// Получаем данные из JSON-запроса
	var requestData map[string]string
	if err := json.NewDecoder(c.Request().Body).Decode(&requestData); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	// Проверяем, что в requestData есть поля login и password и они не пустые
	login, loginExists := requestData["login"]
	password, passwordExists := requestData["password"]

	if !loginExists || !passwordExists || login == "" || password == "" {
		log.Printf("Missing or empty login or password")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing or empty login or password"})
	}

	// Ищем пользователя по логину в базе данных
	user, err := h.s.GetUserByLogin(requestData["login"])
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Сравниваем пароли с помощью сервиса passHash и той реализации, которая передана в h
	err = h.h.Compare([]byte(user.Password), requestData["password"])
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Создаем новый токен
	tokenString, err := jwtmanager.NewToken(user.Login, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Устанавливаем куку в ответе
	cookie := new(http.Cookie)
	cookie.Name = "SESSTOKEN"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(time.Hour * 1)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"message": "Authentication successful"})

}

// GetUserByName обрабатывает GET /user/:name и возвращает пользователей с именем name в формате JSON
func (h HandlUser) GetUserByName(c echo.Context) error {
	// Получаем имя пользователя из URL
	n := c.Param("name")
	if n == "" {
		log.Printf("Name is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is empty"})
	}

	// Ищем пользователей по имени в базе данных через то хранилище, которое передано в h
	users, err := h.s.GetUsersByName(n)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Преобразуем пользователей в формат JSON
	var userJSON []map[string]interface{}
	for _, user := range users {
		userMap := map[string]interface{}{
			"id":   user.ID,
			"name": user.Name,
			"age":  user.Age,
		}
		userJSON = append(userJSON, userMap)
	}

	return c.JSON(http.StatusOK, userJSON)

}
