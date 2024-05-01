package user

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/catness812/faf-hub-backend/gateway/internal/config"
	"github.com/catness812/faf-hub-backend/gateway/internal/user/pb"
	"github.com/catness812/faf-hub-backend/gateway/internal/util"
	"github.com/catness812/faf-hub-backend/gateway/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type IUserService interface {
	SaveJWT(userID int, jwt string) error
}

type UserController struct {
	client  pb.UserServiceClient
	userSvc IUserService
}

func NewUserController(client pb.UserServiceClient, userSvc IUserService) *UserController {
	return &UserController{
		client:  client,
		userSvc: userSvc,
	}
}

func (ctrl *UserController) CreateUser(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.CreateUser(c, &pb.CreateUserRequest{User: &pb.User{
		Email:         user.Email,
		Password:      user.Password,
		PhoneNumber:   int32(user.PhoneNumber),
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		AcademicGroup: user.AcademicGroup,
	}})

	if err != nil {
		slog.Errorf("Error creating user: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("User created successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl *UserController) Login(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.Login(c, &pb.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		slog.Errorf("Error logging user in: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	jwt, err := util.GenerateJWT(int(res.UserId))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	slog.Info("User logged in successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl *UserController) GetUser(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := ctrl.client.GetUser(c, &pb.GetUserRequest{
		UserId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error retrieving user: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userInfo := models.UserInfo{
		Email:         res.User.Email,
		PhoneNumber:   int(res.User.PhoneNumber),
		FirstName:     res.User.FirstName,
		LastName:      res.User.LastName,
		AcademicGroup: res.User.AcademicGroup,
	}

	slog.Info("User retrieved successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"user": userInfo})
}

func (ctrl *UserController) GoogleAuth(ctx *fiber.Ctx) error {
	url := config.Google.Config.AuthCodeURL("authState")

	ctx.Status(fiber.StatusSeeOther)
	ctx.Redirect(url)
	return ctx.JSON(url)
}

func (ctrl *UserController) GoogleCallback(ctx *fiber.Ctx) error {
	state := ctx.Query("state")
	if state != "authState" {
		return ctx.SendString("States do not match")
	}

	code := ctx.Query("code")

	token, err := config.Google.Config.Exchange(context.Background(), code)
	if err != nil {
		return ctx.SendString("Code-Token exchange failed")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return ctx.SendString("User data fetch failed")
	}

	type UserData struct {
		Email      string `json:"email"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.SendString("Failed to read response body")
	}

	userData := UserData{}
	if err := json.Unmarshal(body, &userData); err != nil {
		return ctx.SendString("JSON parsing failed")
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GoogleAuth(c, &pb.GoogleAuthRequest{
		Email:     userData.Email,
		FirstName: userData.GivenName,
		LastName:  userData.FamilyName,
	})

	if err != nil {
		slog.Errorf("Error logging user in: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	jwt, err := util.GenerateJWT(int(res.UserId))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	slog.Info("User logged in successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var newUser models.UserInfo

	if err := ctx.BodyParser(&newUser); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.UpdateUser(c, &pb.UpdateUserRequest{
		UserId: int32(id),
		User: &pb.UserInfo{
			Email:         newUser.Email,
			PhoneNumber:   int32(newUser.PhoneNumber),
			FirstName:     newUser.FirstName,
			LastName:      newUser.LastName,
			AcademicGroup: newUser.AcademicGroup,
		},
	})

	if err != nil {
		slog.Errorf("Error updating user: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("User updated successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res.Message})
}

func (ctrl *UserController) Logout(ctx *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	slog.Info("User logged out successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "user logged out successfully"})
}
