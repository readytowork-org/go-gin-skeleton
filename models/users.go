package models

type User struct {
	Base
	FirebaseUid    string `json:"firebase_uid"`
	NickName       string `json:"nick_name"`
	Email          string `json:"email"`
	FullName       string `json:"full_name"`
	Phone          string `json:"phone"`
	ResidentStatus string `json:"resident_status"`
	Age            string `json:"age"`
}

// TableName gives table name of model
func (m User) TableName() string {
	return "users"
}

// ToMap convert User to map
func (m User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"nick_name":       m.NickName,
		"email":           m.Email,
		"resident_status": m.ResidentStatus,
		"phone":           m.Phone,
		"full_name":       m.FullName,
		"create_at":       m.CreatedAt,
		"age":             m.Age,
		"firebase_uid":    m.FirebaseUid,
	}
}

func (m User) ToMapEmailOnly() map[string]interface{} {
	return map[string]interface{}{
		"email": m.Email,
	}
}
