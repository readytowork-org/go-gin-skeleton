package infrastructure

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	HOST        string `mapstructure:"HOST"`
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
	MailAccesstoken  string `mapstructure:"MAIL_ACCESS_TOKEN"`
	MailRefreshToken string `mapstructure:"MAIL_REFRESH_TOKEN"`

	AwsS3Region  string `mapstructure:"AWS_S3_REGION"`
	AwsS3Bucket  string `mapstructure:"AWS_S3_BUCKET"`
	AwsAccessKey string `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey string `mapstructure:"AWS_SECRET_KEY"`

	TwilioBaseURL            string `mapstructure:"TWILIO_BASE_URL"`
	TwilioSID                string `mapstructure:"TWILIO_SID"`
	TwilioAuthToken          string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioSMSFrom            string `mapstructure:"TWILIO_SMS_FROM"`
	JwtAccessSecret          string `mapstructure:"JWT_ACCESS_SECRET"`
	JwtRefreshSecret         string `mapstructure:"JWT_REFRESH_SECRET"`
	JwtAccessTokenExpiresAt  int    `mapstructure:"JWT_ACCESS_TOKEN_EXPIRES_AT"`
	JwtRefreshTokenExpiresAt int    `mapstructure:"JWT_REFRESH_TOKEN_EXPIRES_AT"`

	RateLimitPeriod   time.Duration `mapstructure:"RATE_LIMIT_PERIOD"`
	RateLimitRequests int64         `mapstructure:"RATE_LIMIT_REQUESTS"`

	ProjectName       string `mapstructure:"PROJECT_NAME"`
	BillingAccountId  string `mapstructure:"BILLING_ACCOUNT_ID"`
	BudgetDisplayName string `mapstructure:"BUDGET_DISPLAY_NAME"`
	BudgetAmount      int64  `mapstructure:"BUDGET_AMOUNT"`
	SetBudget         int    `mapstructure:"SET_BUDGET"`

	OAuthClientId     string `mapstructure:"OAUTH_CLIENT_ID"`
	OAuthClientSecret string `mapstructure:"OAUTH_CLIENT_SECRET"`
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
