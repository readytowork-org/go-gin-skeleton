package services

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"boilerplate-api/lib/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"google.golang.org/api/gmail/v1"
)

type EmailBodyTemplateName string

type EmailParams struct {
	To              string
	SubjectData     string
	SubjectTemplate string
	BodyData        interface{}
	BodyTemplate    string
	Lang            string
}

type gLogger interface {
	Fatal(args ...interface{})
}

type GmailConfig struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	HostURL      string
}

type GmailService struct {
	*gmail.Service
	logger gLogger
}

func NewGmailService(gmailConfig GmailConfig, logger gLogger) *GmailService {
	ctx := context.Background()

	oauthConfig := oauth2.Config{
		ClientID:     gmailConfig.ClientID,
		ClientSecret: gmailConfig.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  gmailConfig.HostURL, // e.g: "http://localhost" or deployed API url
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send"},
	}
	token := oauth2.Token{
		AccessToken:  gmailConfig.AccessToken,
		RefreshToken: gmailConfig.RefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}
	var tokenSource = oauthConfig.TokenSource(ctx, &token)
	_service, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		logger.Fatal("failed to receive gmail client", err.Error())
	}

	return &GmailService{
		Service: _service,
		logger:  logger,
	}
}

func (g GmailService) SendEmail(params EmailParams) (bool, error) {
	to := params.To
	emailBody, err := utils.ParseTemplate(params.BodyTemplate, params.BodyData)
	if err != nil {
		return false, errors.New("unable to parse email body template")
	}

	msgString := ""
	msgString += "To: " + to + "\r\n"
	msgString += "Subject: " + params.SubjectData + "\r\n"
	msgString += "MIME-Version: 1.0\r\n"

	if params.Lang == "ja" {
		msgString += "Content-Type: text/html; charset=ISO-2022-JP\r\n"
		msgString += "Content-Transfer-Encoding: 7bit\r\n"
	}

	if params.Lang == "en" {
		msgString += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	}

	msgString += "\r\n"
	msgString += emailBody

	var msg []byte
	if params.Lang == "ja" {
		msg, _ = utils.ToISO2022JP(msgString)
	}

	if params.Lang == "en" {
		msg = []byte(msgString)
	}

	gmailMessage := gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(msg),
	}

	_, err = g.Users.Messages.Send("me", &gmailMessage).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}
