package main

import (
	"os"

	"github.com/catness812/faf-hub-backend/gateway/internal/config"
	"github.com/catness812/faf-hub-backend/gateway/internal/content"
	"github.com/catness812/faf-hub-backend/gateway/internal/event"
	"github.com/catness812/faf-hub-backend/gateway/internal/notification"
	"github.com/catness812/faf-hub-backend/gateway/internal/repository"
	"github.com/catness812/faf-hub-backend/gateway/internal/service"
	"github.com/catness812/faf-hub-backend/gateway/internal/user"
	"github.com/catness812/faf-hub-backend/gateway/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		slog.Error("Error loading .env file:", err)
	}
	config.GoogleConfig()
}

func main() {
	redisDB := redis.Connect()
	redisRepo := repository.NewRedisRepository(redisDB)
	redisSvc := service.NewRedisService(redisRepo)
	app := fiber.New()

	frontendURL := os.Getenv("FRONT_URL")
	if frontendURL == "" {
		panic("FRONT_URL environment variable not set")
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowCredentials: true,
	}))

	registerRoutes(app, redisSvc)
	err := app.Listen(":" + os.Getenv("APP_PORT"))
	if err != nil {
		slog.Fatalf("Failed to start server: %v", err)
	}
}

func registerRoutes(r *fiber.App, redisSvc *service.RedisService) {
	notificationClient := notification.InitNotificationServiceClient(os.Getenv("NOTIFICATION_SVC_PORT"))
	userClient := user.InitUserServiceClient(os.Getenv("USER_SVC_PORT"))
	userCtrl := user.NewUserController(userClient, notificationClient, redisSvc)
	eventClient := event.InitEventServiceClient(os.Getenv("EVENT_SVC_PORT"))
	eventCtrl := event.NewEventController(eventClient, notificationClient, redisSvc)
	contentClient := content.InitContentServiceClient(os.Getenv("CONTENT_SVC_PORT"))
	contentCtrl := content.NewContentController(contentClient)

	user.RegisterUserRoutes(r, userCtrl, userClient)
	event.RegisterEventRoutes(r, eventCtrl, eventClient)
	content.RegisterContentRoutes(r, contentCtrl, contentClient)
}
