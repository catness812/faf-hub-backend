package event

import (
	"github.com/catness812/faf-hub-backend/gateway/internal/event/pb"
	"github.com/catness812/faf-hub-backend/gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterEventRoutes(r *fiber.App, eventCtrl *EventController, eventClient pb.EventServiceClient) {
	route := r.Group("/event")
	route.Post("/create", middleware.JWTAuth(), middleware.CheckAdmin(), middleware.CheckIfVerified(), eventCtrl.CreateEvent)
	route.Post("/edit", middleware.JWTAuth(), middleware.CheckAdmin(), middleware.CheckIfVerified(), eventCtrl.EditEvent)
	route.Delete("/delete/:id", middleware.JWTAuth(), middleware.CheckAdmin(), middleware.CheckIfVerified(), eventCtrl.DeleteEvent)
	route.Get("/all", eventCtrl.GetAllEvents)
	route.Get("/registered", middleware.JWTAuth(), middleware.CheckIfUser(), eventCtrl.GetUserEvents)
	route.Get("/:id/registration", middleware.JWTAuth(), middleware.CheckIfUser(), eventCtrl.GetUserEventRegistration)
	route.Get("/:id", eventCtrl.GetEvent)
	route.Post("/register/:id", middleware.JWTAuth(), middleware.CheckIfUser(), middleware.CheckIfVerified(), eventCtrl.RegisterForEvent)
	route.Get("/registrations/:id", middleware.JWTAuth(), middleware.CheckAdmin(), middleware.CheckIfVerified(), eventCtrl.GetEventRegistrations)
	route.Post("/registration/edit", middleware.JWTAuth(), middleware.CheckIfUser(), eventCtrl.EditRegistration)
	route.Post("/feedback", middleware.JWTAuth(), middleware.CheckIfUser(), eventCtrl.LeaveFeedback)
}
