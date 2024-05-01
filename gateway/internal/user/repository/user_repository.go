package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisClient *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{redisClient: redisClient}
}

func (repo *RedisRepository) Save(key string, value string) error {
	return repo.redisClient.Set(context.Background(), key, value, 0).Err()
}
