package models

type User struct {
	Base
	Email    string `json:"email" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
}

type UserProfile struct {
	Base
	UserId  int64  `json:"user_id" validate:"required"`
	Address string `json:"address" validate:"required"`
	Contact string `json:"contact" validate:"required"`
}

type UserProfileDetail struct {
	UserProfile
	UserDetail User
}

// TableName gives table name of model
func (m User) TableName() string {
	return "users"
}

func (p UserProfile) TableName() string {
	return "user_profile"
}

// ToMap convert User to map
func (m User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":     m.Email,
		"full_name": m.FullName,
	}
}

func (p UserProfile) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"address": p.Address,
		"contact": p.Contact,
		"user_id": p.UserId,
	}
}
