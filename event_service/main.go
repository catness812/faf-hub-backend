package main

import (
	"fmt"
	"net"
	"os"

	"github.com/catness812/faf-hub-backend/event_service/internal/controller/rpc"
	"github.com/catness812/faf-hub-backend/event_service/internal/pb"
	"github.com/catness812/faf-hub-backend/event_service/internal/repository"
	"github.com/catness812/faf-hub-backend/event_service/internal/service"
	"github.com/catness812/faf-hub-backend/event_service/pkg/database/postgres"
	"github.com/gookit/slog"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	err := godotenv.Load("../.env")
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

	s := grpc.NewServer()
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
