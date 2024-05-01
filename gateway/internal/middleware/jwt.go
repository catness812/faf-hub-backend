package middleware

import (
	"net/http"

	"github.com/catness812/faf-hub-backend/gateway/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func JWTAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := util.ValidateJWT(ctx)
		if err != nil {
			slog.Errorf("Could not validate authentication: %v", err)
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "authentication required"})
		}
		return ctx.Next()
	}
}
