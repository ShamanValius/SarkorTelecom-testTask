package DTO

// UserDTO модель данных пользователя
type UserDTO struct {
	ID       int    `json:"id"`
	Login    string `json:"login" form:"login" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Name     string `json:"name" form:"name" validate:"required"`
	Age      int    `json:"age" form:"age" validate:"required,numeric, min=1, max=150"`
}
