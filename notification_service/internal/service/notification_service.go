package service

import (
	"strings"

	"github.com/catness812/faf-hub-backend/notification_service/internal/util"
	"github.com/gookit/slog"
)

type INotificationRepository interface {
	SendMail(to []string, subject, body string) error
}

type NotificationService struct {
	notificationRepository INotificationRepository
}

func NewNotificationService(
	notificationRepo INotificationRepository,
) *NotificationService {
	return &NotificationService{
		notificationRepository: notificationRepo,
	}
}

func (svc *NotificationService) SendNotificationMail(msg string) {
	parts := strings.SplitN(msg, ";", 3)
	if len(parts) != 3 {
		slog.Errorf("Invalid message format: %s", msg)
		return
	}

	recipients := strings.Trim(parts[0], "[]")
	subject := parts[1]
	body := parts[2]

	to := strings.Split(recipients, ", ")
	for i := range to {
		to[i] = strings.TrimSpace(to[i])
	}

	if err := svc.notificationRepository.SendMail(to, subject, util.FormatMailMessage(body, "notification.html")); err != nil {
		slog.Fatalf("Failed to send message: %v", err)
	}

	slog.Info("Successfully sent message")
}

func (svc *NotificationService) SendVerificationMail(msg string) {
	parts := strings.Split(msg, ";")
	if len(parts) != 2 {
		slog.Fatalf("Invalid message format: %v", msg)
		return
	}

	to := []string{strings.TrimSpace(parts[0])}
	body := strings.TrimSpace(parts[1])

	if err := svc.notificationRepository.SendMail(to, "Verify your email", util.FormatMailMessage(body, "verification.html")); err != nil {
		slog.Fatalf("Failed to send message: %v", err)
	}

	slog.Info("Successfully sent message")
}
