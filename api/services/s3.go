package services

import (
	"boilerplate-api/infrastructure"
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3BucketService -> handles the file upload functions
type S3BucketService struct {
	logger infrastructure.Logger
	client *s3.Client
	env    infrastructure.Env
}

// NewS3BucketService -> initilization for the AWS S3 BucketService struct
func NewS3BucketService(
	logger infrastructure.Logger,
	client *s3.Client,
	env infrastructure.Env,
) S3BucketService {
	return S3BucketService{
		logger: logger,
		client: client,
		env:    env,
	}
}

// UploadFile -> uploads the file to the aws s3 bucket
func (s S3BucketService) UploadtoS3(
	file multipart.File,
	fileHeader *multipart.FileHeader,
	fileName string,
) (string, error) {

	uploader := manager.NewUploader(s.client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.env.AWS_S3_BUCKET),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("content-type")),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		s.logger.Zap.Fatalf("aws s3 cloud bucket upload error: %v", err.Error())
		return "", err
	}
	return result.Location, nil
}
