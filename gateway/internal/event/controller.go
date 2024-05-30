package event

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/catness812/faf-hub-backend/gateway/internal/event/pb"
	pb2 "github.com/catness812/faf-hub-backend/gateway/internal/notification/pb"
	"github.com/catness812/faf-hub-backend/gateway/internal/util"
	"github.com/catness812/faf-hub-backend/gateway/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gookit/slog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IRedisService interface {
	GetNewsletterEmails() ([]string, error)
}

type EventController struct {
	client             pb.EventServiceClient
	notificationClient pb2.NotificationServiceClient
	redisSvc           IRedisService
}

func NewEventController(client pb.EventServiceClient, notificationClient pb2.NotificationServiceClient, redisSvc IRedisService) *EventController {
	return &EventController{
		client:             client,
		notificationClient: notificationClient,
		redisSvc:           redisSvc,
	}
}

func (ctrl *EventController) CreateEvent(ctx *fiber.Ctx) error {
	var event models.Event

	if err := ctx.BodyParser(&event); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.CreateEvent(c, &pb.CreateEventRequest{
		Event: &pb.Event{
			Name:     event.Name,
			Start:    timestamppb.New(event.StartDateTime),
			End:      timestamppb.New(event.EndDateTime),
			Location: event.Location,
			Deadline: timestamppb.New(event.ApplicationDeadline),
			Cover:    event.Cover,
			Desc:     event.Description,
		},
	})

	if err != nil {
		slog.Errorf("Error creating event: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	emails, err := ctrl.redisSvc.GetNewsletterEmails()
	if err != nil {
		slog.Errorf("Error getting newsletter emails: %v", err.Error())
	}

	body := "[" + strings.Join(emails, ", ") + "]" + ";" + event.Name + ";" + event.Description

	_, err = ctrl.notificationClient.Publish(c, &pb2.PublishRequest{
		QueueName: "notification",
		Body:      body,
	})
	if err != nil {
		slog.Errorf("Error publishing newsletter emails: %v", err.Error())
	}

	slog.Info("Event created successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message, "event_id": res.EventId})
}

func (ctrl *EventController) EditEvent(ctx *fiber.Ctx) error {
	type NewEvent struct {
		EventID uint         `json:"event_id"`
		Event   models.Event `json:"event"`
	}

	var event NewEvent

	if err := ctx.BodyParser(&event); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.EditEvent(c, &pb.EditEventRequest{
		Event: &pb.Event{
			EventId:  int32(event.EventID),
			Name:     event.Event.Name,
			Start:    timestamppb.New(event.Event.StartDateTime),
			End:      timestamppb.New(event.Event.EndDateTime),
			Location: event.Event.Location,
			Deadline: timestamppb.New(event.Event.ApplicationDeadline),
			Cover:    event.Event.Cover,
			Desc:     event.Event.Description,
		},
	})

	if err != nil {
		slog.Errorf("Error updating event: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("Event updated successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl *EventController) DeleteEvent(ctx *fiber.Ctx) error {
	sid := ctx.Params("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		slog.Errorf("Error converting string to int: %v", err)
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.DeleteEvent(c, &pb.DeleteEventRequest{
		EventId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error deleting event: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("Event deleted successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl *EventController) GetEvent(ctx *fiber.Ctx) error {
	sid := ctx.Params("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		slog.Errorf("Error converting string to int: %v", err)
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetEvent(c, &pb.GetEventRequest{
		EventId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error retrieving event: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	event := models.Event{
		Name:                res.Event.Name,
		StartDateTime:       res.Event.Start.AsTime(),
		EndDateTime:         res.Event.End.AsTime(),
		Location:            res.Event.Location,
		ApplicationDeadline: res.Event.Deadline.AsTime(),
		Cover:               res.Event.Cover,
		Description:         res.Event.Desc,
	}
	event.ID = uint(res.Event.EventId)

	slog.Info("Event retrieved successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"event": event})
}

func (ctrl *EventController) GetAllEvents(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetEvents(c, &empty.Empty{})

	if err != nil {
		slog.Errorf("Error retrieving events: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	events := make([]models.Event, len(res.Events))

	for i := range events {
		events[i] = models.Event{
			Name:                res.Events[i].Name,
			StartDateTime:       res.Events[i].Start.AsTime(),
			EndDateTime:         res.Events[i].End.AsTime(),
			Location:            res.Events[i].Location,
			ApplicationDeadline: res.Events[i].Deadline.AsTime(),
			Cover:               res.Events[i].Cover,
			Description:         res.Events[i].Desc,
		}
		events[i].ID = uint(res.Events[i].EventId)
	}

	slog.Info("Events retrieved successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"events": events})
}

func (ctrl *EventController) RegisterForEvent(ctx *fiber.Ctx) error {
	var reg models.Registration

	sid := ctx.Params("id")
	event_id, err := strconv.Atoi(sid)
	if err != nil {
		slog.Errorf("Error converting string to int: %v", err)
		return err
	}

	user_id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ctx.BodyParser(&reg); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.RegisterForEvent(c, &pb.RegisterForEventRequest{
		Registration: &pb.EventRegistration{
			EventId:       int32(event_id),
			UserId:        int32(user_id),
			FirstName:     reg.FirstName,
			LastName:      reg.LastName,
			Email:         reg.Email,
			PhoneNumber:   int32(reg.PhoneNumber),
			AcademicGroup: reg.AcademicGroup,
			TeamMembers:   reg.TeamMembers,
			ShirtSize:     reg.ShirtSize,
			FoodPref:      reg.FoodPreferences,
			Motivation:    reg.Motivation,
			Questions:     reg.Questions,
		}})

	if err != nil {
		slog.Errorf("Error registering for event: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("Event registration completed successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl *EventController) GetEventRegistrations(ctx *fiber.Ctx) error {
	sid := ctx.Params("id")
	event_id, err := strconv.Atoi(sid)
	if err != nil {
		slog.Errorf("Error converting string to int: %v", err)
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetEventRegistrations(c, &pb.GetEventRegistrationsRequest{
		EventId: int32(event_id),
	})

	if err != nil {
		slog.Errorf("Error retrieving event registrations: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("Event registrations retrieved successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"registrations": res.Registrations})
}

func (ctrl EventController) GetUserEventRegistration(ctx *fiber.Ctx) error {
	sid := ctx.Params("id")
	event_id, err := strconv.Atoi(sid)
	if err != nil {
		slog.Errorf("Error converting string to int: %v", err)
		return err
	}

	user_id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetEventUserRegistration(c, &pb.GetEventUserRegistrationRequest{
		EventId: int32(event_id),
		UserId:  int32(user_id),
	})

	if err != nil {
		slog.Errorf("Error retrieving user event registration: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("User event registration retrieved successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"registration": res.Registration})
}

func (ctrl EventController) GetUserEvents(ctx *fiber.Ctx) error {
	user_id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetUserEvents(c, &pb.GetUserEventsRequest{
		UserId: int32(user_id),
	})

	if err != nil {
		slog.Errorf("Error retrieving events user registered for: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	events := make([]models.Event, len(res.Events))

	for i := range events {
		events[i] = models.Event{
			Name:                res.Events[i].Name,
			StartDateTime:       res.Events[i].Start.AsTime(),
			EndDateTime:         res.Events[i].End.AsTime(),
			Location:            res.Events[i].Location,
			ApplicationDeadline: res.Events[i].Deadline.AsTime(),
			Cover:               res.Events[i].Cover,
			Description:         res.Events[i].Desc,
		}
		events[i].ID = uint(res.Events[i].EventId)
	}

	slog.Info("Events user registered for retrieved successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"events": events})
}

func (ctrl EventController) EditRegistration(ctx *fiber.Ctx) error {
	var registration models.Registration

	if err := ctx.BodyParser(&registration); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user_id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.EditRegistration(c, &pb.EditRegistrationRequest{
		Registration: &pb.EventRegistration{
			EventId:       int32(registration.EventID),
			UserId:        int32(user_id),
			FirstName:     registration.FirstName,
			LastName:      registration.LastName,
			Email:         registration.Email,
			PhoneNumber:   int32(registration.PhoneNumber),
			AcademicGroup: registration.AcademicGroup,
			TeamMembers:   registration.TeamMembers,
			ShirtSize:     registration.ShirtSize,
			FoodPref:      registration.FoodPreferences,
			Motivation:    registration.Motivation,
			Questions:     registration.Questions,
		},
	})

	if err != nil {
		slog.Errorf("Error updating user event registration: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("User event registration updated successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl EventController) LeaveFeedback(ctx *fiber.Ctx) error {
	type Feedback struct {
		EventID  uint   `json:"event_id"`
		Feedback string `json:"feedback"`
	}

	var feedback Feedback

	if err := ctx.BodyParser(&feedback); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user_id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	_, err = ctrl.client.EditRegistration(c, &pb.EditRegistrationRequest{
		Registration: &pb.EventRegistration{
			EventId:  int32(feedback.EventID),
			UserId:   int32(user_id),
			Feedback: feedback.Feedback,
		},
	})

	if err != nil {
		slog.Errorf("Error updating user event registration: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("Event feedback added successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "event feedback added successfully"})
}
