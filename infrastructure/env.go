package infrastructure

import (
	"os"
)

// Env has environment stored
type Env struct {
	ServerPort  string
	Environment string
	LogOutput   string
	DBUsername  string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string
	SentryDSN   string

	StorageBucketName string

	AdminEmail string
	AdminPass  string

	MailClientID     string
	MailClientSecret string
	MailAccesstoken  string
	MailRefreshToken string

	AWS_S3_REGION  string
	AWS_S3_BUCKET  string
	AWS_ACCESS_KEY string
	AWS_SECRET_KEY string

	TwilioBaseURL   string
	TwilioSID       string
	TwilioAuthToken string
	TwilioSMSFrom   string
}

// NewEnv creates a new environment
func NewEnv() Env {
	env := Env{}
	env.LoadEnv()
	return env
}

// LoadEnv loads environment
func (env *Env) LoadEnv() {
	env.ServerPort = os.Getenv("ServerPort")
	env.Environment = os.Getenv("Environment")
	env.LogOutput = os.Getenv("LogOutput")

	env.DBUsername = os.Getenv("DBUsername")
	env.DBPassword = os.Getenv("DBPassword")
	env.DBHost = os.Getenv("DBHost")
	env.DBPort = os.Getenv("DBPort")
	env.DBName = os.Getenv("DBName")

	env.SentryDSN = os.Getenv("SentryDSN")
	env.StorageBucketName = os.Getenv("StorageBucketName")

	env.AdminEmail = os.Getenv("AdminEmail")
	env.AdminPass = os.Getenv("AdminPass")

	env.MailClientID = os.Getenv("MailClientID")
	env.MailClientSecret = os.Getenv("MailClientSecret")
	env.MailAccesstoken = os.Getenv("MailAccesstoken")
	env.MailRefreshToken = os.Getenv("MailRefreshToken")

	env.AWS_S3_REGION = os.Getenv("AWS_S3_REGION")
	env.AWS_S3_BUCKET = os.Getenv("AWS_S3_BUCKET")
	env.AWS_ACCESS_KEY = os.Getenv("AWS_ACCESS_KEY")
	env.AWS_SECRET_KEY = os.Getenv("AWS_SECRET_KEY")

	env.TwilioBaseURL = os.Getenv("TwilioBaseURL")
	env.TwilioAuthToken = os.Getenv("TwilioAuthToken")
	env.TwilioSID = os.Getenv("TwilioSID")
	env.TwilioSMSFrom = os.Getenv("TwilioSMSFrom")
}
