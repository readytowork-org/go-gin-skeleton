package firebase

import (
	"boilerplate-api/internal/config"
	"context"
	"path/filepath"

	"firebase.google.com/go"
	"google.golang.org/api/option"
)

// NewFirebaseApp creates new firebase app instance
func NewFirebaseApp(logger config.Logger) *firebase.App {
	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		logger.Panic("Unable to load serviceAccountKey.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Fatalf("Firebase NewApp: %v", err)
	}
	logger.Info("✅ Firebase app initialized.")
	return app
}
