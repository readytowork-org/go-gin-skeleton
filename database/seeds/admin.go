package seeds

import (
	"context"

	"boilerplate-api/external_services/firebase"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
)

// AdminSeed  Admin seeding
type AdminSeed struct {
	logger          config.Logger
	firebaseService firebase.AuthService
	env             config.Env
}

// NewAdminSeed creates admin seed
func NewAdminSeed(
	logger config.Logger,
	firebaseSerivce firebase.AuthService,
	env config.Env,
) AdminSeed {
	return AdminSeed{
		logger:          logger,
		firebaseService: firebaseSerivce,
		env:             env,
	}
}

// Run the seed data
func (c AdminSeed) Run() {

	email := c.env.AdminEmail
	password := c.env.AdminPass
	name := c.env.AdminName

	c.logger.Info("🌱 seeding  admin data...")

	_, err := c.firebaseService.GetUserByEmail(context.Background(), email)

	if err != nil {
		firebaseAuthUser := firebase.AuthUser{
			Password:    password,
			Email:       email,
			Role:        string(constants.Roles.SuperAdmin),
			DisplayName: &name,
		}

		_, err = c.firebaseService.CreateUser(firebaseAuthUser)
		if err != nil {
			c.logger.Error("Firebase Admin user can't be created: ", err.Error())
			return
		}

		c.logger.Info("Firebase Admin UserName Created, email: ", email, " password: ", password)
	}

	c.logger.Info("Admin already exist")
}
