package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	HOST        string `mapstructure:"HOST"`
	TimeZone    string `mapstructure:"TZ"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"ENVIRONMENT" validate:"required,oneof=local development production test"`
	LogOutput   string `mapstructure:"LOG_OUTPUT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`

	DBType     string `mapstructure:"DB_TYPE" validate:"required"`
	DBUsername string `mapstructure:"DB_USERNAME" validate:"required"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST" validate:"required"`
	DBPort     string `mapstructure:"DB_PORT" validate:"required"`
	DBName     string `mapstructure:"DB_NAME" validate:"required"`

	DBMaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`

	SentryDSN string `mapstructure:"SENTRY_DSN"`

	CorsAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`

	StorageBucketName string `mapstructure:"STORAGE_BUCKET_NAME"`
	ServiceAccountKey string `mapstructure:"SERVICE_ACCOUNT_KEY"`

	AdminEmail string `mapstructure:"ADMIN_EMAIL"`
	AdminPass  string `mapstructure:"ADMIN_PASS"`
	AdminName  string `mapstructure:"ADMIN_NAME"`

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
	JwtAccessSecret          string `mapstructure:"JWT_ACCESS_SECRET" validate:"required,min=16"`
	JwtRefreshSecret         string `mapstructure:"JWT_REFRESH_SECRET" validate:"required,min=16"`
	JwtAccessTokenExpiresAt  int    `mapstructure:"JWT_ACCESS_TOKEN_EXPIRES_AT" validate:"required,min=1"`
	JwtRefreshTokenExpiresAt int    `mapstructure:"JWT_REFRESH_TOKEN_EXPIRES_AT" validate:"required,min=1"`

	IdempotencyStore string        `mapstructure:"IDEMPOTENCY_STORE" validate:"omitempty,oneof=none mysql redis"`
	IdempotencyTTL   time.Duration `mapstructure:"IDEMPOTENCY_TTL"`
	RedisAddr        string        `mapstructure:"REDIS_ADDR"`
	RedisPassword    string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB          int           `mapstructure:"REDIS_DB"`

	RateLimitPeriod   time.Duration `mapstructure:"RATE_LIMIT_PERIOD"`
	RateLimitRequests int64         `mapstructure:"RATE_LIMIT_REQUESTS"`

	ProjectName       string `mapstructure:"PROJECT_NAME"`
	BillingAccountId  string `mapstructure:"BILLING_ACCOUNT_ID"`
	BudgetDisplayName string `mapstructure:"BUDGET_DISPLAY_NAME"`
	BudgetAmount      int64  `mapstructure:"BUDGET_AMOUNT"`
	SetBudget         int    `mapstructure:"SET_BUDGET"`

	StripeSecretKey   string `mapstructure:"STRIPE_SECRET_KEY"`
	StripeProductID   string `mapstructure:"STRIPE_PRODUCT_ID"`
	StripeWebhookKey  string `mapstructure:"STRIPE_WEBHOOK_KEY"`
	StripeRedirectUrl string `mapstructure:"STRIPE_REDIRECT_URL"`
}

type EnvPath string

func (p EnvPath) ToString() string {
	return string(p)
}

// NewEnv creates a new environment
func NewEnv(envPath EnvPath) Env {
	env := Env{}
	_ = godotenv.Load(envPath.ToString())
	viper.SetConfigFile(envPath.ToString())

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("☠️ Env config file not found: %+v", err)
		} else {
			log.Fatalf("☠️ Env config file error: %+v", err)
		}
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatalf("☠️ environment can't be loaded: %+v", err)
	}

	if env.TimeZone == "" {
		env.TimeZone = "UTC"
	}

	if err := validateEnv(&env); err != nil {
		log.Fatalf("☠️ environment validation failed:\n%s", err.Error())
	}

	return env
}

// validateEnv runs struct validation and aggregates all issues into a single
// error message so missing/invalid vars are reported together at startup.
func validateEnv(env *Env) error {
	v := validator.New()
	err := v.Struct(env)
	if err == nil {
		return nil
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var b strings.Builder
	for _, fe := range verrs {
		fmt.Fprintf(&b, "  - %s failed %q (got %q)\n", fe.Field(), fe.Tag(), fmt.Sprint(fe.Value()))
	}
	return fmt.Errorf("%s", b.String())
}
