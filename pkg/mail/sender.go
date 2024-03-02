package mail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(subject, content string, to, cc, bcc []string) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

type GmailSenderOptions struct {
	Name              string
	FromEmailAddress  string
	FromEmailPassword string
}

func NewGmailSender(options *GmailSenderOptions) EmailSender {
	return &GmailSender{
		name:              options.Name,
		fromEmailAddress:  options.FromEmailAddress,
		fromEmailPassword: options.FromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(subject, content string, to, cc, bcc []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
