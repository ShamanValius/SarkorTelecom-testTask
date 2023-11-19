package handlers

import (
	"SarkorTelecom-testTask/internal/app/service/passHash"
	"SarkorTelecom-testTask/internal/storage"
)

// HandlerSet - набор обработчиков
type HandlerSet struct {
	UserHandlers  *HandlUser
	PhoneHandlers *HandlPhone
}

// NewHandlerSet инциализирует набор обработчиков
func NewHandlerSet(storage storage.Storage, hasher passHash.Hasher) *HandlerSet {
	return &HandlerSet{
		UserHandlers:  NewUserHandler(storage, hasher),
		PhoneHandlers: NewHandlerPhone(storage),
	}
}
