package handlers

import (
	"SarkorTelecom-testTask/internal/app/DTO"
	"SarkorTelecom-testTask/internal/app/service/jwtmanager"
	"SarkorTelecom-testTask/internal/storage"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

// HandlPhone обрабатывает запросы, связанные с UserPhoneDTO
type HandlPhone struct {
	s storage.Storage
}

// NewHandlerPhone создает новый обработчик номеров телефонов с хранилищем s
func NewHandlerPhone(s storage.Storage) *HandlPhone {
	return &HandlPhone{s: s}
}

// CreateUserPhone обрабатывает POST /user/phone и сохраняет номер телефона в базе данных,
// возвращая StatusCreated в случае успеха
func (h HandlPhone) CreateUserPhone(c echo.Context) error {
	var requestData map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&requestData); err != nil {
		log.Printf("failed to decode JSON: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	tokenString, err := c.Cookie("SESSTOKEN")
	if err != nil {
		log.Printf("failed to read cookie: %s", err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "read cookie failed"})
	}

	// Получаем текущего пользователя из контекста
	payload, err := jwtmanager.GetPayload(tokenString.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	// Получаем ID пользователя из токена
	userID := int(payload["user_id"].(float64))

	// Создаем объект UserPhoneDTO
	phoneData := DTO.UserPhoneDTO{
		Phone:       requestData["phone"].(string),
		Description: requestData["description"].(string),
		IsFax:       requestData["is_fax"] == "true", // Преобразование строки в булев тип
		UserId:      userID,
	}

	// Добавление номера пользователя
	err = h.s.CreatePhoneNumber(&phoneData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "phone number created successfully"})
}

// GetUserPhoneByNumber возвращает все номера телефонов из базы данных по запросу номера,
// возвращая StatusOK в случае успеха и JSON-объект с номерами телефонов  и их данными
func (h HandlPhone) GetUserPhoneByNumbers(c echo.Context) error {
	// Получаем значение параметра запроса
	query := c.QueryParam("q")
	if query == "" {
		log.Printf("invalid query parameter")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid query parameter"})
	}

	// Получаем номера телефонов из хранилища
	phones, err := h.s.GetPhoneNumbers(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Преобразуем в нужный формат JSON
	var response []map[string]interface{}
	for _, phone := range phones {
		response = append(response, map[string]interface{}{
			"user_id":     phone.UserId,
			"phone":       phone.Phone,
			"description": phone.Description,
			"is_fax":      phone.IsFax,
		})
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateUserPhone обновляет номер телефона в базе данных, возвращая StatusOK в случае успеха
func (h HandlPhone) UpdateUserPhone(c echo.Context) error {
	// Получаем текущего пользователя из Cookie
	tokenString, err := c.Cookie("SESSTOKEN")
	if err != nil {
		log.Printf("failed to read cookie: %s", err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unauthorized"})
	}

	// Получаем ID пользователя из токена
	payload, err := jwtmanager.GetPayload(tokenString.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	userID := int(payload["user_id"].(float64))

	// Получаем ID номера телефона из параметра запроса
	phoneID, err := strconv.Atoi(c.QueryParam("phone_id"))
	if err != nil {
		log.Printf("invalid phone ID in QueryParam: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid phone ID in QueryParam"})
	}

	phone := c.QueryParam("phone")
	if phone == "" {
		log.Printf("invalid phone in QueryParam: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid phone in QueryParam"})
	}

	description := c.QueryParam("description")
	if description == "" {
		log.Printf("invalid description in QueryParam: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid description in QueryParam"})
	}

	isFax, err := strconv.ParseBool(c.QueryParam("is_fax"))
	if err != nil {
		log.Printf("invalid is_fax in QueryParam %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid is_fax in QueryParam"})
	}

	// Создаем объект UserPhoneDTO
	updatedPhone := &DTO.UserPhoneDTO{
		ID:          phoneID,
		Phone:       phone,
		IsFax:       isFax,
		Description: description,
		UserId:      userID,
	}

	// Обновляем данные в базе данных
	err = h.s.UpdatePhoneNumber(updatedPhone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "phone number updated successfully"})
}

// DeleteUserPhoneByID удаляет номер телефона из базы данных по его ID, возвращая StatusOK в случае успеха
func (h HandlPhone) DeleteUserPhoneByID(c echo.Context) error {
	phoneID, err := strconv.Atoi(c.Param("phone_id"))
	if err != nil {
		log.Printf("invalid phone_id in QueryParam")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid phone_id in QueryParam"})
	}

	// Удаляем номер телефона из базы данных
	err = h.s.DeletePhoneNumber(phoneID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "phone number deleted successfully"})
}
