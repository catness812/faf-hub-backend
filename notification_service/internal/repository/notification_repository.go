package repository

import (
	"fmt"
	"net/smtp"
	"os"
)

type NotificationRepository struct {
	auth smtp.Auth
}

func NewNotificationRepository(auth smtp.Auth) *NotificationRepository {
	return &NotificationRepository{
		auth: auth,
	}
}

func (repo *NotificationRepository) SendMail(to []string, subject, body string) error {
	addr := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")

	msg := []byte(fmt.Sprintf("Subject: FAF Hub: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s", subject, body))

	err := smtp.SendMail(addr, repo.auth, os.Getenv("SMTP_MAIL"), to, msg)
	if err != nil {
		return err
	}

	return nil
}
