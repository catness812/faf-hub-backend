package repository

import (
	"fmt"

	"github.com/catness812/faf-hub-backend/event_service/internal/models"
	"gorm.io/gorm"
)

type RegistrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) *RegistrationRepository {
	return &RegistrationRepository{
		db: db,
	}
}

func (repo *RegistrationRepository) SaveRegistration(registration models.Registration) error {
	err := repo.db.Create(&registration).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *RegistrationRepository) GetUserEventRegistration(eventID uint, userID uint) (models.Registration, error) {
	var registration models.Registration
	if err := repo.db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&registration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return registration, fmt.Errorf("registration for this event doesn't exist")
		} else {
			return registration, err
		}
	}
	return registration, nil
}

func (repo *RegistrationRepository) GetEventRegistrations(eventID uint) ([]models.Registration, error) {
	var registrations []models.Registration
	// add pagination
	if err := repo.db.Where("event_id = ?", eventID).Find(&registrations).Error; err != nil {
		return nil, err
	}
	return registrations, nil
}

func (repo *RegistrationRepository) GetUserEventIDs(userID uint) ([]uint, error) {
	var eventIDs []uint
	if err := repo.db.Model(&models.Registration{}).Where("user_id = ?", userID).Pluck("event_id", &eventIDs).Error; err != nil {
		return nil, err
	}
	return eventIDs, nil
}

func (repo *RegistrationRepository) UpdateRegistration(registration models.Registration) error {
	err := repo.db.Updates(&registration).Error
	if err != nil {
		return err
	}
	return nil
}
