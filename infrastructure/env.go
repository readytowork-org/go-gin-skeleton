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

	AdminEmail        string
	AdminPassword     string
	RecruitAdminEmail string

	MailClientID     string
	MailClientSecret string
	MailRefreshToken string
	MailAccesstoken  string
	FirebaseApiKey   string

	UserURI        string
	CompanyURI     string
	AdminURI       string
	FirebaseAppURI string
	ApiURL         string
	GcpProjectId   string

	AWS_ACCESS_KEY string
	AWS_SECRET_KEY string
	AWS_S3_REGION  string
	AWS_S3_BUCKET  string

	MerchantRegisterUrl string
}

// NewEnv creates a new environment
func NewEnv() Env {
	env := Env{}
	env.LoadEnv()
	return env
}

// LoadEnv loads environment
func (env *Env) LoadEnv() {
	env.ServerPort = os.Getenv("SERVER_PORT")
	env.Environment = os.Getenv("ENVIRONMENT")
	env.LogOutput = os.Getenv("LOG_OUTPUT")

	env.DBUsername = os.Getenv("DB_USER")
	env.DBPassword = os.Getenv("DB_PASS")
	env.DBHost = os.Getenv("DB_HOST")
	env.DBPort = os.Getenv("DB_PORT")
	env.DBName = os.Getenv("DB_NAME")

	env.SentryDSN = os.Getenv("SENTRY_DSN")
	env.StorageBucketName = os.Getenv("STORAGE_BUCKET_NAME")

	env.AdminEmail = os.Getenv("ADMIN_EMAIL")
	env.AdminPassword = os.Getenv("ADMIN_PASSWORD")
	env.RecruitAdminEmail = os.Getenv("RECRUIT_ADMIN_EMAIL")

	env.MailClientID = os.Getenv("MAIL_CLIENT_ID")
	env.MailClientSecret = os.Getenv("MAIL_CLIENT_SECRET")
	env.MailRefreshToken = os.Getenv("MAIL_REFRESH_TOKEN")
	env.MailAccesstoken = os.Getenv("MAIL_ACCESS_TOKEN")
	env.FirebaseApiKey = os.Getenv("FIREBASE_API_KEY")
	env.UserURI = os.Getenv("USER_URI")
	env.CompanyURI = os.Getenv("COMPANY_URI")
	env.AdminURI = os.Getenv("ADMIN_URI")
	env.FirebaseAppURI = os.Getenv("FIREBASE_APP_URI")
	env.ApiURL = os.Getenv("API_URL")
	env.GcpProjectId = os.Getenv("GCP_PROJECT_ID")

	env.AWS_ACCESS_KEY = os.Getenv("AWS_ACCESS_KEY")
	env.AWS_SECRET_KEY = os.Getenv("AWS_SECRET_KEY")
	env.AWS_S3_REGION = os.Getenv("AWS_S3_REGION")
	env.AWS_S3_BUCKET = os.Getenv("AWS_S#_BUCKET")
	env.MerchantRegisterUrl = os.Getenv("MERCHANT_REGISTER_API_URL")

}
