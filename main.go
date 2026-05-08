package main

import (
	"boilerplate-api/bootstrap"
	_ "boilerplate-api/swagger"

	"go.uber.org/fx"
)

//	@title						Demo API
//	@version					1.0
//	@description				An API in Go using Gin framework
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used
func main() {
	fx.New(bootstrap.Module).Run()
}
