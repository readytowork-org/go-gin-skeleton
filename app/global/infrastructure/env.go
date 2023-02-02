package infrastructure

import (
	"log"

	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	LogOutput   string `mapstructure:"LogOutput"`
	DBUsername  string `mapstructure:"DB_USERNAME"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBName      string `mapstructure:"DB_NAME"`
	SentryDSN   string `mapstructure:"SENTRY_DSN"`

	StorageBucketName string `mapstructure:"STORAGE_BUCKET_NAME"`

	AdminEmail string `mapstructure:"ADMIN_EMAIL"`
	AdminPass  string `mapstructure:"ADMIN_PASS"`

	MailClientID     string `mapstructure:"MAIL_CLIENT_ID"`
	MailClientSecret string `mapstructure:"MAIL_CLIENT_SECRET"`
	MailAccesstoken  string `mapstructure:"MailAccessToken"`
	MailRefreshToken string `mapstructure:"MailRefreshToken"`

	AWS_S3_REGION  string `mapstructure:"AWS_S3_REGION"`
	AWS_S3_BUCKET  string `mapstructure:"AWS_S3_BUCKET"`
	AWS_ACCESS_KEY string `mapstructure:"AWS_ACCESS_KEY"`
	AWS_SECRET_KEY string `mapstructure:"AWS_SECRET_KEY"`

	TwilioBaseURL   string `mapstructure:"TWILIO_BASE_URL"`
	TwilioSID       string `mapstructure:"TWILIO_SID"`
	TwilioAuthToken string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioSMSFrom   string `mapstructure:"TWILIO_SMS_FROM"`
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
