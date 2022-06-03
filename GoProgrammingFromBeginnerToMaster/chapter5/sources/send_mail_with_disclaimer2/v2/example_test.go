package mail_test

import (
	"net/smtp"
	"testing"

	mail "github.com/GoProgrammingFromBeginnerToMaster/chapter5/sources/send_mail_with_disclaimer2/v2"
	"github.com/jordan-wright/email"
)

type MailSenderAdapter struct {
	e *email.Email
}

func (s *MailSenderAdapter) Send(subject, from string, to, bcc, cc []string, text, server string, a smtp.Auth) error {
	s.e.Subject = subject
	s.e.From = from
	s.e.To = to
	s.e.Bcc = bcc
	s.e.Cc = cc
	s.e.Text = []byte(text)
	return s.e.Send(server, a)
}

func TestSendMailWithDisclaimerWithMailSenderAdapter(t *testing.T) {
	type args struct {
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
		s       mail.MailSender
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"base case",
			&MailSenderAdapter{e: email.NewEmail()},
			args{"subject", "test@gmail.com", []string{"test@gmail.com"}, []string{"test@gmail.com"}, []string{"test@gmail.com"}, "content", "smtp.163.com:25", smtp.PlainAuth("", "test@gmail.com", "123", "smtp.163.com")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mail.SendMailWithDisclaimer(tt.s, tt.args.subject, tt.args.from, tt.args.to, tt.args.bcc, tt.args.cc, tt.args.text, tt.args.server, tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("SendMailWithDisclaimer() error = %v, wantErr %v", err, tt.wantErr)
			}
			want := tt.args.text + "\n\n" + mail.DISCLAIMER
			got := string((tt.s.(*MailSenderAdapter)).e.Text)
			if want != got {
				t.Errorf("want: %s, got: %s", want, got)
			}
		})
	}
}
