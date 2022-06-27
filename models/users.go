package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	BinaryBase
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

// TableName gives table name of model
func (u User) TableName() string {
	return "users"
}

// BeforeCreate
func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	id, err := uuid.NewRandom()
	u.ID = BINARY16(id)
	return err
}

// ToMap convert User to map
func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":     u.Email,
		"full_name": u.FullName,
	}
}
