package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/catness812/faf-hub-backend/event_service/internal/controller/rpc"
	"github.com/catness812/faf-hub-backend/event_service/internal/pb"
	"github.com/catness812/faf-hub-backend/event_service/internal/repository"
	"github.com/catness812/faf-hub-backend/event_service/internal/service"
	"github.com/catness812/faf-hub-backend/event_service/pkg/database/postgres"
	"github.com/gookit/slog"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		slog.Error("Error loading .env file:", err)
	}
}

func main() {
	db := postgres.LoadDatabase()
	eventRepo := repository.NewEventRepository(db)
	registrationRepo := repository.NewRegistrationRepository(db)
	eventSvc := service.NewEventService(eventRepo)
	registrationSvc := service.NewRegistrationService(registrationRepo, eventRepo)

	grpcStart(eventSvc, registrationSvc)
}

func grpcStart(eventSvc rpc.IEventService, registrationSvc rpc.IRegistrationService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", os.Getenv("EVENT_SVC_PORT")))
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
		EventService:        eventSvc,
		RegistrationService: registrationSvc,
	}

	pb.RegisterEventServiceServer(s, server)

	slog.Infof("gRPC Server listening at %v\n", lis.Addr())

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
