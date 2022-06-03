package mail_test

import (
	"net/smtp"
	"testing"

	mail "github.com/GoProgrammingFromBeginnerToMaster/chapter5/sources/send_mail_with_disclaimer2/v1"
)

func TestSendMailWithDisclaimer(t *testing.T) {
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
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"base case",
			args{"subject", "test@gmail.com", []string{"test@gmail.com"}, []string{"test@gmail.com"}, []string{"test@gmail.com"}, "text11", "smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mail.SendMailWithDisclaimer(tt.args.subject, tt.args.from, tt.args.to, tt.args.bcc, tt.args.cc, tt.args.text, tt.args.server, tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("SendMailWithDisclaimer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
