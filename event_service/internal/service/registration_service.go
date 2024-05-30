package service

import (
	"fmt"

	"github.com/catness812/faf-hub-backend/event_service/internal/models"
	"github.com/gookit/slog"
)

type IRegistrationRepository interface {
	SaveRegistration(registration models.Registration) error
	GetUserEventRegistration(eventID uint, userID uint) (models.Registration, error)
	GetEventRegistrations(eventID uint) ([]models.Registration, error)
	GetUserEventIDs(userID uint) ([]uint, error)
	UpdateRegistration(registration models.Registration) error
}

type RegistrationService struct {
	registrationRepository IRegistrationRepository
	eventRepository        IEventRepository
}

func NewRegistrationService(
	registrationRepo IRegistrationRepository,
	eventRepository IEventRepository,
) *RegistrationService {
	return &RegistrationService{
		registrationRepository: registrationRepo,
		eventRepository:        eventRepository,
	}
}

func (svc *RegistrationService) RegisterForEvent(registration models.Registration) error {
	if _, err := svc.eventRepository.GetEventByID(registration.EventID); err != nil {
		slog.Errorf("Could not retrieve event: %v", err)
		return err
	}

	if _, err := svc.registrationRepository.GetUserEventRegistration(registration.EventID, registration.UserID); err == nil {
		slog.Error("Registration already exists")
		return fmt.Errorf("registration already exists")
	}

	if err := svc.registrationRepository.SaveRegistration(registration); err != nil {
		slog.Errorf("Could not save registration: %v", err)
		return err
	}

	slog.Info("Event registration successfully saved")
	return nil
}

func (svc *RegistrationService) GetEventRegistrations(eventID uint) ([]models.Registration, error) {
	if _, err := svc.eventRepository.GetEventByID(eventID); err != nil {
		slog.Errorf("Could not retrieve event: %v", err)
		return nil, err
	}

	regs, err := svc.registrationRepository.GetEventRegistrations(eventID)
	if err != nil {
		slog.Errorf("Could not retrieve event registrations: %v", err)
		return nil, err
	}

	slog.Info("Event registrations successfully retrieved")
	return regs, nil
}

func (svc *RegistrationService) GetUserEventRegistration(eventID uint, userID uint) (models.Registration, error) {
	registration, err := svc.registrationRepository.GetUserEventRegistration(eventID, userID)
	if err != nil {
		slog.Errorf("Could not retrieve user event registration: %v", err)
		return registration, err
	}

	slog.Info("User event registration successfully retrieved")
	return registration, err
}

func (svc *RegistrationService) GetUserEvents(userID uint) ([]models.Event, error) {
	eventsIDs, err := svc.registrationRepository.GetUserEventIDs(userID)
	if err != nil {
		slog.Errorf("Could not retrieve event ids user registered for: %v", err)
		return nil, err
	}

	events, err := svc.eventRepository.GetEventsByIDs(eventsIDs)
	if err != nil {
		slog.Errorf("Could not retrieve events user registered for: %v", err)
		return nil, err
	}

	slog.Info("Events user registered for successfully retrieved")
	return events, nil
}

func (svc *RegistrationService) UpdateRegistration(registration models.Registration) error {
	newRegistration, err := svc.registrationRepository.GetUserEventRegistration(registration.EventID, registration.UserID)
	if err != nil {
		slog.Errorf("Could not retrieve user registration: %v", err)
		return err
	}

	if registration.FirstName != "" {
		newRegistration.FirstName = registration.FirstName
	}
	if registration.LastName != "" {
		newRegistration.LastName = registration.LastName
	}
	if registration.Email != "" {
		newRegistration.Email = registration.Email
	}
	if registration.PhoneNumber != 0 {
		newRegistration.PhoneNumber = registration.PhoneNumber
	}
	if registration.AcademicGroup != "" {
		newRegistration.AcademicGroup = registration.AcademicGroup
	}
	if registration.TeamMembers != "" {
		newRegistration.TeamMembers = registration.TeamMembers
	}
	if registration.ShirtSize != "" {
		newRegistration.ShirtSize = registration.ShirtSize
	}
	if registration.FoodPreferences != "" {
		newRegistration.FoodPreferences = registration.FoodPreferences
	}
	if registration.Motivation != "" {
		newRegistration.Motivation = registration.Motivation
	}
	if registration.Questions != "" {
		newRegistration.Questions = registration.Questions
	}
	if registration.Feedback != "" {
		newRegistration.Feedback = registration.Feedback
	}

	if err := svc.registrationRepository.UpdateRegistration(newRegistration); err != nil {
		slog.Errorf("Could not update user registration: %v", err)
		return err
	}

	slog.Info("User registration successfully updated")
	return nil
}
