package service

import (
	"github.com/catness812/faf-hub-backend/content_provider/internal/models"
	"github.com/gookit/slog"
)

type IContentRepository interface {
	SaveContent(content *models.Content) (uint, error)
	UpdateContent(content models.Content) error
	GetContentByID(contentID uint) (models.Content, error)
	DeleteContent(contentID uint) error
	UpdateViews(contentID uint) error
	GetContentByType(contentType string) ([]models.Content, error)
}

type ContentService struct {
	contentRepository IContentRepository
}

func NewContentService(
	contentRepo IContentRepository,
) *ContentService {
	return &ContentService{
		contentRepository: contentRepo,
	}
}

func (svc *ContentService) AddContent(content models.Content) (uint, error) {
	contentID, err := svc.contentRepository.SaveContent(&content)
	if err != nil {
		slog.Errorf("Could not add new content: %v", err)
		return 0, err
	}

	slog.Infof("Content successfully created: %s", content.Name)
	return contentID, err
}

func (svc *ContentService) GetContent(contentID uint) (models.Content, error) {
	content, err := svc.contentRepository.GetContentByID(contentID)
	if err != nil {
		slog.Errorf("Could not retrieve content: %v", err)
		return content, err
	}

	slog.Info("Content successfully retrieved")
	return content, err
}

func (svc *ContentService) UpdateContent(contentID uint, content models.Content) error {
	newContent, err := svc.contentRepository.GetContentByID(contentID)
	if err != nil {
		slog.Errorf("Could not retrieve content: %v", err)
		return err
	}

	if content.Name != "" {
		newContent.Name = content.Name
	}
	if content.Authors != "" {
		newContent.Authors = content.Authors
	}
	if content.Cover != "" {
		newContent.Cover = content.Cover
	}
	if content.Text != "" {
		newContent.Text = content.Text
	}
	if content.Images != "" {
		newContent.Images = content.Images
	}

	if err := svc.contentRepository.UpdateContent(newContent); err != nil {
		slog.Errorf("Could not update content: %v", err)
		return err
	}

	slog.Info("Content successfully updated")
	return nil
}

func (svc *ContentService) DeleteContent(contentID uint) error {
	if _, err := svc.contentRepository.GetContentByID(contentID); err != nil {
		slog.Errorf("Could not retrieve content: %v", err)
		return err
	}

	if err := svc.contentRepository.DeleteContent(contentID); err != nil {
		slog.Errorf("Could not delete content: %v", err)
		return err
	}

	slog.Info("Content successfully deleted")
	return nil
}

func (svc *ContentService) IncrementViews(contentID uint) error {
	if err := svc.contentRepository.UpdateViews(contentID); err != nil {
		slog.Errorf("Could not increment content views: %v", err)
		return err
	}

	slog.Info("Content views successfully incremented")
	return nil
}

func (svc *ContentService) GetContentByType(contentType string) ([]models.Content, error) {
	content, err := svc.contentRepository.GetContentByType(contentType)
	if err != nil {
		slog.Errorf("Could not retrieve content: %v", err)
		return nil, err
	}

	slog.Info("Content successfully retrieved")
	return content, nil
}
