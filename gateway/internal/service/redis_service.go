package service

import (
	"fmt"

	"github.com/gookit/slog"
)

type IRedisRepository interface {
	Save(key string, value string) error
	CheckKeyAndCompare(key string, value string) error
	SAdd(key string, value string) error
	SRem(key string, value string) error
	GetValuesFromSet(key string) ([]string, error)
	Delete(key string) error
}

type RedisService struct {
	redisRepo IRedisRepository
}

func NewRedisService(redisRepo IRedisRepository) *RedisService {
	return &RedisService{
		redisRepo: redisRepo,
	}
}

func (svc *RedisService) Save(email string, passcode string) error {
	if err := svc.redisRepo.Save(email, passcode); err != nil {
		slog.Errorf("Could not save info in Redis: %v", err)
		return err
	}
	slog.Info("Successfully saved info in Redis")
	return nil
}

func (svc *RedisService) ValidatePasscode(email string, passcode string) error {
	if err := svc.redisRepo.CheckKeyAndCompare(email, passcode); err != nil {
		if err.Error() == "key doesn't exist" {
			slog.Error("Could not validate passcode: passcode expired")
			return fmt.Errorf("passcode expired")
		} else if err.Error() == "value isn't right" {
			slog.Error("Could not validate passcode: incorrect passcode")
			return fmt.Errorf("incorrect passcode")
		} else {
			return err
		}
	}

	if err := svc.redisRepo.Delete(email); err != nil {
		slog.Errorf("Could not delete from Redis: %v", err)
		return err
	}

	slog.Info("Successfully validated email verification passcode")
	return nil
}

func (svc *RedisService) Subscribe(email string) error {
	if err := svc.redisRepo.SAdd("newsletter", email); err != nil {
		slog.Errorf("Could not save info in Redis: %v", err)
		return err
	}
	slog.Info("Successfully updated info in Redis")
	return nil
}

func (svc *RedisService) Unsubscribe(email string) error {
	if err := svc.redisRepo.SRem("newsletter", email); err != nil {
		slog.Errorf("Could not save info in Redis: %v", err)
		return err
	}
	slog.Info("Successfully updated info in Redis")
	return nil
}

func (svc *RedisService) GetNewsletterEmails() ([]string, error) {
	emails, err := svc.redisRepo.GetValuesFromSet("newsletter")
	if err != nil {
		slog.Errorf("Could not retrieve emails from Redis: %v", err)
		return nil, err
	}
	slog.Info("Successfully retrieved emails from Redis")
	return emails, nil
}
