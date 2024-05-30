package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisClient *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{redisClient: redisClient}
}

func (repo *RedisRepository) Save(key string, value string) error {
	return repo.redisClient.Set(context.Background(), key, value, 10*time.Minute).Err()
}

func (repo *RedisRepository) SAdd(key string, value string) error {
	return repo.redisClient.SAdd(context.Background(), key, value).Err()
}

func (repo *RedisRepository) SRem(key string, value string) error {
	return repo.redisClient.SRem(context.Background(), key, value).Err()
}

func (repo *RedisRepository) GetValuesFromSet(key string) ([]string, error) {
	return repo.redisClient.SMembers(context.Background(), key).Result()
}

func (repo *RedisRepository) CheckKeyAndCompare(key string, value string) error {
	val, err := repo.redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key doesn't exist")
	} else if err != nil {
		return err
	}

	if val == value {
		return nil
	}

	return fmt.Errorf("value isn't right")
}

func (repo *RedisRepository) Delete(key string) error {
	if err := repo.redisClient.Del(context.Background(), key).Err(); err != nil {
		return err
	}
	return nil
}
