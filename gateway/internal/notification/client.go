package notification

import (
	"github.com/catness812/faf-hub-backend/gateway/internal/notification/pb"
	"github.com/gookit/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitNotificationServiceClient(notificationSvcPort string) pb.NotificationServiceClient {
	conn, err := grpc.NewClient("notification_svc"+":"+notificationSvcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		slog.Fatalf("Could not connect: %v", err)
		return nil
	}

	return pb.NewNotificationServiceClient(conn)
}
