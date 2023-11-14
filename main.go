package main

import (
	"boilerplate-api/config"

	"go.uber.org/fx"
)

// @title 		Boilerplate API
// @version		1.0
// @description An API in Go using Gin framework
// @host 		localhost:8000
func main() {
	fx.New(config.Module).Run()
}
