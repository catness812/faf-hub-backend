package redis

import (
	"os"
	"strconv"

	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
)

func Connect() *redis.Client {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		slog.Fatalf("Failed to parse REDIS_DB: %v\n", err)
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	if client == nil {
		slog.Fatal("Failed to connect to Redis")
	}

	slog.Info("Successfully connected to Redis")
	return client
}
