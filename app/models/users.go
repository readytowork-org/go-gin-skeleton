package models

type User struct {
	Base
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone"  validate:"required,phone"`
	Gender   string `json:"gender" validate:"required,gender"`
}

// TableName gives table name of model
func (u User) TableName() string {
	return "users"
}

// ToMap convert User to map
func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":     u.Email,
		"full_name": u.FullName,
		"phone":     u.Phone,
		"gender":    u.Gender,
	}
}
