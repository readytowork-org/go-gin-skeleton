package main

import (
	"boilerplate-api/bootstrap"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	godotenv.Load()
	fx.New(bootstrap.Module).Run()
}
