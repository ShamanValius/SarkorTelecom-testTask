package bcrypt

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Hash реализует интерфейс Hasher, используя bcrypt
type Hash struct {
	cost int
}

// New возвращает новый bcrypt хэшер с указанным cost, который является стоимостью хэширования
func New(cost int) *Hash {
	return &Hash{cost: cost}
}

// Password хэширует пароль
func (h *Hash) Password(password string) ([]byte, error) {
	// Хэшируем пароль с помощью bcrypt встроенной в Go
	b, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, err
	}
	return b, nil
}

// Compare сравнивает хэш с паролем
func (h *Hash) Compare(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	switch {
	// Хэш не соответствует паролю
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		log.Printf("ErrNotValidPassword: %v", err)
		return errors.New("ErrNotValidPassword")
	// Не удалось сравнить пароли
	case err != nil:
		log.Printf("Error comparing passwords: %v", err)
		return errors.New("error comparing passwords")
	}

	return nil
}
