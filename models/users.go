package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	BinaryBase
	FirebaseUID string `json:"firebase_uid"`
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	FullName    string `json:"full_name" validate:"required"`
	Address     string `json:"address" validate:"required"`
}

// TableName gives table name of model
func (m *User) TableName() string {
	return "user"
}

// ToMap convert User to map
func (m User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":     m.Email,
		"username":  m.Username,
		"phone":     m.Phone,
		"full_name": m.FullName,
		"address":   m.Address,
	}
}

// Runs before inserting a row into table
func (m *User) BeforeCreate(db *gorm.DB) error {
	id, err := uuid.NewRandom()

	m.ID = BINARY16(id)
	fmt.Println(m.ID, "--id----")
	b, err := json.MarshalIndent(m, "", "")
	if err == nil {
		fmt.Println(string(b))
	}
	return err
}
