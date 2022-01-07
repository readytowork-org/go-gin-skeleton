package models

type EmailParams struct {
	To              string
	SubjectData     string
	SubjectTemplate string
	BodyData        interface{}
	BodyTemplate    string
	Lang            string
}

type EmailSubject struct {
	Data string
}
type EmailBody struct {
	UserName     string
	UserEmail    string
	InquiryType  string
	DashboardUrl string
	EmailAddress string
	LPUrl string
}
