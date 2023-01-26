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

