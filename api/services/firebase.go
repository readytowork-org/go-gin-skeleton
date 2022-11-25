package services

import (
	"boilerplate-api/constants"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

// FirebaseService structure
type FirebaseService struct {
	client *auth.Client
	logger infrastructure.Logger
	env    infrastructure.Env
}

// NewFirebaseService creates new firebase service
func NewFirebaseService(client *auth.Client, logger infrastructure.Logger,
	env infrastructure.Env) FirebaseService {
	return FirebaseService{
		client: client,
		logger: logger,
		env:    env,
	}
}

// CreateUser creates a new user with email and password
func (fb *FirebaseService) CreateUser(newUser models.FirebaseAuthUser) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(newUser.Email).
		Password(newUser.Password).
		DisplayName(newUser.DisplayName).
		PhoneNumber(newUser.PhoneNumber).
		Disabled(false)
	u, err := fb.client.CreateUser(context.Background(), params)
	if err != nil {
		return "", err
	}
	claims := map[string]interface{}{
		"role":   "user",
		"fb_uid": u.UID,
		"id":     newUser.UserId,
	}
	err = fb.client.SetCustomUserClaims(context.Background(), u.UID, claims)
	if err != nil {
		return "Internal Server Error", err
	}
	return u.UID, err
}

// CreateAdminUser creates admin user
func (fb *FirebaseService) CreateAdminUser(newUser models.FirebaseAuthUser) error {
	uid, err := fb.CreateUser(newUser)
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
func (fb *FirebaseService) GetUserByEmail(email string) string {
	user, err := fb.client.GetUserByEmail(context.Background(), email)
	if err != nil {
		return ""
	}
	if user != nil {
		return user.UID
	}
	return ""
}

// UpdateUser -> update firebase user.
func (fb *FirebaseService) UpdateUser(UID string, user models.UserToUpdate) (*auth.UserRecord, error) {
	userToUpdateParam := &auth.UserToUpdate{}
	userToUpdateParam.Email(user.Email)
	userToUpdateParam.DisplayName(user.FullName)
	userToUpdateParam.PhoneNumber(user.Phone)
	return fb.client.UpdateUser(context.Background(), UID, userToUpdateParam)
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

// Check User Email and Password
func (fb *FirebaseService) VerifyUserCredentials(email string, password string) error {
	post_body, err := json.Marshal(map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	})
	if err != nil {
		return err
	}
	req_body := bytes.NewBuffer(post_body)
	response, err := http.Post(
		constants.FirebaseUseLoginUrl+fb.env.FirebaseApiKey,
		"application/json",
		req_body,
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.BadRequest.New("Incorrect user credentials")
	}

	return nil
}

// Delete all auth user in firebase
func (fb *FirebaseService) DeleteAllUsers(uid string) error {
	if fb.env.Environment != "local" {
		return errors.Unauthorized.New("Not Authorised")
	}
	iter := fb.client.Users(context.Background(), "")
	var uid_list []string
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error listing users: %s\n", err)
		}
		if user.UID != uid {
			uid_list = append(uid_list, user.UID)
		}
	}

	deleteUsersResult, err := fb.client.DeleteUsers(context.Background(), uid_list)
	if err != nil {
		log.Fatalf("error deleting users: %v\n", err)
	}

	log.Printf("Successfully deleted %d users", deleteUsersResult.SuccessCount)
	log.Printf("Failed to delete %d users", deleteUsersResult.FailureCount)
	return nil
}
