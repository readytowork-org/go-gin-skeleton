package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewFBApp creates new firebase app instance
func NewS3Client(logger Logger, config aws.Config, env Env) *s3.Client {
	client := s3.New(s3.Options{Credentials: config.Credentials, Region: env.AWS_S3_REGION})
	logger.Zap.Info("✅  AWS S3 client created.")
	return client
}
