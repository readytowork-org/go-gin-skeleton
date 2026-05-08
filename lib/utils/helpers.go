package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPlainPassword(HashedPassword, PlainPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(HashedPassword), []byte(PlainPassword)); err != nil {
		return false
	}
	return true
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// StringToPtr converts a string to a *string
func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
