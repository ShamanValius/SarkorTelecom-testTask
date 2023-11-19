package app

import (
	"SarkorTelecom-testTask/internal/app/http/handlers"
	"SarkorTelecom-testTask/internal/app/http/routes"
	"SarkorTelecom-testTask/internal/app/service/passHash"
	"SarkorTelecom-testTask/internal/storage"
	"github.com/labstack/echo/v4"
)

// App представляет собой экземпляр приложения
type App struct {
	echo    *echo.Echo
	storage storage.Storage
}

// New создает экземпляр приложения
func New(storage storage.Storage, hasher passHash.Hasher) (*App, error) {

	// Создаем экземпляр фреймворка Echo
	e := echo.New()

	// Создаем структуру, которая содержит все обработчики запросов и передаем ей хранилище и хэшер
	handlerSet := handlers.NewHandlerSet(storage, hasher)

	// Инициализируем маршруты
	routes.InitializeHTTPRoutes(e, handlerSet)

	// Возвращаем экземпляр приложения
	return &App{
		echo:    e,
		storage: storage,
	}, nil
}

// Run запускает приложение на указанном порту
func (a *App) Run(adress string) error {
	return a.echo.Start(adress)
}
