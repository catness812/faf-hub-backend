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
