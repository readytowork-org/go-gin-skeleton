package models

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Base
	Email    string `gorm:"email" json:"email" validate:"required,email"`
	FullName string `gorm:"full_name" json:"full_name"`
	Phone    string `gorm:"phone" json:"phone"  validate:"phone"`
	Gender   string `gorm:"gender" json:"gender" validate:"gender"`
	Password string `gorm:"password" json:"password" validate:"required"`
}

type OAuthUser struct {
	Base
	OAuthId       string `json:"id"`
	Name          string `json:"name"`
	Email         string `gorm:"email" json:"email" validate:"required,email"`
	VerifiedEmail bool   `gorm:"-" json:"verified_email" validate:"required"`
}

// TableName gives table name of model
func (u *User) TableName() string {
	return "users"
}

// ToMap convert User to map
func (u *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":     u.Email,
		"full_name": u.FullName,
		"phone":     u.Phone,
		"gender":    u.Gender,
	}
}

// BeforeCreate Runs before inserting a row into table
func (u *User) BeforeCreate(db *gorm.DB) error {
	var Zap *zap.SugaredLogger
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(password)
	if err != nil {
		Zap.Error("Error decrypting plain password to hash", err.Error())
		return err
	}
	return nil
}
