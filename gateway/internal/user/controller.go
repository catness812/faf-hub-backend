package user

import (
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/catness812/faf-hub-backend/gateway/internal/config"
	pb2 "github.com/catness812/faf-hub-backend/gateway/internal/notification/pb"
	"github.com/catness812/faf-hub-backend/gateway/internal/user/pb"
	"github.com/catness812/faf-hub-backend/gateway/internal/util"
	"github.com/catness812/faf-hub-backend/gateway/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type IRedisService interface {
	Save(email string, passcode string) error
	ValidatePasscode(email string, passcode string) error
	Subscribe(email string) error
	Unsubscribe(email string) error
}

type UserController struct {
	client             pb.UserServiceClient
	notificationClient pb2.NotificationServiceClient
	redisSvc           IRedisService
}

func NewUserController(client pb.UserServiceClient, notificationClient pb2.NotificationServiceClient, redisSvc IRedisService) *UserController {
	return &UserController{
		client:             client,
		notificationClient: notificationClient,
		redisSvc:           redisSvc,
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

func (ctrl *UserController) SendVerification(ctx *fiber.Ctx) error {
	id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetUser(c, &pb.GetUserRequest{
		UserId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error retrieving user: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	passcode := strconv.Itoa(rand.Intn(900000) + 100000)
	if err := ctrl.redisSvc.Save(res.User.Email, passcode); err != nil {
		slog.Errorf("Error saving verification credentials: %v", err.Error())
	}

	body := res.User.Email + ";" + passcode

	_, err = ctrl.notificationClient.Publish(c, &pb2.PublishRequest{
		QueueName: "verification",
		Body:      body,
	})
	if err != nil {
		slog.Errorf("Error publishing verification email: %v", err.Error())
	}

	slog.Info("User verification sent successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "user verification sent"})
}

func (ctrl *UserController) CompleteVerification(ctx *fiber.Ctx) error {
	type Passcode struct {
		Passcode string `json:"passcode"`
	}

	var passcode Passcode

	if err := ctx.BodyParser(&passcode); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetUser(c, &pb.GetUserRequest{
		UserId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error retrieving user: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ctrl.redisSvc.ValidatePasscode(res.User.Email, passcode.Passcode); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res2, err := ctrl.client.VerifyUser(c, &pb.VerifyRequest{
		UserId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error updating verified user status: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("User verified successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": res2.Message})
}

func (ctrl *UserController) Subscribe(ctx *fiber.Ctx) error {
	id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetUser(c, &pb.GetUserRequest{
		UserId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error retrieving user: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ctrl.redisSvc.Subscribe(res.User.Email); err != nil {
		slog.Errorf("Error subscribing user: %v", err.Error())
		return err
	}

	slog.Info("User subscribed to newsletter successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "user subscribed to newsletter"})
}

func (ctrl *UserController) Unsubscribe(ctx *fiber.Ctx) error {
	id, err := util.CurrentUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := ctrl.client.GetUser(c, &pb.GetUserRequest{
		UserId: int32(id),
	})

	if err != nil {
		slog.Errorf("Error retrieving user: %v", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ctrl.redisSvc.Unsubscribe(res.User.Email); err != nil {
		slog.Errorf("Error unsubscribing user: %v", err.Error())
		return err
	}

	slog.Info("User unsubscribed from newsletter successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "user unsubscribed to newsletter"})
}
