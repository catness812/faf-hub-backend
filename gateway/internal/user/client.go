package user

import (
	"os"

	"github.com/catness812/faf-hub-backend/gateway/internal/user/pb"
	"github.com/gookit/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserServiceClient(userSvcPort string) pb.UserServiceClient {
	conn, err := grpc.Dial(os.Getenv("APP_HOST")+":"+userSvcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		slog.Fatalf("Could not connect: %v", err)
		return nil
	}

	return pb.NewUserServiceClient(conn)
}
