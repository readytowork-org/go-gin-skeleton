package services

import (
	"context"
)

// AuthErrorResponse structure
type AuthErrorResponse struct {
	Message   string `json:"message"`
	ErrorType int    `json:"error_type"`
}

// FirebaseToken Replace this with firebase.google.com/go/auth auth.Token
type FirebaseToken struct {
	AuthTime int64                  `json:"auth_time"`
	Issuer   string                 `json:"iss"`
	Audience string                 `json:"aud"`
	Expires  int64                  `json:"exp"`
	IssuedAt int64                  `json:"iat"`
	Subject  string                 `json:"sub,omitempty"`
	UID      string                 `json:"uid,omitempty"`
	Claims   map[string]interface{} `json:"-"`
}

type IFirebaseAdminSeed interface {
	GetUserByEmail(context context.Context, email string) (interface{}, error)
	CreateUser(displayName, email, password, role string) (string, *AuthErrorResponse)
}

type IFirebaseMiddlewareService interface {
	VerifyToken(idToken string) (*FirebaseToken, *AuthErrorResponse)
}

// TODO :: Uncomment this and add it to module.go if you are using firebase middleware
// FirebaseAuthContract is a contract implementation that combines
// IFirebaseMiddlewareService and go_firebase_service.IAuthService.
// It is used to provide a middleware service for handling Firebase
// authentication functionalities like token verification.

//type FirebaseAuthContract struct {
//	IFirebaseMiddlewareService
//	go_firebase_service.IAuthService
//}
//
//func NewFirebaseAuthContract(authService go_firebase_service.IAuthService) IFirebaseMiddlewareService {
//	return FirebaseAuthContract{
//		IAuthService: authService,
//	}
//}
//
//func (c FirebaseAuthContract) VerifyToken(idToken string) (*FirebaseToken, *AuthErrorResponse) {
//	token, err := c.IAuthService.VerifyToken(idToken)
//	if err != nil {
//		return nil, &AuthErrorResponse{ErrorType: err.ErrorType, Message: err.Message}
//	}
//
//	return &FirebaseToken{
//		AuthTime: token.AuthTime,
//		Issuer:   token.Issuer,
//		Audience: token.Audience,
//		Expires:  token.Expires,
//		IssuedAt: token.IssuedAt,
//		Subject:  token.Subject,
//		UID:      token.UID,
//		Claims:   token.Claims,
//	}, nil
//}
