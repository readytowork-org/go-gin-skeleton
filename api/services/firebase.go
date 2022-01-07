package services

import (
	"boilerplate-api/constants"
	"context"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// FirebaseService structure
type FirebaseService struct {
	client *auth.Client
}

// NewFirebaseService creates new firebase service
func NewFirebaseService(client *auth.Client) FirebaseService {
	return FirebaseService{
		client: client,
	}
}

// CreateUser creates a new user with email and password
func (fb *FirebaseService) CreateUser(email, password string) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	u, err := fb.client.CreateUser(context.Background(), params)
	if err != nil {
		return "", err
	}
	return u.UID, err
}

// CreateAdminUser creates admin user
func (fb *FirebaseService) CreateAdminUser(email, password string) error {
	uid, err := fb.CreateUser(email, password)
	if err != nil {
		return err
	}
	claims := gin.H{
		"role": constants.RoleAdmin,
	}
	return fb.SetClaim(uid, claims)
}

// VerifyToken verify passed firebase id token
func (fb *FirebaseService) VerifyToken(idToken string) (*auth.Token, error) {
	token, err := fb.client.VerifyIDToken(context.Background(), idToken)
	return token, err
}

// GetUserByEmail gets the user data corresponding to the specified email.
func (fb *FirebaseService) GetUserByEmail(email string) (*auth.UserRecord, error) {
	user, err := fb.client.GetUserByEmail(context.Background(), email)
	return user, err
}

// UpdateUser -> update firebase user.
func (fb *FirebaseService) UpdateUser(UID string, user *auth.UserToUpdate) (*auth.UserRecord, error) {

	return fb.client.UpdateUser(context.Background(), UID, user)
}

// GetUser gets firebase user from uid
func (fb *FirebaseService) GetUser(uid string) (*auth.UserRecord, error) {
	user, err := fb.client.GetUser(context.Background(), uid)
	return user, err
}

// SetClaim set's claim to firebase user
func (fb *FirebaseService) SetClaim(uid string, claims gin.H) error {

	err := fb.client.SetCustomUserClaims(context.Background(), uid, claims)
	return err

}

// UpdateEmailVerification update firebase user email verify
func (fb *FirebaseService) UpdateEmailVerification(uid string) error {
	params := (&auth.UserToUpdate{}).
		EmailVerified(true)
	_, err := fb.client.UpdateUser(context.Background(), uid, params)
	return err
}

// DisableUser  -> Whether or not the user to  disabled. true for disabled; false for enabled.
func (fb *FirebaseService) DisableUser(uid string, enable bool) error {
	params := (&auth.UserToUpdate{}).
		Disabled(enable)
	_, err := fb.client.UpdateUser(context.Background(), uid, params)
	return err
}

// DeleteUser deletes firebase user
func (fb *FirebaseService) DeleteUser(uid string) error {
	err := fb.client.DeleteUser(context.Background(), uid)
	return err
}
