package sqlite

import (
	"SarkorTelecom-testTask/internal/app/DTO"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

type Storage struct {
	db *sql.DB
}

// New создает новое хранилище, подключаясь к базе данных SQLite по пути storagePath
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	db, err := sql.Open("sqlite3", storagePath) // Подключаемся к БД
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем таблицу, если ее еще нет
	stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS users
	(
    id INTEGER PRIMARY KEY,
    login TEXT NOT NULL UNIQUE CHECK(login != ''),
    password TEXT NOT NULL,
    name TEXT NOT NULL CHECK(name != ''),
    age INTEGER NOT NULL CHECK(age > 0)
	);
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
    CREATE TABLE IF NOT EXISTS users_phone
    (
        id INTEGER PRIMARY KEY,
        phone TEXT NOT NULL UNIQUE CHECK(phone != '' AND LENGTH(phone) <= 12),
        description TEXT NOT NULL CHECK(description != ''),
        is_fax INTEGER DEFAULT 0 CHECK(is_fax == 0 OR is_fax == 1),
        users_id INTEGER,
        FOREIGN KEY (users_id) REFERENCES users (id) ON DELETE CASCADE
    );
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// CreateUser создает нового пользователя в базе данных
func (s *Storage) CreateUser(user *DTO.UserDTO) error {
	q := "INSERT INTO users (login, password, name, age) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(q, user.Login, user.Password, user.Name, user.Age)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			// UserDTO с таким логином уже существует
			log.Printf("UserDTO already exists: %v", err)
			return errors.New("user already exists")
		} else if strings.Contains(err.Error(), "err: Login is empty") {
			// Пользователь не ввел логин или он пустой (проверка на стороне базы данных)
			log.Printf("UserDTO login is empty: %v", err)
			return errors.New("user login is empty")
		} else {
			// База данных не смогла создать пользователя по какой-то другой причине
			log.Printf("Error creating user: %v", err)
			return errors.New("error creating user")
		}
	}
	return nil
}

// GetUserByLogin ищет пользователя в базе данных по логину
func (s *Storage) GetUserByLogin(login string) (*DTO.UserDTO, error) {
	// Ищем пользователя по логину
	var user DTO.UserDTO
	q := "SELECT * FROM users WHERE login = ?"
	// Записываем в user данные из строки запроса
	err := s.db.QueryRow(q, login).Scan(&user.ID, &user.Login, &user.Password, &user.Name, &user.Age)
	if errors.Is(err, sql.ErrNoRows) {
		// Пользователь не найден
		log.Printf("UserDTO not found for login: %s", login)
		return nil, errors.New("user not found for login: " + login)
	} else if err != nil {
		// Произошла какая-то другая ошибка при поиске пользователя
		return nil, errors.New("error getting user by login: " + login)
	}
	return &user, nil
}

// GetUsersByName ищет пользователей в базе данных по имени
func (s *Storage) GetUsersByName(name string) ([]DTO.UserDTO, error) {
	// Ищем пользователей по имени
	q := "SELECT * FROM users WHERE name = ?"
	rows, err := s.db.Query(q, name)
	if err != nil {
		log.Printf("Error getting users by name: %v", err)
		return nil, errors.New("error getting users by name: " + name)
	}

	// Записываем в users данные из строк запроса
	var users []DTO.UserDTO
	for rows.Next() {
		var user DTO.UserDTO
		err := rows.Scan(&user.ID, &user.Login, &user.Password, &user.Name, &user.Age)
		if err != nil {
			log.Printf("Error reading user row: %v", err)
			return nil, errors.New("error searching users by name")
		}
		users = append(users, user)
	}

	// Проверяем, не было ли ошибок при чтении строк
	if err := rows.Err(); err != nil {
		log.Printf("Error reading rows user: %v", err)
		return nil, errors.New("error searching users by name")
	}

	return users, nil
}

// CreatePhoneNumber создает новый номер телефона в базе данных
func (s *Storage) CreatePhoneNumber(phoneNumber *DTO.UserPhoneDTO) error {
	q := "INSERT INTO users_phone (phone, description, is_fax, users_id) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(q, phoneNumber.Phone, phoneNumber.Description, phoneNumber.IsFax, phoneNumber.UserId)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			// Телефон с таким номером уже существует
			log.Printf("Phone number already exists: %v", err)
			return errors.New("phone number already exists")
		} else if strings.Contains(err.Error(), "err: Phone is empty") {
			// Номер телефона не введен или он пустой (проверка на стороне базы данных)
			log.Printf("Phone number is empty: %v", err)
			return errors.New("phone number is empty")
		} else {
			// База данных не смогла создать номер телефона по какой-то другой причине
			log.Printf("Error creating phone number: %v", err)
			return errors.New("error creating phone number")
		}
	}
	return nil
}

// GetPhoneNumbers возвращает все номера телефонов из базы данных по запросу номера
func (s *Storage) GetPhoneNumbers(query string) ([]DTO.UserPhoneDTO, error) {
	// Ищем номера телефонов по запросу номера
	q := "SELECT * FROM users_phone WHERE phone LIKE '%' || ? || '%'"
	rows, err := s.db.Query(q, query)
	if err != nil {
		log.Printf("Error getting phone numbers: %v", err)
		return nil, errors.New("error getting phone numbers by query: " + query)
	}

	// Записываем в phones данные из строк запроса
	var phones []DTO.UserPhoneDTO
	for rows.Next() {
		var phone DTO.UserPhoneDTO
		err = rows.Scan(&phone.ID, &phone.Phone, &phone.Description, &phone.IsFax, &phone.UserId)
		if err != nil {
			log.Printf("Error reading phone row: %v", err)
			return nil, errors.New("error searching phone numbers by query")
		}
		phones = append(phones, phone)
	}

	// Проверяем, не было ли ошибок при чтении строк
	if err := rows.Err(); err != nil {
		log.Printf("Error reading rows phone: %v", err)
		return nil, errors.New("error searching phone numbers by query")
	}

	return phones, nil
}

// UpdatePhoneNumber обновляет данные номера телефона в базе данных
func (s *Storage) UpdatePhoneNumber(phoneNumber *DTO.UserPhoneDTO) error {
	q := "REPLACE INTO users_phone (id, phone, description, is_fax, users_id) VALUES (?, ?, ?, ?, ?)"
	_, err := s.db.Exec(q, phoneNumber.ID, phoneNumber.Phone, phoneNumber.Description, phoneNumber.IsFax, phoneNumber.UserId)
	if err != nil {
		log.Printf("Error updating phone: %v", err)
		return errors.New("error updating phone")
	}
	return nil
}

// DeletePhoneNumber удаляет номер телефона из базы данных
func (s *Storage) DeletePhoneNumber(phoneID int) error {
	q := "DELETE FROM users_phone WHERE id = ?"
	_, err := s.db.Exec(q, phoneID)
	if err != nil {
		log.Printf("Error deleting phone: %v", err)
		return errors.New("error deleting phone")
	}
	return nil
}

// Close закрывает соединение с базой данных
func (s *Storage) Close() error {
	return s.db.Close()
}
