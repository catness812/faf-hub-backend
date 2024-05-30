package rpc

import (
	"context"

	"github.com/catness812/faf-hub-backend/notification_service/internal/controller"
	"github.com/catness812/faf-hub-backend/notification_service/internal/pb"
	"github.com/gookit/slog"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/types/known/emptypb"
)

type INotificationService interface {
	SendNotificationMail(msg string)
	SendVerificationMail(msg string)
}

type Server struct {
	pb.NotificationServiceServer
	NotificationService INotificationService
	Consumer            *controller.Consumer
}

func (s *Server) NotificationMail(name string) {
	q, err := s.Consumer.Channel.QueueDeclare(name, false, false, false, false, nil)
	if err != nil {
		slog.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := s.Consumer.Channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		slog.Panic(err)
	}

	slog.Infof("Consumer '%s' started", name)
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			s.NotificationService.SendNotificationMail(string(msg.Body))
		}
	}()

	<-forever
}

func (s *Server) VerificationMail(name string) {
	q, err := s.Consumer.Channel.QueueDeclare(name, false, false, false, false, nil)
	if err != nil {
		slog.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := s.Consumer.Channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		slog.Panic(err)
	}

	slog.Infof("Consumer '%s' started", name)
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			s.NotificationService.SendVerificationMail(string(msg.Body))
		}
	}()

	<-forever
}

func (s *Server) Publish(_ context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	if err := s.Consumer.Channel.Publish(
		"",            // exchange
		req.QueueName, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(req.Body),
		}); err != nil {
		slog.Fatalf("Failed to publish message: %v", err)
		return nil, err
	}
	slog.Infof("Successfully published message")
	return nil, nil
}
