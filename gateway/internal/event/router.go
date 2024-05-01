package event

import (
	"github.com/catness812/faf-hub-backend/gateway/internal/event/pb"
	"github.com/catness812/faf-hub-backend/gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterEventRoutes(r *fiber.App, eventCtrl *EventController, eventClient pb.EventServiceClient) {
	route := r.Group("/event")
	route.Post("/create", middleware.JWTAuth(), middleware.CheckAdmin(), eventCtrl.CreateEvent)
	route.Post("/edit", middleware.JWTAuth(), middleware.CheckAdmin(), eventCtrl.EditEvent)
	route.Delete("/delete/:id", middleware.JWTAuth(), middleware.CheckAdmin(), eventCtrl.DeleteEvent)
	route.Get("/all", eventCtrl.GetAllEvents)
	route.Get("/:id", eventCtrl.GetEvent)
	route.Post("/register/:id", middleware.JWTAuth(), middleware.CheckIfUser(), eventCtrl.RegisterForEvent)
	route.Get("/registrations/:id", middleware.JWTAuth(), middleware.CheckAdmin(), eventCtrl.GetEventRegistrations)
}
