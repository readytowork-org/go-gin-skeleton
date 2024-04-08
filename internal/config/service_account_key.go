package config

import (
	"google.golang.org/api/option"
	"path/filepath"
)

type GCPClientOption struct {
	option.ClientOption
}

func NewGCPClientOption(logger Logger) GCPClientOption {
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		logger.Panic("Unable to load serviceAccountKey.json file")
	}

	return GCPClientOption{
		option.WithCredentialsFile(serviceAccountKeyFilePath),
	}
}
