#!/bin/bash

# Prompt the user for the target folder
read -p "Enter the target folder: " folder

# Create the folder if it doesn't exist
mkdir -p "$folder"

# Extract the last folder name (e.g., "jobs" from "api/admin/jobs")
last_folder=$(basename "$folder")
echo "Using package name: $last_folder"

module_name=$(grep '^module ' go.mod | awk '{print $2}')

# Create repository.go with its template
cat <<EOF > "$folder/repository.go"
package $last_folder

import "$module_name/lib/config"

type IRepository interface{}

type Repository struct {
	db     config.Database
	logger config.Logger
}

func NewRepository(db config.Database, logger config.Logger) Repository {
	return Repository{
		db:     db,
		logger: logger,
	}
}
EOF

# Create service.go with its template
cat <<EOF > "$folder/service.go"
package $last_folder

type IService interface{}

type Service struct {
	repo IRepository
}

func NewService(repo IRepository) Service {
	return Service{
		repo: repo,
	}
}
EOF

# Create controller.go with its template
cat <<EOF > "$folder/controller.go"
package $last_folder

import (
	"$module_name/lib/config"
	"$module_name/lib/request_validator"
)

type Controller struct {
	logger    config.Logger
	env       config.Env
	validator request_validator.Validator
	service   IService
}

func NewController(
	logger config.Logger,
	env config.Env,
	validator request_validator.Validator,
	service IService,
) Controller {
	return Controller{
		logger:    logger,
		env:       env,
		validator: validator,
		service:   service,
	}
}
EOF

# Create routes.go with its template
cat <<EOF > "$folder/routes.go"
package $last_folder

import (
	"$module_name/lib/config"
	"$module_name/lib/middlewares"
	"$module_name/lib/router"
)

// SetupRoutes user routes
func SetupRoutes(
	logger config.Logger,
	router router.Router,
	controller Controller,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
) {
	logger.Info(" Setting up pic routes")
}
EOF

# Create modules.go with its template
cat <<EOF > "$folder/modules.go"
package $last_folder

import "go.uber.org/fx"

var Module = fx.Module(
	"$last_folder",
	fx.Options(
		fx.Provide(
			fx.Annotate(NewRepository, fx.As(new(IRepository))),
			fx.Annotate(NewService, fx.As(new(IService))),
			NewController,
		),
		fx.Invoke(SetupRoutes),
	),
)

EOF

echo "Files created in folder: $folder"
