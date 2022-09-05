package models

type Mailer struct {
	Sender     string
	SenderName string
	To         []string
	Subject    string
	Body       string
}
