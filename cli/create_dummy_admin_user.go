package cli

import (
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"

	"github.com/manifoldco/promptui"
)

// CreateDummyAdminUser command
type CreateDummyAdminUser struct {
	logger          infrastructure.Logger
	firebaseSerivce services.FirebaseService
}

// NewCreateDummyAdminUser creates instance of admin user
func NewCreateDummyAdminUser(
	logger infrastructure.Logger,
	firebaseService services.FirebaseService,
) CreateDummyAdminUser {
	return CreateDummyAdminUser{
		logger:          logger,
		firebaseSerivce: firebaseService,
	}
}

// Run runs command
func (c CreateDummyAdminUser) Run() {
	c.logger.Zap.Info("+ Creating dummy admin user...")
	emailPrompt := promptui.Prompt{
		Label: "Dummy Admin Email",
	}

	email, _ := emailPrompt.Run()

	passwordPrompt := promptui.Prompt{
		Label: "Password",
	}

	password, _ := passwordPrompt.Run()

	_, err := c.firebaseSerivce.CreateUser(email, password)

	if err != nil {
		c.logger.Zap.Error("firebase dummy admin user can't be created: ", err.Error())
		return
	}

	c.logger.Zap.Info("Firebase dummy admin user created, email: ", email, " password: ", password)

}

// Name return name of command
func (c CreateDummyAdminUser) Name() string {
	return "CREATE_DUMMY_ADMIN_USER"
}
