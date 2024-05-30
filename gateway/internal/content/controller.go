package content

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/catness812/faf-hub-backend/gateway/internal/content/pb"
	"github.com/catness812/faf-hub-backend/gateway/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type ContentController struct {
	client pb.ContentServiceClient
}

func NewContentController(client pb.ContentServiceClient) *ContentController {
	return &ContentController{
		client: client,
	}
}

func (ctrl *ContentController) PostContent(ctx *fiber.Ctx) error {
	var content models.Content

	if err := ctx.BodyParser(&content); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if content.Type != "article" && content.Type != "project" {
		slog.Error("invalid content type")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid content type"})
	}

	var images []string
	if content.Images != "" {
		imagesStr := strings.ReplaceAll(content.Images, " ", "")
		images = strings.Split(imagesStr, ",")
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.PostContent(c, &pb.PostContentRequest{
		Content: &pb.Content{
			Type:    content.Type,
			Name:    content.Name,
			Authors: content.Authors,
			Cover:   content.Cover,
			Text:    content.Text,
			Images:  images,
		},
	})

	if err != nil {
		slog.Errorf("Error posting content: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("Content posted successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message, "content_id": res.ContentId})
}

func (ctrl *ContentController) EditContent(ctx *fiber.Ctx) error {
	type NewContent struct {
		ContentID uint           `json:"content_id"`
		Content   models.Content `json:"content"`
	}

	var content NewContent

	if err := ctx.BodyParser(&content); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var images []string
	if content.Content.Images != "" {
		imagesStr := strings.ReplaceAll(content.Content.Images, " ", "")
		images = strings.Split(imagesStr, ",")
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.EditContent(c, &pb.EditContentRequest{
		ContentId: int32(content.ContentID),
		Name:      content.Content.Name,
		Authors:   content.Content.Authors,
		Cover:     content.Content.Cover,
		Text:      content.Content.Text,
		Images:    images,
	})

	if err != nil {
		slog.Errorf("Error editting content: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("Content editted successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl *ContentController) DeleteContent(ctx *fiber.Ctx) error {
	sid := ctx.Params("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		slog.Errorf("Error converting string to int: %v", err)
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.DeleteContent(c, &pb.DeleteContentRequest{
		ContentId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error deleting content: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("Content deleted successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl *ContentController) GetContent(ctx *fiber.Ctx) error {
	sid := ctx.Params("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		slog.Errorf("Error converting string to int: %v", err)
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetContent(c, &pb.GetContentRequest{
		ContentId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error retrieving content: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var images []string
	if err := json.Unmarshal([]byte(strings.Join(res.Content.Images, ", ")), &images); err != nil {
		slog.Errorf("Error JSON unmarshal: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	content := models.Content{
		Type:    res.Content.Type,
		Name:    res.Content.Name,
		Authors: res.Content.Authors,
		Cover:   res.Content.Cover,
		Text:    res.Content.Text,
		Views:   int(res.Content.Views),
		Images:  strings.Join(images, ", "),
	}
	content.ID = uint(id)
	content.CreatedAt = res.Content.Date.AsTime()

	slog.Info("Content retrieved successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"content": content})
}

func (ctrl *ContentController) GetArticles(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetAllContent(c, &pb.GetAllContentRequest{
		Type: "article",
	})

	if err != nil {
		slog.Errorf("Error retrieving articles: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	articles := make([]models.Content, len(res.Content))

	for i := range articles {
		var images []string
		if err := json.Unmarshal([]byte(strings.Join(res.Content[i].Images, ", ")), &images); err != nil {
			slog.Errorf("Error JSON unmarshal: %v", err.Error())
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		articles[i] = models.Content{
			Name:    res.Content[i].Name,
			Authors: res.Content[i].Authors,
			Cover:   res.Content[i].Cover,
			Text:    res.Content[i].Text,
			Views:   int(res.Content[i].Views),
			Images:  strings.Join(images, ", "),
		}
		articles[i].ID = uint(res.Content[i].ContentId)
		articles[i].CreatedAt = res.Content[i].Date.AsTime()
	}

	slog.Info("Articles retrieved successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"articles": articles})
}

func (ctrl *ContentController) GetProjects(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetAllContent(c, &pb.GetAllContentRequest{
		Type: "project",
	})

	if err != nil {
		slog.Errorf("Error retrieving projects: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	projects := make([]models.Content, len(res.Content))

	for i := range projects {
		var images []string
		if err := json.Unmarshal([]byte(strings.Join(res.Content[i].Images, ", ")), &images); err != nil {
			slog.Errorf("Error JSON unmarshal: %v", err.Error())
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		projects[i] = models.Content{
			Name:    res.Content[i].Name,
			Authors: res.Content[i].Authors,
			Cover:   res.Content[i].Cover,
			Text:    res.Content[i].Text,
			Views:   int(res.Content[i].Views),
			Images:  strings.Join(images, ", "),
		}
		projects[i].ID = uint(res.Content[i].ContentId)
		projects[i].CreatedAt = res.Content[i].Date.AsTime()
	}

	slog.Info("Projects retrieved successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"projects": projects})
}
