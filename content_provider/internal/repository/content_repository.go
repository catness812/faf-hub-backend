package repository

import (
	"fmt"

	"github.com/catness812/faf-hub-backend/content_provider/internal/models"
	"gorm.io/gorm"
)

type ContentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{
		db: db,
	}
}

func (repo *ContentRepository) SaveContent(content *models.Content) (uint, error) {
	err := repo.db.Create(content).Error
	if err != nil {
		return 0, err
	}
	return content.ID, nil
}

func (repo *ContentRepository) UpdateContent(content models.Content) error {
	err := repo.db.Updates(&content).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *ContentRepository) UpdateViews(contentID uint) error {
	if err := repo.db.Model(&models.Content{}).Where("id = ?", contentID).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ContentRepository) GetContentByID(contentID uint) (models.Content, error) {
	var content models.Content
	if err := repo.db.Where("id = ?", contentID).First(&content).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return content, fmt.Errorf("content doesn't exist")
		} else {
			return content, err
		}
	}
	return content, nil
}

func (repo *ContentRepository) DeleteContent(contentID uint) error {
	var content models.Content
	err := repo.db.Unscoped().Where("id = ?", contentID).Delete(&content).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *ContentRepository) GetContentByType(contentType string) ([]models.Content, error) {
	var content []models.Content
	if err := repo.db.Where("type = ?", contentType).Find(&content).Error; err != nil {
		return nil, err
	}
	return content, nil
}
