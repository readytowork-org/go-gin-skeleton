package seeds

import (
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
)

// AdminSeed  Admin seeding
type AdminSeed struct {
	logger          infrastructure.Logger
	firebaseSerivce services.FirebaseService
	env             infrastructure.Env
}

// NewAdminSeed creates admin seed
func NewAdminSeed(
	logger infrastructure.Logger,
	firebaseSerivce services.FirebaseService,
	env infrastructure.Env,
) AdminSeed {
	return AdminSeed{
		logger:          logger,
		firebaseSerivce: firebaseSerivce,
		env:             env,
	}
}

// Run the seed data
func (c AdminSeed) Run() {

	email := c.env.AdminEmail
	password := c.env.AdminPass
	adminUser := models.FirebaseAuthUser{}
	adminUser.Email = email
	adminUser.Password = password
	c.logger.Zap.Info("🌱 seeding  admin data...")

	// err := c.firebaseSerivce.GetUserByEmail(email)

	// if err != "" {
	// 	err := c.firebaseSerivce.CreateAdminUser(adminUser)
	// 	if err != nil {
	// 		c.logger.Zap.Error("Firebase Admin user can't be created: ", err.Error())
	// 		return
	// 	}

	// 	c.logger.Zap.Info("Firebase Admin User Created, email: ", email, " password: ", password)
	// }

	// c.logger.Zap.Info("Admin already exist")

}
