package dtos

import "boilerplate-api/models"

// Request body data to create user
type CreateUserRequestData struct {
	Email           string `json:"email" validate:"required,email"`
	FullName        string `json:"full_name" validate:"required"`
	Phone           string `json:"phone"  validate:"required,phone"`
	Gender          string `json:"gender" validate:"required,gender"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

func (u *CreateUserRequestData) GetUser() models.User {
	return models.User{
		Email:    u.Email,
		FullName: u.FullName,
		Phone:    u.Phone,
		Gender:   u.Gender,
		Password: u.Password,
	}
}