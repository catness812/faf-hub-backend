package content

import (
	"github.com/catness812/faf-hub-backend/gateway/internal/content/pb"
	"github.com/gookit/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitContentServiceClient(contentSvcPort string) pb.ContentServiceClient {
	conn, err := grpc.NewClient("content_svc"+":"+contentSvcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		slog.Fatalf("Could not connect: %v", err)
		return nil
	}

	return pb.NewContentServiceClient(conn)
}
