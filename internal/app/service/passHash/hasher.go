package passHash

// Сервис хеширования, отвечающий за работу с паролями.
type Hasher interface {
	// Password хеширует пароль.
	Password(password string) ([]byte, error)
	// Compare сравнивает хешированный пароль с паролем в открытом виде.
	Compare(hashedPassword []byte, password string) error
}
