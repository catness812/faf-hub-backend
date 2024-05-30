package repository

import (
	"fmt"

	"github.com/catness812/faf-hub-backend/event_service/internal/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (repo *EventRepository) SaveEvent(event *models.Event) (uint, error) {
	err := repo.db.Create(event).Error
	if err != nil {
		return 0, err
	}
	return event.ID, nil
}

func (repo *EventRepository) UpdateEvent(event models.Event) error {
	err := repo.db.Updates(&event).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) GetEventByID(eventID uint) (models.Event, error) {
	var event models.Event
	if err := repo.db.Where("id = ?", eventID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return event, fmt.Errorf("event doesn't exist")
		} else {
			return event, err
		}
	}
	return event, nil
}

func (repo *EventRepository) DeleteEvent(eventID uint) error {
	var event models.Event
	err := repo.db.Unscoped().Where("id = ?", eventID).Delete(&event).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) GetEvents() ([]models.Event, error) {
	var events []models.Event
	// add pagination
	err := repo.db.Find(&events)
	if err.Error != nil {
		return nil, err.Error
	}
	return events, nil
}

// needs rework
func (repo *EventRepository) GetEventsByIDs(eventIDs []uint) ([]models.Event, error) {
	var events []models.Event
	if err := repo.db.Where("id IN ?", eventIDs).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
