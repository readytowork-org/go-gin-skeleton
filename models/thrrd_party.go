package models

type MerchanntRegisterInput struct {
	Name        string `form:"name" `
	Gender      string `form:"gender"`
	BirthDate   string `form:"birth_date"`
	Pin         string `form:"pin"`
	Email       string `form:"email"`
	Password    string `form:"password"`
	DeviceToken string `form:"device_token"`
}
