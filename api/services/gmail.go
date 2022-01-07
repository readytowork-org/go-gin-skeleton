package services

import (
	"encoding/base64"
	"errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"

	"google.golang.org/api/gmail/v1"
)

type GmailService struct {
	gmailService *gmail.Service
	logger       infrastructure.Logger
}

func NewGmailService(gmailService *gmail.Service, logger infrastructure.Logger) GmailService {
	return GmailService{
		gmailService: gmailService,
		logger:       logger,
	}
}

func (g GmailService) SendEmail(params models.EmailParams) (bool, error) {
	to := params.To
	emailBody, err := utils.ParseTemplate(params.BodyTemplate, params.BodyData)
	if err != nil {
		return false, errors.New("unable to parse email body template")
	}
	var msgString string
	emailTo := "To: " + to + "\r\n"
	msgString = emailTo
	subject := "Subject: " + params.SubjectData + "\n"
	msgString = msgString + subject
	msgString = msgString + "\n" + emailBody
	var msg []byte

	if params.Lang != "en" {
		msgStringJP, _ := utils.ToISO2022JP(msgString)
		msg = []byte(msgStringJP)
	} else {
		msg = []byte(msgString)
	}
	message := gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(msg)),
	}
	_, err = g.gmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}
