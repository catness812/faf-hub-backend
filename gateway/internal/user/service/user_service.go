package service

import (
	"fmt"

	"github.com/gookit/slog"
)

type IRedisRepository interface {
	Save(key string, value string) error
}

type UserService struct {
	redisRepo IRedisRepository
}

func NewUserService(redisRepo IRedisRepository) *UserService {
	return &UserService{
		redisRepo: redisRepo,
	}
}

func (svc *UserService) SaveJWT(userID int, jwt string) error {
	if err := svc.redisRepo.Save(jwt, fmt.Sprintf("%d", userID)); err != nil {
		slog.Errorf("Could not save JWT in Redis: %v", err)
		return err
	}
	return nil
}
