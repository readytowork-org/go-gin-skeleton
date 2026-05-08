package gcp_services

import (
	"fmt"
	"path/filepath"

	"boilerplate-api/lib/config"

	"google.golang.org/api/option"
)

func NewGCPClientOption(logger config.Logger, env config.Env) option.ClientOption {
	serviceAccountKeyFilePath, err := filepath.Abs(fmt.Sprintf("./%v", env.ServiceAccountKey))
	if err != nil {
		logger.Panic("Unable to load serviceAccountKey.json file")
	}

	return option.WithCredentialsFile(serviceAccountKeyFilePath)
}
