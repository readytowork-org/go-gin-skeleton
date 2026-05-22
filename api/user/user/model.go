package user

import (
	"encoding/json"

	"boilerplate-api/database/dao"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CUser struct {
	dao.User
}

// MarshalJSON hides the password hash from API responses.
func (u CUser) MarshalJSON() ([]byte, error) {
	type alias dao.User
	return json.Marshal(struct {
		alias
		Password string `json:"password,omitempty"`
	}{alias: alias(u.User)})
}

// BeforeCreate Runs before inserting a row into table
func (u *CUser) BeforeCreate(db *gorm.DB) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}
