package middleware

import (
	"net/http"

	"github.com/catness812/faf-hub-backend/gateway/internal/util"
	"github.com/gofiber/fiber/v2"
)

func CheckAdmin() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		status, err := util.CheckAdmin(ctx)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err})
		} else {
			if status != true {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "not admin"})
			}
			return ctx.Next()
		}
	}
}

func CheckIfUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		status, err := util.CheckAdmin(ctx)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err})
		} else {
			if status != false {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "not user"})
			}
			return ctx.Next()
		}
	}
}
