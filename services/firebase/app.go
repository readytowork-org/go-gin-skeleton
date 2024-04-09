package firebase

import (
	"boilerplate-api/internal/config"
	"context"
	"firebase.google.com/go"
)

// NewFirebaseApp creates new firebase app instance
func NewFirebaseApp(logger config.Logger, opt config.GCPClientOption) *firebase.App {
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Fatalf("Firebase NewApp: %v", err)
	}
	logger.Info("✅ Firebase app initialized.")
	return app
}
