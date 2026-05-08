package services

import (
	"context"
	"io"
)

type GcpStorageBucketService interface {
	UploadFile(ctx context.Context, file io.Reader, fileName string) (string, error)
}
