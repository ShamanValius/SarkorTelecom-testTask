package DTO

// UserPhoneDTO модель данных телефона пользователя
type UserPhoneDTO struct {
	ID          int    `json:"id"`
	Phone       string `json:"number" form:"number" validate:"required, min=12, max=12"`
	Description string `json:"desc" form:"desc" validate:"required"`
	IsFax       bool   `json:"is_fax" form:"is_fax" validate:"required"`
	UserId      int    `json:"user_id" form:"user_id" validate:"required"`
}
