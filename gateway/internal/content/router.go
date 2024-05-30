package content

import (
	"github.com/catness812/faf-hub-backend/gateway/internal/content/pb"
	"github.com/catness812/faf-hub-backend/gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterContentRoutes(r *fiber.App, contentCtrl *ContentController, contentClient pb.ContentServiceClient) {
	route := r.Group("/content")
	route.Get("/articles", contentCtrl.GetArticles)
	route.Get("/projects", contentCtrl.GetProjects)
	route.Post("/post", middleware.JWTAuth(), middleware.CheckAdmin(), middleware.CheckIfVerified(), contentCtrl.PostContent)
	route.Post("/edit", middleware.JWTAuth(), middleware.CheckAdmin(), middleware.CheckIfVerified(), contentCtrl.EditContent)
	route.Delete("/delete/:id", middleware.JWTAuth(), middleware.CheckAdmin(), middleware.CheckIfVerified(), contentCtrl.DeleteContent)
	route.Get("/:id", contentCtrl.GetContent)
}
