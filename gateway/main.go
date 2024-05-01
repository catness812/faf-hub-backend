package main

import (
	"os"

	"github.com/catness812/faf-hub-backend/gateway/internal/config"
	"github.com/catness812/faf-hub-backend/gateway/internal/event"
	"github.com/catness812/faf-hub-backend/gateway/internal/user"
	"github.com/catness812/faf-hub-backend/gateway/internal/user/repository"
	"github.com/catness812/faf-hub-backend/gateway/internal/user/service"
	"github.com/catness812/faf-hub-backend/gateway/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("Error loading .env file:", err)
	}
	config.GoogleConfig()
}

func main() {
	redisDB := redis.Connect()
	userRepo := repository.NewRedisRepository(redisDB)
	userSvc := service.NewUserService(userRepo)
	r := fiber.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("APP_URL"),
		AllowCredentials: true,
	}))
	// r.Use(middleware.RateLimiterMiddleware())
	registerRoutes(r, userSvc)
	err := r.Listen(":" + os.Getenv("APP_PORT"))
	if err != nil {
		slog.Fatalf("Failed to start server: %v", err)
	}
}

func registerRoutes(r *fiber.App, userSvc *service.UserService) {
	userClient := user.InitUserServiceClient(os.Getenv("USER_SVC_PORT"))
	userCtrl := user.NewUserController(userClient, userSvc)
	eventClient := event.InitEventServiceClient(os.Getenv("EVENT_SVC_PORT"))
	eventCtrl := event.NewEventController(eventClient)

	user.RegisterUserRoutes(r, userCtrl, userClient)
	event.RegisterEventRoutes(r, eventCtrl, eventClient)
}
