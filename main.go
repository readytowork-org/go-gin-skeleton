package main

import (
	"boilerplate-api/bootstrap"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// @title 		Boilerplate API
// @version		1.0
// @description An API in Go using Gin framework
// @host 		localhost:8000/
func main() {
	_ = godotenv.Load()
	fx.New(bootstrap.Module).Run()
}
