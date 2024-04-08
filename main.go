package main

import (
	"boilerplate-api/bootstrap"
	"go.uber.org/fx"
)

//	@title						Golf Simulation API
//	@version					1.0
//	@description				An API in Go using Gin framework
//	@host						localhost:8000
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used
func main() {
	fx.New(bootstrap.Module).Run()
}
