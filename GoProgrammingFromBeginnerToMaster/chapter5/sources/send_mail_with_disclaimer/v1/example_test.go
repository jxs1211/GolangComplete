package mail

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

type RealEmailSender struct {
	e *email.Email
}

func (s *RealEmailSender) Send(from, subject, text, mailserver string, to, bcc, cc []string, a smtp.Auth) error {
	s.e.From = from
	s.e.Subject = subject
	s.e.Text = []byte(text)
	s.e.To = to
	s.e.Bcc = bcc
	s.e.Cc = cc
	return s.e.Send(mailserver, a)
}

func TestSendMail(t *testing.T) {
	from := "jxs1211@gmail.com"
	to := []string{from}
	subject := "go email test subject"
	text := "go email test body"
	mailserver := "smtp.gmail.com:578"
	a := smtp.PlainAuth("", from, "Jxs93503", "smtp.gmail.com")

	s := &RealEmailSender{e: &email.Email{}}

	err := SendEmailWithDisclaimer2(s, subject, from, text, mailserver, to, []string{}, []string{}, a)
	if err != nil {
		t.Errorf("want: nil, actual: %s\n", err)
	}
	want := fmt.Sprintf("%s \n\n %s", text, DISCLAIMER)
	got := string(s.e.Text)
	if got != want {
		t.Fatalf("want: %s, got: %s\n", want, got)
	}
}
