package mail_test

import (
	"net/smtp"
	"testing"

	mail "github.com/GoProgrammingFromBeginnerToMaster/chapter5/sources/send_mail_with_disclaimer2/v2"
)

type FakeMailSender struct {
	subject string
	from    string
	to      []string
	bcc     []string
	cc      []string
	text    string
	server  string
	a       smtp.Auth
}

func (s *FakeMailSender) Send(subject, from string, to, bcc, cc []string, text, server string, a smtp.Auth) error {
	s.subject = subject
	s.from = from
	s.to = to
	s.bcc = bcc
	s.cc = cc
	s.text = text
	s.server = server
	s.a = a
	return nil
}

func TestSendMailWithDisclaimer(t *testing.T) {
	type args struct {
		sender  mail.MailSender
		subject string
		from    string
		to      []string
		bcc     []string
		cc      []string
		text    string
		server  string
		a       smtp.Auth
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"base case",
			args{&FakeMailSender{}, "subject", "test@gmail.com", []string{"test@gmail.com"}, []string{"test@gmail.com"}, []string{"test@gmail.com"}, "text11", "smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mail.SendMailWithDisclaimer(tt.args.sender, tt.args.subject, tt.args.from, tt.args.to, tt.args.bcc, tt.args.cc, tt.args.text, tt.args.server, tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("SendMailWithDisclaimer() error = %v, wantErr %v", err, tt.wantErr)
			}
			want := tt.args.text + "\n\n" + mail.DISCLAIMER
			got := (tt.args.sender.(*FakeMailSender)).text
			if got != want {
				t.Errorf("want: %s, got: %s", want, got)
			}
		})
	}
}

func TestSendMailWithDisclaimerWithoutTableDrivenTest(t *testing.T) {
	sender := &FakeMailSender{}
	text := "hello world"
	err := mail.SendMailWithDisclaimer(
		sender,
		"subject", "test@gmail.com", []string{"test@gmail.com"}, []string{"test@gmail.com"}, []string{"test@gmail.com"}, text, "smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com"),
	)
	if err != nil {
		t.Errorf("SendMailWithDisclaimer() err: %v\n", err)
	}
	want := text + "\n\n" + mail.DISCLAIMER
	got := sender.text
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
