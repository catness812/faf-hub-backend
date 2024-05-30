package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/catness812/faf-hub-backend/notification_service/internal/controller"
	"github.com/catness812/faf-hub-backend/notification_service/internal/controller/rpc"
	"github.com/catness812/faf-hub-backend/notification_service/internal/pb"
	"github.com/catness812/faf-hub-backend/notification_service/internal/repository"
	"github.com/catness812/faf-hub-backend/notification_service/internal/service"
	"github.com/catness812/faf-hub-backend/notification_service/pkg/rabbitMQ"
	sm "github.com/catness812/faf-hub-backend/notification_service/pkg/smtp"
	"github.com/gookit/slog"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"

	"github.com/joho/godotenv"
)

var ch *amqp.Channel

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		slog.Error("Error loading .env file:", err)
	}
}

func main() {
	time.Sleep(5 * time.Second)
	defer func() {
		if r := recover(); r != nil {
			time.Sleep(time.Second * 30)
			slog.Infof("Recovered. Error:\t", r)
			ch = rabbitMQ.ConnectAMQP(os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASS"), os.Getenv("RABBITMQ_HOST"), os.Getenv("AMQP_PORT"))
		}
	}()

	ch = rabbitMQ.ConnectAMQP(os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASS"), os.Getenv("RABBITMQ_HOST"), os.Getenv("AMQP_PORT"))
	auth := sm.SmtpAuth(os.Getenv("SMTP_MAIL"), os.Getenv("SMTP_PASS"))
	notificationRepo := repository.NewNotificationRepository(auth)
	notificationSvc := service.NewNotificationService(notificationRepo)
	consumer := controller.NewConsumer(ch)

	grpcStart(notificationSvc, consumer)
}

func grpcStart(notificationSvc rpc.INotificationService, consumer *controller.Consumer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", os.Getenv("NOTIFICATION_SVC_PORT")))
	if err != nil {
		slog.Error(err)
		panic(err)
	}

	srvMetrics := srvMetrics()
	s := grpc.NewServer(
		grpc.StreamInterceptor(srvMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(srvMetrics.UnaryServerInterceptor()),
	)
	server := &rpc.Server{
		NotificationService: notificationSvc,
		Consumer:            consumer,
	}

	pb.RegisterNotificationServiceServer(s, server)

	slog.Infof("gRPC Server listening at %v\n", lis.Addr())

	go server.NotificationMail("notification")
	go server.VerificationMail("verification")

	if err := s.Serve(lis); err != nil {
		slog.Error(err)
		panic(err)
	}
}

func srvMetrics() *grpcprom.ServerMetrics {
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	return srvMetrics
}
