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
	env.TwilioBaseURL = os.Getenv("TwilioBaseURL")
	env.TwilioAuthToken = os.Getenv("TwilioAuthToken")
	env.TwilioSID = os.Getenv("TwilioSID")
	env.TwilioSMSFrom = os.Getenv("TwilioSMSFrom")
}
