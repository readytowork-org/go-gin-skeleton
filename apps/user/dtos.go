package user

import "boilerplate-api/helpers"

// CreateUserRequestData Request body data to create user
type CreateUserRequestData struct {
	User
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// GetUserResponse Dtos for User model
type GetUserResponse struct {
	helpers.Base
	Email    string `gorm:"email" json:"email"`
	FullName string `gorm:"full_name" json:"full_name"`
	Phone    string `gorm:"phone" json:"phone"`
	Gender   string `gorm:"gender" json:"gender"`
}
