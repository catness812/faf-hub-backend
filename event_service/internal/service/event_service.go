package service

import (
	"time"

	"github.com/catness812/faf-hub-backend/event_service/internal/models"
	"github.com/gookit/slog"
)

type IEventRepository interface {
	SaveEvent(event *models.Event) (uint, error)
	UpdateEvent(event models.Event) error
	GetEventByID(eventID uint) (models.Event, error)
	DeleteEvent(eventID uint) error
	GetEvents() ([]models.Event, error)
	GetEventsByIDs(eventIDs []uint) ([]models.Event, error)
}

type EventService struct {
	eventRepository IEventRepository
}

func NewEventService(
	eventRepo IEventRepository,
) *EventService {
	return &EventService{
		eventRepository: eventRepo,
	}
}

func (svc *EventService) CreateEvent(event models.Event) (uint, error) {
	eventID, err := svc.eventRepository.SaveEvent(&event)
	if err != nil {
		slog.Errorf("Could not create new event: %v", err)
		return 0, err
	}

	slog.Infof("Event successfully created: %s", event.Name)
	return eventID, err
}

func (svc *EventService) UpdateEvent(eventID uint, event models.Event) error {
	newEvent, err := svc.eventRepository.GetEventByID(eventID)
	if err != nil {
		slog.Errorf("Could not retrieve event: %v", err)
		return err
	}

	if event.Name != "" {
		newEvent.Name = event.Name
	}
	if !event.StartDateTime.Equal(time.Time{}) {
		newEvent.StartDateTime = event.StartDateTime
	}
	if !event.EndDateTime.Equal(time.Time{}) {
		newEvent.EndDateTime = event.EndDateTime
	}
	if event.Location != "" {
		newEvent.Location = event.Location
	}
	if !event.ApplicationDeadline.Equal(time.Time{}) {
		newEvent.ApplicationDeadline = event.ApplicationDeadline
	}
	if event.Cover != "" {
		newEvent.Cover = event.Cover
	}
	if event.Description != "" {
		newEvent.Description = event.Description
	}

	if err := svc.eventRepository.UpdateEvent(newEvent); err != nil {
		slog.Errorf("Could not update event: %v", err)
		return err
	}

	slog.Info("Event successfully updated")
	return nil
}

func (svc *EventService) DeleteEvent(eventID uint) error {
	if _, err := svc.eventRepository.GetEventByID(eventID); err != nil {
		slog.Errorf("Could not retrieve event: %v", err)
		return err
	}

	if err := svc.eventRepository.DeleteEvent(eventID); err != nil {
		slog.Errorf("Could not delete event: %v", err)
		return err
	}

	slog.Info("Event successfully deleted")
	return nil
}

func (svc *EventService) GetEventByID(eventID uint) (models.Event, error) {
	event, err := svc.eventRepository.GetEventByID(eventID)
	if err != nil {
		slog.Errorf("Could not retrieve event: %v", err)
		return event, err
	}

	slog.Info("Event successfully retrieved")
	return event, nil
}

func (svc *EventService) GetAllEvents() ([]models.Event, error) {
	events, err := svc.eventRepository.GetEvents()
	if err != nil {
		slog.Errorf("Could not retrieve events: %v", err)
		return nil, err
	}

	slog.Info("Events successfully retrieved")
	return events, nil
}
