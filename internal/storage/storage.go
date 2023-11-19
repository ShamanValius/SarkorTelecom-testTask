package storage

import "SarkorTelecom-testTask/internal/app/DTO"

type Storage interface {
	CreateUser(user *DTO.UserDTO) error
	GetUserByLogin(login string) (*DTO.UserDTO, error)
	GetUsersByName(name string) ([]DTO.UserDTO, error)
	CreatePhoneNumber(phoneNumber *DTO.UserPhoneDTO) error
	GetPhoneNumbers(query string) ([]DTO.UserPhoneDTO, error)
	UpdatePhoneNumber(phoneNumber *DTO.UserPhoneDTO) error
	DeletePhoneNumber(phoneID int) error
}
