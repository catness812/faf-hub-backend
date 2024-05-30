package repository

import (
	"fmt"

	"github.com/catness812/faf-hub-backend/user_service/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) SaveUser(user models.User) error {
	err := repo.db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, fmt.Errorf("user doesn't exist")
		} else {
			return user, err
		}
	}
	return user, nil
}

func (repo *UserRepository) GetUserByID(userID uint) (models.User, error) {
	var user models.User
	if err := repo.db.Where("ID = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, fmt.Errorf("user doesn't exist")
		} else {
			return user, err
		}
	}
	return user, nil
}

func (repo *UserRepository) UpdateUser(user models.User) error {
	err := repo.db.Updates(&user).Error
	if err != nil {
		return err
	}
	return nil
}
