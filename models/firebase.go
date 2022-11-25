package models

type FirebaseAuthUser struct {
	Email       string
	Password    string
	DisplayName string
	Role        string
	Enabled     int
	UserId      string
	PhoneNumber string
}
