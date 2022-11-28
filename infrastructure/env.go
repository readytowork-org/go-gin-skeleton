package infrastructure

import (
	"log"

	"github.com/spf13/viper"
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
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("☠️ Env config file not found: ", err)
		} else {
			log.Fatal("☠️ Env config file error: ", err)
		}
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	log.Printf("%+v \n", env)
	return env
}
