package util

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/catness812/faf-hub-backend/gateway/internal/user/pb"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * time.Duration(12)).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_PRIVATE_KEY")))
}

func ValidateJWT(ctx *fiber.Ctx) error {
	token, err := getToken(ctx)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func CurrentUserID(ctx *fiber.Ctx) (int, error) {
	err := ValidateJWT(ctx)
	if err != nil {
		return 0, err
	}
	token, err := getToken(ctx)
	if err != nil {
		return 0, nil
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	return int(userID), nil
}

func getToken(ctx *fiber.Ctx) (*jwt.Token, error) {
	tokenString := ctx.Cookies("jwt")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_PRIVATE_KEY")), nil
	})
	return token, err
}

// needs reformat
func CheckAdmin(ctx *fiber.Ctx) (bool, error) {
	userID, err := CurrentUserID(ctx)
	if err != nil {
		return false, err
	}

	conn, err := grpc.Dial(os.Getenv("APP_HOST")+":"+os.Getenv("USER_SVC_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err
	}
	defer conn.Close()

	c, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	res, err := pb.NewUserServiceClient(conn).CheckAdmin(c, &pb.CheckAdminRequest{
		UserId: int32(userID),
	})

	if err != nil {
		return false, err
	}

	return res.Admin, nil
}
