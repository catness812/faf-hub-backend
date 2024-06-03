package rpc

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/catness812/faf-hub-backend/content_provider/internal/models"
	"github.com/catness812/faf-hub-backend/content_provider/internal/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IContentService interface {
	AddContent(content models.Content) (uint, error)
	GetContent(contentID uint) (models.Content, error)
	UpdateContent(contentID uint, content models.Content) error
	DeleteContent(contentID uint) error
	IncrementViews(contentID uint) error
	GetContentByType(contentType string) ([]models.Content, error)
}

type Server struct {
	pb.ContentServiceServer
	ContentService IContentService
}

func (s *Server) PostContent(_ context.Context, req *pb.PostContentRequest) (*pb.PostContentResponse, error) {
	imagesJSON, err := json.Marshal(req.Content.Images)
	if err != nil {
		return nil, err
	}

	newContent := models.Content{
		Type:    req.Content.Type,
		Name:    req.Content.Name,
		Authors: req.Content.Authors,
		Cover:   req.Content.Cover,
		Text:    req.Content.Text,
		Views:   int(req.Content.Views),
		Images:  string(imagesJSON),
	}

	contentID, err := s.ContentService.AddContent(newContent)
	if err != nil {
		return nil, err
	}

	return &pb.PostContentResponse{
		Message:   "content posted successfully",
		ContentId: int32(contentID),
	}, nil
}

func (s *Server) EditContent(_ context.Context, req *pb.EditContentRequest) (*pb.EditContentResponse, error) {
	imagesJSON, err := json.Marshal(req.Images)
	if err != nil {
		return nil, err
	}

	content := models.Content{
		Name:    req.Name,
		Authors: req.Authors,
		Cover:   req.Cover,
		Text:    req.Text,
		Images:  string(imagesJSON),
	}

	if err := s.ContentService.UpdateContent(uint(req.ContentId), content); err != nil {
		return nil, err
	}

	return &pb.EditContentResponse{
		Message: "content editted successfully",
	}, nil
}

func (s *Server) DeleteContent(_ context.Context, req *pb.DeleteContentRequest) (*pb.DeleteContentResponse, error) {
	if err := s.ContentService.DeleteContent(uint(req.ContentId)); err != nil {
		return nil, err
	}

	return &pb.DeleteContentResponse{
		Message: "content deleted successfully",
	}, nil
}

func (s *Server) GetContent(_ context.Context, req *pb.GetContentRequest) (*pb.GetContentResponse, error) {
	content, err := s.ContentService.GetContent(uint(req.ContentId))
	if err != nil {
		return nil, err
	}

	if err := s.ContentService.IncrementViews(content.ID); err != nil {
		return nil, err
	}

	return &pb.GetContentResponse{
		Content: &pb.Content{
			Type:    content.Type,
			Name:    content.Name,
			Date:    timestamppb.New(content.CreatedAt),
			Authors: content.Authors,
			Cover:   content.Cover,
			Text:    content.Text,
			Views:   int32(content.Views) + 1,
			Images:  strings.Split(content.Images, ","),
		},
	}, nil
}

func (s *Server) GetAllContent(_ context.Context, req *pb.GetAllContentRequest) (*pb.GetAllContentResponse, error) {
	content, err := s.ContentService.GetContentByType(req.Type)
	if err != nil {
		return nil, err
	}

	getContentRes := make([]*pb.Content, len(content))

	for i := range getContentRes {
		c := content[i]
		getContentRes[i] = &pb.Content{
			ContentId: int32(c.ID),
			Name:      c.Name,
			Date:      timestamppb.New(c.CreatedAt),
			Authors:   c.Authors,
			Cover:     c.Cover,
			Text:      c.Text,
			Views:     int32(c.Views),
			Images:    strings.Split(c.Images, ","),
		}
	}

	return &pb.GetAllContentResponse{
		Content: getContentRes,
	}, nil
}
