package service

import (
	"github.com/iscritic/archive-api/internal/config"
	"gopkg.in/gomail.v2"
	"io"
)

type EmailService struct {
	Dialer *gomail.Dialer
	From   string
}

func NewEmailService(cfg *config.SMTPConfig) *EmailService {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)
	return &EmailService{
		Dialer: dialer,
		From:   cfg.User,
	}
}

func (s *EmailService) SendEmail(to string, subject string, body string, attachmentFilename string, attachmentReader io.Reader, attachmentType string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if attachmentReader != nil && attachmentFilename != "" && attachmentType != "" {
		m.Attach(attachmentFilename, gomail.SetHeader(map[string][]string{
			"Content-Type": {attachmentType},
		}), gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, attachmentReader)
			return err
		}))
	}

	return s.Dialer.DialAndSend(m)
}
