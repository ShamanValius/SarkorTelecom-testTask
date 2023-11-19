package routes

import (
	"SarkorTelecom-testTask/internal/app/http/handlers"
	"SarkorTelecom-testTask/internal/app/http/middleware/auth"
	"github.com/labstack/echo/v4"
)

// InitializeHTTPRoutes - инициализация всех роутов
func InitializeHTTPRoutes(e *echo.Echo, h *handlers.HandlerSet) {

	e.POST("/user/register", h.UserHandlers.RegisterUser)
	e.POST("/user/auth", h.UserHandlers.AuthenticateUser)
	e.GET("/user/:name", h.UserHandlers.GetUserByName, auth.CheckAuth)

	phone := e.Group("/user/phone")
	phone.Use(auth.CheckAuth)
	phone.POST("", h.PhoneHandlers.CreateUserPhone)
	phone.GET("", h.PhoneHandlers.GetUserPhoneByNumbers)
	phone.PUT("", h.PhoneHandlers.UpdateUserPhone)
	phone.DELETE("/:phone_id", h.PhoneHandlers.DeleteUserPhoneByID)

}
