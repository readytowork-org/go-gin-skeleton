package main

import (
	"boilerplate-api/bootstrap"

	"go.uber.org/fx"
)

// @title 									Boilerplate API
//
//	@version								1.0
//	@description							An API in Go using Gin framework
//	@host									localhost:8000
//	@securityDefinitions.apikey				Bearer
//	@in										header
//	@name									Authorization
//	@description							Description for what is this security definition being used
//	@securitydefinitions.oauth2.implicit	firebase
//	@authorizationUrl						https://accounts.google.com/o/oauth2/v2/auth
//	@x-google-issuer						{"key": "https://securetoken.google.com/<key-here>"}
//	@x-google-jwks_uri						{"key": "https://www.googleapis.com/service_accounts/v1/metadata/x509/securetoken@system.gserviceaccount.com"}
//	@x-google-audiences						{"key": "key-here"}
//	@x-google-jwt-locations					{"header": "Authorization", "value_prefix": "Bearer "}
func main() {
	fx.New(bootstrap.Module).Run()
}
