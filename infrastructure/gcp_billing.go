package infrastructure

import (
	"context"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/option"
)

type GCPBilling struct {
	Client *cloudbilling.APIService
}

// NewGCPBilling creates a new gcp billing api client
func NewGCPBilling(logger Logger, env Env) GCPBilling {
	ctx := context.Background()

	client, err := cloudbilling.NewService(ctx, option.WithCredentialsFile("serviceAccountKey.json"))

	if err != nil {
		logger.Zap.Error("Failed to create cloud billing api client: %v \n", err)
	}
	return GCPBilling{
		Client: client,
	}
}
