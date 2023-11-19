package main

import (
	"SarkorTelecom-testTask/internal/app"
	"SarkorTelecom-testTask/internal/app/service/passHash/bcrypt"
	"SarkorTelecom-testTask/internal/storage/sqlite"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Загрузка переменных конфигурации из файла .env
	err := godotenv.Load("internal/config/.env")
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	// Получение переменных конфигурации
	serverPort := os.Getenv("HTTP_SERVER_PORT")
	dbPath := os.Getenv("STORAGE_PATH")
	hasher := bcrypt.New(10)

	// Создание экземпляра хранилища SQLite
	storage, err := sqlite.New(dbPath)
	if err != nil {
		panic("Error creating SQLite storage: " + err.Error())
	}
	// Закрытие соединения с БД при завершении работы приложения
	defer storage.Close()

	// Создание экземпляра приложения и передача ему хранилища и хэшера
	myApp, err := app.New(storage, hasher)
	if err != nil {
		panic("Error creating app: " + err.Error())
	}

	// Запуск приложения
	log.Fatal(myApp.Run(":" + serverPort))
}
