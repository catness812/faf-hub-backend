package user

import (
	"github.com/catness812/faf-hub-backend/gateway/internal/middleware"
	"github.com/catness812/faf-hub-backend/gateway/internal/user/pb"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(r *fiber.App, userCtrl *UserController, userClient pb.UserServiceClient) {
	route := r.Group("/user")
	route.Post("/register", userCtrl.CreateUser)
	route.Post("/login", userCtrl.Login)
	route.Get("", middleware.JWTAuth(), userCtrl.GetUser)
	route.Get("/google_auth", userCtrl.GoogleAuth)
	route.Get("/google_callback", userCtrl.GoogleCallback)
	route.Post("/update", middleware.JWTAuth(), userCtrl.UpdateUser)
	route.Get("/logout", middleware.JWTAuth(), userCtrl.Logout)
}
