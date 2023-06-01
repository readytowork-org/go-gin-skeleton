package infrastructure

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// NewAWSConfig creates new config instance from default aws profile in ~/.aws/credentials file
func NewAWSConfig(logger Logger, env Env) aws.Config {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				env.AwsAccessKey,
				env.AwsSecretKey,
				""),
		),
	)
	if err != nil {
		logger.Zap.Panic("Unable to load aws configuration from .env file")
	}
	logger.Zap.Info("✅ AWS config created.")
	return cfg
}
