package main

import (
	"fmt"
	"net"
	"os"

	"github.com/catness812/faf-hub-backend/user_service/internal/controller/rpc"
	"github.com/catness812/faf-hub-backend/user_service/internal/pb"
	"github.com/catness812/faf-hub-backend/user_service/internal/repository"
	"github.com/catness812/faf-hub-backend/user_service/internal/service"
	"github.com/catness812/faf-hub-backend/user_service/pkg/database/postgres"
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
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	grpcStart(userSvc)
}

func grpcStart(userSvc rpc.IUserService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", os.Getenv("USER_SVC_PORT")))
	if err != nil {
		slog.Error(err)
		panic(err)
	}

	s := grpc.NewServer()
	server := &rpc.Server{
		UserService: userSvc,
	}

	pb.RegisterUserServiceServer(s, server)

	slog.Infof("gRPC Server listening at %v\n", lis.Addr())

	if err := s.Serve(lis); err != nil {
		slog.Error(err)
		panic(err)
	}
}
