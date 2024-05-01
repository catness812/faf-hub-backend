package event

import (
	"os"

	"github.com/catness812/faf-hub-backend/gateway/internal/event/pb"
	"github.com/gookit/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitEventServiceClient(eventSvcPort string) pb.EventServiceClient {
	conn, err := grpc.Dial(os.Getenv("APP_HOST")+":"+eventSvcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		slog.Fatalf("Could not connect: %v", err)
		return nil
	}

	return pb.NewEventServiceClient(conn)
}
