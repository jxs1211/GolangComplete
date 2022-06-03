package mail

import "net/smtp"

const (
	DISCLAIMER = "disclaimer content here"
)

type MailSender interface {
	Send(subject, from string, to, bcc, cc []string, text, server string, a smtp.Auth) error
}

func attachDisclaimer(s string) string {
	return s + "\n\n" + DISCLAIMER
}

func SendMailWithDisclaimer(sender MailSender, subject, from string, to, bcc, cc []string, text, server string, a smtp.Auth) error {
	return sender.Send(subject, from, to, bcc, cc, attachDisclaimer(text), server, a)
}

// "-->" means dependent on
// SendMailWithDisclaimer --> MailSender --> FakeMailSender
//  --> MailSenderAdapter --> "github.com/jordan-wright/email"
