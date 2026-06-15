package services

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"monthly-journal/internal/config"
)

type EmailService struct {
	SMTPHost   string
	SMTPPort   int
	SMTPUser   string
	SMTPPass   string
	EmailFrom  string
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		SMTPHost:  cfg.SMTPHost,
		SMTPPort:  cfg.SMTPPort,
		SMTPUser:  cfg.SMTPUser,
		SMTPPass:  cfg.SMTPPass,
		EmailFrom: cfg.EmailFrom,
	}
}

func (es *EmailService) SendEmail(recipients []string, subject string, body string) error {
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	auth := smtp.PlainAuth("", es.SMTPUser, es.SMTPPass, es.SMTPHost)
	addr := fmt.Sprintf("%s:%d", es.SMTPHost, es.SMTPPort)

	headers := fmt.Sprintf("From: %s\r\n", es.EmailFrom)
	headers += fmt.Sprintf("To: %s\r\n", strings.Join(recipients, ","))
	headers += fmt.Sprintf("Subject: %s\r\n", subject)
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=UTF-8\r\n\r\n"

	msg := headers + body

	log.Printf("Sending email to %v with subject: %s", recipients, subject)

	if err := smtp.SendMail(addr, auth, es.EmailFrom, recipients, []byte(msg)); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %v", recipients)
	return nil
}
