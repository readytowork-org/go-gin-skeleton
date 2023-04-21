package infrastructure

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// NewBucketStorage -> creates a new storage client
func NewBucketStorage(logger Logger, env Env) *storage.Client {
	bucketName := env.StorageBucketName
	ctx := context.Background()
	if bucketName == "" {
		logger.Zap.Error("Please check your env file for STORAGE_BUCKET_NAME")
	}
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("serviceAccountKey.json"))
	if err != nil {
		logger.Zap.Errorf(err.Error())
	}

	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		logger.Zap.Errorf("Provided bucket %v doesn't exists", bucketName)
	}
	if err != nil {
		logger.Zap.Errorf("Cloud bucket error: %v", err.Error())
	}
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		CORS: []storage.CORS{
			{
				MaxAge:          600,
				Methods:         []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
				Origins:         []string{"*"},
				ResponseHeaders: []string{"Content-Type"},
			}},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		logger.Zap.Errorf("Cloud bucket update error: %v", err.Error())
	}
	return client
}
