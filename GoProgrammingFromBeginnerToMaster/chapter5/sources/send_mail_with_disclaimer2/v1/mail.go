package mail

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

func attachDisclaimer(s string) string {
	return "\n\n" + s
}

func SendMailWithDisclaimer(subject, from string, to, bcc, cc []string, text, server string, a smtp.Auth) error {
	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Bcc = bcc
	e.Cc = cc
	e.Subject = subject
	e.Text = []byte(attachDisclaimer(text))
	return e.Send(server, a)
}
